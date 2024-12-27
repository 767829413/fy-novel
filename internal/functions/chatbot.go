package functions

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type FyChatbot struct {
	log *logrus.Logger
}

func NewFyChatbot(l *logrus.Logger) *FyChatbot {
	return &FyChatbot{log: l}
}

func (f *FyChatbot) StartChatbot(ctx context.Context, userInput string) (string, error) {
	llm, err := ollama.New(ollama.WithModel("llama2"))

	if err != nil {
		return "", fmt.Errorf("failed to create Ollama client after multiple attempts: %w", err)
	}

	completion, err := llm.Call(
		ctx,
		userInput,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			return nil
		}),
	)
	if err != nil {
		return "", fmt.Errorf("Ollama API call failed: %w", err)
	}

	return completion, nil
}
