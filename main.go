package main

import (
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

var openAIKey string
var gcpKey string

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

		story, err := GenerateStory(storyParams, openAIKey)
		if err != nil {
			logrus.Errorf("erreur lors de la génération de l'histoire : %s", err.Error())
			http.Error(w, "Erreur lors de la génération de l'histoire", http.StatusInternalServerError)
			return
		}

		logrus.Infof("story generated : %s", story)

		audio, err := GenerateAudio(r.Context(), story, gcpKey)
		if err != nil {
			logrus.Errorf("erreur lors de la génération du son : %s", err.Error())
			http.Error(w, "Erreur lors de la génération du son", http.StatusInternalServerError)
			return
		}

		// Send JSON
		w.Header().Set("Content-Type", "application/json")
		data := map[string]any{
			"story": story,
			"audio": audio,
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
	openAIKey = os.Getenv("OPENAI_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_KEY n'est pas défini")
		return
	}
	gcpKey = os.Getenv("GCP_KEY")
	if gcpKey == "" {
		log.Fatal("GCP_KEY n'est pas défini")
		return
	}

	http.HandleFunc("/generateStory", storyHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Serving index.html")
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/background.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "background.png")
	})

	fmt.Println("Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
