package functions

import (
	"context"
	"fmt"
	"fy-novel/internal/config"
	ollamaTool "fy-novel/internal/tools/ollama"
	progressTool "fy-novel/internal/tools/progress"
	"fy-novel/pkg/utils"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type FyChatbot struct {
	log             *logrus.Logger
	deepseekMessage []deepseek.ChatCompletionMessage
}

func NewFyChatbot(l *logrus.Logger) *FyChatbot {
	return &FyChatbot{log: l}
}

func (f *FyChatbot) FindOllamaContainer(
	ctx context.Context,
) (has bool, isInit bool, isSetModel bool, err error) {

	// 首先判断是不是已经在设置模型了
	completedInit, totalInit, existsInit := progressTool.GetProgress(
		ollamaTool.OllamaInitConTaskKey,
	)
	isInit = (existsInit && (completedInit < totalInit))
	completedModel, totalModel, existsModel := progressTool.GetProgress(
		ollamaTool.OllamaInitConSetModelTaskKey,
	)
	isSetModel = (existsModel && (completedModel < totalModel))
	if !isInit || !isSetModel {
		has, err = ollamaTool.FindOllamaContainer(ctx)
	} else {
		has = true
	}
	return has, isInit, isSetModel, err
}

func (f *FyChatbot) InitOllama(ctx context.Context) {
	go ollamaTool.InitOllamaContainer(ctx)
}

func (f *FyChatbot) GetInitOllamaProgress(ctx context.Context) (int64, int64, bool) {
	return progressTool.GetProgress(ollamaTool.OllamaInitConTaskKey)
}

func (f *FyChatbot) InitSetOllamaModelTask(ctx context.Context) error {
	return ollamaTool.InitSetOllamaModelTask(ctx)
}

func (f *FyChatbot) SetOllamaModel(ctx context.Context, model string) {
	go ollamaTool.OllamaContainerSetModel(ctx, model)
}

func (f *FyChatbot) GetCurrentUseModel(ctx context.Context) string {
	conf := config.GetConf()
	return conf.Chatbot.Model
}

func (f *FyChatbot) GetSelectModelList(ctx context.Context) []string {
	return ollamaTool.OllamaModelTypes.GetModels()
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

func (f *FyChatbot) CreateChatCompletion(ctx context.Context, inputText, deepseekApiKey string) (string, error) {
	// Set up the Deepseek client
	client := deepseek.NewClient(deepseekApiKey)

	if f.deepseekMessage == nil {
		f.deepseekMessage = []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: inputText},
		}
	} else {
		f.deepseekMessage = append(
			f.deepseekMessage,
			deepseek.ChatCompletionMessage{
				Role:    constants.ChatMessageRoleUser,
				Content: inputText},
		)
	}

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Messages: f.deepseekMessage,
	}

	// Send the request and handle the response
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", fmt.Errorf("deepseek CreateChatCompletion failed: %w", err)
	} else {
		f.deepseekMessage = append(
			f.deepseekMessage,
			deepseek.ChatCompletionMessage{
				Role:    constants.ChatMessageRoleAssistant,
				Content: response.Choices[0].Message.Content},
		)
	}
	utils.PrintAsJSON(response.Choices[0].Message.Content)
	return response.Choices[0].Message.Content, nil
}
