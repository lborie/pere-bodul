package ai

import (
	"context"
)

type Wizard string

const GeminiWizard Wizard = "gemini-1.5-flash-001"
const TextBizonWizard Wizard = "text-bison"
const OpenAIWizard Wizard = "OpenAI"
const Llama3Wizard Wizard = "llama3"

type StoryParams struct {
	Hero      string
	Villain   string
	Location  string
	Objects   string
	Scene     *[]byte
	SceneType string
	Wizard    Wizard
}

type PereBodulClient interface {
	GenerateStory(ctx context.Context, params StoryParams) (string, error)
	GenerateAudio(ctx context.Context, params StoryParams, story string) ([]byte, error)
	GenerateImage(ctx context.Context, params StoryParams, imagePrompt string) (string, error)
	GenerateImagePrompt(ctx context.Context, params StoryParams, story string) (string, error)
}

var OpenAI PereBodulClient
var AIPlatform PereBodulClient // Deprecated
var VertexAI PereBodulClient
