package functions

import (
	"context"
	"fmt"
	"fy-novel/internal/config"
	ollamaTool "fy-novel/internal/tools/ollama"
	progressTool "fy-novel/internal/tools/progress"

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

func (f *FyChatbot) FindOllamaContainer(ctx context.Context) (bool, error) {
	return ollamaTool.FindOllamaContainer(ctx)
}

func (f *FyChatbot) InitOllama(ctx context.Context) error {
	return ollamaTool.InitOllamaContainer(ctx)
}

func (f *FyChatbot) GetInitOllamaProgress(ctx context.Context) (int64, int64, bool) {
	return progressTool.GetProgress(ollamaTool.OllamaInitConTaskKey)
}

func (f *FyChatbot) SetOllamaModel(ctx context.Context, model string) error {
	return ollamaTool.OllamaContainerSetModel(ctx, model)
}

func (f *FyChatbot) GetSetOllamaModelProgress(ctx context.Context) (int64, int64, bool) {
	return progressTool.GetProgress(ollamaTool.OllamaInitConSetModelTaskKey)
}

func (f *FyChatbot) StartChatbot(ctx context.Context, userInput string) (string, error) {
	conf := config.GetConf()
	llm, err := ollama.New(ollama.WithModel(conf.Chatbot.Model))

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
