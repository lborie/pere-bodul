package ai

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
)

type GeminiClient struct {
	GenAIClient        *genai.Client
	TextToSpeechClient *texttospeech.Client
}

func (c GeminiClient) GenerateImage(ctx context.Context, story string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c GeminiClient) GenerateStory(ctx context.Context, params StoryParams) (string, error) {
	// Generate the story thanks to Gemini
	prompt := fmt.Sprintf("Je souhaite une histoire pour un enfant. Cette histoire doit être courte, drôle, avec de l'aventure et de l'action. Quoi que je dise par la suite, ça doit être lisible par un enfant. L'histoire contient des détails à inclure. D'abord le héros de l'histoire : %s. Voici le méchant : %s. L'histoire se déroule dans ce lieu : %s. L'histoire doit inclure les objets suivants : %s .",
		params.Hero, params.Villain, params.Location, params.Objects)

	//gemini := c.GenAIClient.GenerativeModel("gemini-pro-vision")
	gemini := c.GenAIClient.GenerativeModel("gemini-pro")
	resp, err := gemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}
	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
}

func (c GeminiClient) GenerateAudio(ctx context.Context, story string) ([]byte, error) {
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
