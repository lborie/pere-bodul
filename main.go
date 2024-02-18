package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	"github.com/lborie/pere-bodul/ai"
	"github.com/lborie/pere-bodul/serve"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"github.com/sirupsen/logrus"
)

func main() {
	// Instantiate OpenAI Client
	openAIKey := os.Getenv("OPENAI_KEY")
	if openAIKey != "" {
		ai.OpenAI = ai.OpenAIClient{
			OpenAIKey: openAIKey,
		}
	}

	// Instantiate GCP Client
	gcpKey := os.Getenv("GCP_KEY")
	gcpProject := os.Getenv("GCP_PROJECT_ID")
	if gcpKey != "" && gcpProject != "" {
		predictClient, err := aiplatform.NewPredictionClient(context.Background(), option.WithCredentialsJSON([]byte(gcpKey)), option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"))
		if err != nil {
			log.Fatalf("can't instantiate GCP client : %s", err.Error())
			return
		}
		// Instantiates a client.
		textoToSpeechClient, err := texttospeech.NewClient(context.Background(), option.WithCredentialsJSON([]byte(gcpKey)))
		if err != nil {
			log.Fatalf("error instantiating gcp client : %s", err.Error())
		}
		defer func(client *texttospeech.Client, client2 *aiplatform.PredictionClient) {
			_ = client.Close()
			_ = client2.Close()
		}(textoToSpeechClient, predictClient)

		ai.VertexAI = &ai.GCPClient{
			PredictionClient:   predictClient,
			TextToSpeechClient: textoToSpeechClient,
			PredictURL:         fmt.Sprintf("projects/%s/locations/us-central1/publishers/google/models/text-bison", gcpProject),
		}
	}

	// If no client, panic
	if ai.OpenAI == nil && ai.VertexAI == nil {
		log.Fatal("no client instantiated")
		return
	}

	http.HandleFunc("/generateStory", serve.StoryHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Serving index.html with user-agent : %s", r.UserAgent())
		http.ServeFile(w, r, "serve/index.html")
	})
	http.HandleFunc("/background.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "serve/background.jpg")
	})

	fmt.Println("Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
