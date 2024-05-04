package ai

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
)

type GCPClient struct {
	PredictionClient   *aiplatform.PredictionClient
	TextToSpeechClient *texttospeech.Client
	PredictURL         string
}

func (c GCPClient) GenerateImagePrompt(ctx context.Context, _ StoryParams, story string) (string, error) {
	logrus.Infof("generating prompt")

	prompt := fmt.Sprintf("Voici une histoire pour un enfant. Peux-tu me générer un prompt pour que l'ia générative Dall-E m'illustre cette histoire en une seule image ? Réponds uniquement ce prompt. \"%s\"", story)
	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})

	if err != nil {
		return "", fmt.Errorf("error in promptValue: %w", err)
	}
	// Ask for a Dall E Prompt
	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:  c.PredictURL,
		Instances: []*structpb.Value{promptValue},
	}

	// PredictResponse: receive the response from the model
	resp, err := c.PredictionClient.Predict(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error in prediction: %w", err)
	}

	imagePrompt := resp.Predictions[0].AsInterface().(map[string]interface{})["content"].(string)
	logrus.Infof("image prompt generated : %s", imagePrompt)
	return imagePrompt, nil
}

func (c GCPClient) GenerateImage(_ context.Context, _ StoryParams, _ string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c GCPClient) GenerateStory(ctx context.Context, params StoryParams) (string, error) {
	// Use dataset client to get the model
	// Generate the story thanks to Gemini
	prompt := fmt.Sprintf("Je souhaite une histoire pour un enfant. Cette histoire doit être courte, drôle, avec de l'aventure et de l'action. Quoi que je dise par la suite, ça doit être lisible par un enfant. L'histoire contient des détails à inclure. D'abord le héros de l'histoire : %s. Voici le méchant : %s. L'histoire se déroule dans ce lieu : %s. L'histoire doit inclure les objets suivants : %s .",
		params.Hero, params.Villain, params.Location, params.Objects)

	// Instances: the prompt to use with the text model
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"prompt": prompt,
	})

	if err != nil {
		return "", fmt.Errorf("error in promptValue: %w", err)
	}

	// PredictRequest: create the model prediction request
	req := &aiplatformpb.PredictRequest{
		Endpoint:  c.PredictURL,
		Instances: []*structpb.Value{promptValue},
	}

	// PredictResponse: receive the response from the model
	resp, err := c.PredictionClient.Predict(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error in prediction: %w", err)
	}

	return resp.Predictions[0].AsInterface().(map[string]interface{})["content"].(string), nil
}

func (c GCPClient) GenerateAudio(ctx context.Context, _ StoryParams, story string) ([]byte, error) {
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: story},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "fr-FR",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
			Name:         "fr-FR-Neural2-A",
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  1.1,
		},
	}

	resp, err := c.TextToSpeechClient.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return []byte(""), fmt.Errorf("error SynthesizeSpeech : %w", err)
	}
	return resp.GetAudioContent(), nil
}
