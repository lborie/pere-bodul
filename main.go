package main

import (
	"context"
	"fmt"
	"github.com/lborie/pere-bodul/ai"
	"github.com/lborie/pere-bodul/serve"
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
	if gcpKey != "" {
		client, err := aiplatform.NewDatasetClient(context.Background())
		if err != nil {
			log.Fatalf("can't instantiate GCP client : %s", err.Error())
			return
		}
		ai.VertexAI = &ai.GCPClient{
			GCPKey:        gcpKey,
			DataSetClient: client,
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
		http.ServeFile(w, r, "front/index.html")
	})
	http.HandleFunc("/background.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "front/background.jpg")
	})

	fmt.Println("Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
