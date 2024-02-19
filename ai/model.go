package ai

import "context"

type StoryParams struct {
	Hero     string
	Villain  string
	Location string
	Objects  string
}

type PereBodulClient interface {
	GenerateStory(ctx context.Context, params StoryParams) (string, error)
	GenerateAudio(ctx context.Context, story string) ([]byte, error)
	GenerateImage(ctx context.Context, imagePrompt string) (string, error)
	GenerateImagePrompt(ctx context.Context, story string) (string, error)
}

type PereBodulImpl string

var OpenAIImpl PereBodulImpl = "openai"
var GCPImpl PereBodulImpl = "gcp"
var GeminiImpl PereBodulImpl = "gemini"

var OpenAI PereBodulClient
var AIPlatform PereBodulClient
var Gemini PereBodulClient
