package ai

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type VertexAIClient struct {
	GenAIClient        *genai.Client
	TextToSpeechClient *texttospeech.Client
}

func (c VertexAIClient) GenerateImagePrompt(ctx context.Context, params StoryParams, story string) (string, error) {
	logrus.Infof("generating prompt")

	prompt := fmt.Sprintf("Voici une histoire pour un enfant. Peux-tu me générer un prompt pour que l'ia générative Dall-E m'illustre cette histoire en une seule image ? Réponds uniquement ce prompt. \"%s\"", story)
	// Ask for a Dall E Prompt
	gemini := c.GenAIClient.GenerativeModel(string(params.Wizard))
	resp, err := gemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error during image prompt generation : %w", err)
	}
	imagePrompt := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
	logrus.Infof("image prompt generated : %s", imagePrompt)
	return imagePrompt, nil
}

func (c VertexAIClient) GenerateImage(_ context.Context, _ StoryParams, _ string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c VertexAIClient) GenerateStory(ctx context.Context, params StoryParams) (string, error) {
	// Generate the story thanks to Gemini
	prompt := fmt.Sprintf("Je souhaite une histoire pour un enfant. Cette histoire doit être courte, drôle, avec de l'aventure et de l'action. Quoi que je dise par la suite, ça doit être lisible par un enfant. L'histoire contient des détails à inclure. D'abord le héros de l'histoire : %s. Voici le méchant : %s. L'histoire se déroule dans ce lieu : %s. L'histoire doit inclure les objets suivants : %s .",
		params.Hero, params.Villain, params.Location, params.Objects)
	gemini := c.GenAIClient.GenerativeModel(string(params.Wizard))
	resp, err := gemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}
	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
}

func (c VertexAIClient) GenerateAudio(ctx context.Context, _ StoryParams, story string) ([]byte, error) {
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
