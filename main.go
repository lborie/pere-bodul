package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type StoryParams struct {
	Age      int
	Hero     string
	Villain  string
	Location string
	Objects  []string
}

type GPT3Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

var openAIKey string

func generateStory(params StoryParams) (string, error) {
    prompt := fmt.Sprintf("Créez une histoire pour un enfant de %d ans. Le héros de l'histoire est %s. Le méchant est %s. L'histoire se déroule à %s. L'histoire doit inclure les objets suivants : %s.",
        params.Age, params.Hero, params.Villain, params.Location, strings.Join(params.Objects, ", "))

    requestBody, _ := json.Marshal(map[string]interface{}{
        "model": "gpt-3.5-turbo",
        "messages": []map[string]string{
            {
                "role": "system",
                "content": "Créer une histoire.",
            },
            {
                "role": "user",
                "content": prompt,
            },
        },
    })

    request, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", "Bearer " + openAIKey)

    client := &http.Client{}
    response, err := client.Do(request)

    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return "", errors.New("La requête à GPT-3 a échoué avec le status : " + response.Status)
    }

    body, _ := io.ReadAll(response.Body)

    var gpt3Response GPT3Response
    json.Unmarshal(body, &gpt3Response)

    if len(gpt3Response.Choices) > 0 {
        return gpt3Response.Choices[0].Text, nil
    }

    return "", errors.New("GPT-3 n'a pas renvoyé de texte")
}
func storyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
			return
		}

		age, err := strconv.Atoi(r.Form.Get("age"))
		if err != nil {
			http.Error(w, "Erreur lors de la conversion de l'âge en entier", http.StatusBadRequest)
			return
		}

		hero := r.Form.Get("hero")
		villain := r.Form.Get("villain")
		location := r.Form.Get("location")

		objectsStr := r.Form.Get("objects")
		objects := strings.Split(objectsStr, ",")

		storyParams := StoryParams{
			Age:      age,
			Hero:     hero,
			Villain:  villain,
			Location: location,
			Objects:  objects,
		}

		story, err := generateStory(storyParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erreur lors de la génération de l'histoire : %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(story))
	} else {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {

	openAIKey = os.Getenv("OPENAI_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_KEY n'est pas défini")
	}

	http.HandleFunc("/generateStory", storyHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("Le serveur tourne sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
