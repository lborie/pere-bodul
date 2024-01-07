package ai

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"fmt"
	"google.golang.org/api/option"
)

type GCPClient struct {
	GCPKey        string
	DataSetClient *aiplatform.DatasetClient
}

func (c GCPClient) GenerateImage(ctx context.Context, story string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c GCPClient) GenerateStory(params StoryParams) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c GCPClient) GenerateAudio(ctx context.Context, story string) ([]byte, error) {
	// Instantiates a client.
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsJSON([]byte(c.GCPKey)))
	if err != nil {
		return []byte(""), fmt.Errorf("error instantiating gcp client : %w", err)
	}
	defer func(client *texttospeech.Client) {
		_ = client.Close()
	}(client)

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

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return []byte(""), fmt.Errorf("error SynthesizeSpeech : %w", err)
	}
	return resp.GetAudioContent(), nil
}
