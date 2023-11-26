package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type StoryParams struct {
	Hero     string
	Villain  string
	Location string
	Objects  string
}

type PereBodulClient interface {
	GenerateStory(params StoryParams) (string, error)
	GenerateAudio(ctx context.Context, story string) ([]byte, error)
	GenerateImage(ctx context.Context, story string) (string, error)
}

var pereBodulClient PereBodulClient

func storyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			logrus.Error("erreur lors de l'analyse du formulaire : %s", err.Error())
			http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
			return
		}

		storyParams := StoryParams{
			Hero:     r.Form.Get("hero"),
			Villain:  r.Form.Get("villain"),
			Location: r.Form.Get("location"),
			Objects:  r.Form.Get("objects"),
		}

		logrus.Infof("new form submitted: %v", storyParams)

		story, err := pereBodulClient.GenerateStory(storyParams)
		if err != nil {
			logrus.Errorf("erreur lors de la génération de l'histoire : %s", err.Error())
			http.Error(w, "Erreur lors de la génération de l'histoire", http.StatusInternalServerError)
			return
		}

		logrus.Infof("story generated : %s", story)

		// Parallelize steps 2 and 3 with channels
		audioChan := make(chan []byte, 1)
		defer close(audioChan)
		imageChan := make(chan string, 1)
		defer close(imageChan)

		go func(ctx context.Context, story string) {
			logrus.Info("generating audio")
			audio, err := pereBodulClient.GenerateAudio(ctx, story)
			if err != nil {
				logrus.Errorf("erreur lors de la génération du son : %s", err.Error())
				audioChan <- []byte("")
			}
			logrus.Info("audio generated")
			audioChan <- audio
		}(r.Context(), story)

		go func(ctx context.Context, story string) {
			logrus.Info("generating image")
			image, err := pereBodulClient.GenerateImage(ctx, story)
			if err != nil {
				logrus.Errorf("erreur lors de la génération de l'image : %s", err.Error())
				imageChan <- ""
			}
			logrus.Infof("image generated")
			imageChan <- image
		}(r.Context(), story)

		audio := <-audioChan
		image := <-imageChan

		logrus.Infof("image link : %s", image)

		// Send JSON
		w.Header().Set("Content-Type", "application/json")
		data := map[string]any{
			"story":    story,
			"audio":    audio,
			"imageUrl": image,
		}
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			logrus.Errorf("Erreur lors de la sérialisation en JSON : %s", err.Error())
			http.Error(w, "Erreur lors de la sérialisation en JSON ", http.StatusInternalServerError)
			return
		}
	} else {
		logrus.Error("Méthode non autorisée")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func main() {
	openAIKey := os.Getenv("OPENAI_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_KEY n'est pas défini")
		return
	}
	pereBodulClient = OpenAIClient{
		OpenAIKey: openAIKey,
	}
	/*
		gcpKey := os.Getenv("GCP_KEY")
		if gcpKey == "" {
			log.Fatal("GCP_KEY n'est pas défini")
			return
		}
		gcpClient = &GCPClient{
			GCPKey: gcpKey,
		}
	*/
	http.HandleFunc("/generateStory", storyHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Serving index.html with user-agent : %s", r.UserAgent())
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/background.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "background.png")
	})

	fmt.Println("Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
