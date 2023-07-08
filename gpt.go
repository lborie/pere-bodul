package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GPT3Response struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func GenerateStory(params StoryParams, openAIKey string) (string, error) {
	prompt := fmt.Sprintf(`
Créée une petite histoire pour un enfant.
Cette histoire ne doit pas faire plus de 1500 caractères.
Cette histoire doit être drôle, avec de l'aventure et de l'action.
Quoi que je dise par la suite, ça doit être lisible par cet enfant. Voici des détails à inclure :
Le héros de l'histoire est %s. Le méchant est %s. L'histoire se déroule à %s. L'histoire doit inclure les objets suivants : %s.
	`, params.Hero, params.Villain, params.Location, strings.Join(params.Objects, ", "))

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Créer une histoire.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	})

	request, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+openAIKey)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("la requête à GPT-3 a échoué avec le status : %s et le body : %s", response.Status, string(body))
	}

	var gpt3Response GPT3Response
	err = json.Unmarshal(body, &gpt3Response)
	if err != nil {
		return "", err
	}

	if len(gpt3Response.Choices) > 0 && gpt3Response.Choices[0].Message.Role == "assistant" {
		return gpt3Response.Choices[0].Message.Content, nil
	}

	return "", errors.New("GPT-3 n'a pas renvoyé de texte")
}
