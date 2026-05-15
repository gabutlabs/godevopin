package ai

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

// NewLLMClient creates an LLM client based on the provider
func NewLLMClient(ctx context.Context, provider, apiKey, modelName, baseURL string) (llms.Model, error) {
	switch provider {
	case "openai":
		return openai.New(
			openai.WithToken(apiKey),
			openai.WithModel(modelName),
		)
	case "gemini":
		return googleai.New(ctx,
			googleai.WithAPIKey(apiKey),
			googleai.WithDefaultModel(modelName),
		)
	case "huggingface":
		return openai.New(
			openai.WithToken(apiKey),
			openai.WithModel(modelName),
			openai.WithBaseURL(baseURL),
		)
	case "custom":
		// For user-deployed models that are OpenAI-API compatible (e.g., Ollama, vLLM, LocalAI)
		return openai.New(
			openai.WithToken(apiKey),
			openai.WithModel(modelName),
			openai.WithBaseURL(baseURL),
		)
	default:
		return nil, errors.New("unsupported AI provider")
	}
}
