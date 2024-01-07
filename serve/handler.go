package serve

import (
	"context"
	"encoding/json"
	"github.com/lborie/pere-bodul/ai"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			logrus.Error("erreur lors de l'analyse du formulaire : %s", err.Error())
			http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
			return
		}

		storyParams := ai.StoryParams{
			Hero:     r.Form.Get("hero"),
			Villain:  r.Form.Get("villain"),
			Location: r.Form.Get("location"),
			Objects:  r.Form.Get("objects"),
		}

		logrus.Infof("new form submitted: %v", storyParams)

		// For this request, default client is GCPClient from main.go

		pereBodulClient := ai.VertexAI

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
