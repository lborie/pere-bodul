package ai

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"context"
	"google.golang.org/api/option"
	"os"
	"testing"
)

// Permits to list the voices available on GCP
func TestGenerateAudio(t *testing.T) {
	gcpKey := os.Getenv("GCP_KEY")
	client, err := texttospeech.NewClient(context.Background(), option.WithCredentialsJSON([]byte(gcpKey)))
	if err != nil {
		t.Errorf("error instantiating gcp client : %v", err)
	}
	defer func(client *texttospeech.Client) {
		_ = client.Close()
	}(client)

	voices, err := client.ListVoices(context.Background(), &texttospeechpb.ListVoicesRequest{
		LanguageCode: "fr-FR",
	})
	if err != nil {
		t.Errorf("error instantiating gcp client : %v", err)
	}
	for _, voice := range voices.Voices {
		if voice.SsmlGender == texttospeechpb.SsmlVoiceGender_FEMALE {
			t.Logf("voice : %v", voice.Name)
		}
	}
}
