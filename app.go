package main

import (
	"context"
	"fmt"
	"fy-novel/internal/functions"
	"fy-novel/internal/model"
	progressTool "fy-novel/internal/tools/progress"
	"os"

	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx          context.Context
	log          *logrus.Logger
	checkUpdater *functions.CheckUpdater
	downloader   *functions.Downloader
	confHandler  *functions.ConfHandler
	getHint      *functions.GetHint
	chatbot      *functions.FyChatbot
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log := logrus.New()
	// Setup logger
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	// Set log level
	log.SetLevel(logrus.ErrorLevel)
	a.checkUpdater = functions.NewCheckUpdate(log, 5000)
	a.downloader = functions.NewDownload(log)
	a.confHandler = functions.NewGetConf(log)
	a.getHint = functions.NewGetHint(log)
	a.chatbot = functions.NewFyChatbot(log)
	a.log = log
}

func (a *App) SerachNovel(name string) []*model.SearchResult {
	res, err := a.downloader.Serach(name)
	if err != nil {
		a.log.Errorf("app SerachNovel error: %v", err)
		return nil
	}
	return res
}

func (a *App) DownLoadNovel(sr *model.SearchResult) *model.CrawlResult {
	res, err := a.downloader.DownLoad(sr)
	if err != nil {
		a.log.Errorf("app DownLoadNovel error: %v", err)
		return nil
	}
	return res
}

func (a *App) GetUpdateInfo() *model.GetUpdateInfoResult {
	return a.checkUpdater.CheckUpdate()
}

func (a *App) GetUsageInfo() *model.GetUsageInfoResult {
	return a.getHint.GetUsageInfo()
}

func (a *App) GetConfig() *model.GetConfigResult {
	res := &model.GetConfigResult{}
	res.Config = a.confHandler.GetConfig()
	return res
}

func (a *App) GetDownloadProgress(sr *model.SearchResult) *model.ProgressResult {
	res := &model.ProgressResult{}
	completed, total, exists := progressTool.GetProgress(sr.Url)
	if exists {
		res.Exists = exists
		res.Completed = int(completed)
		res.Total = int(total)
	}
	return res
}

func (a *App) SetConfig(conf string) string {
	if err := a.confHandler.SetConfig(conf); err != nil {
		return err.Error()
	}
	return ""
}

func (a *App) StartChatbot(userInput string) *model.StartChatbotResult {
	res := &model.StartChatbotResult{}
	resp, err := a.chatbot.StartChatbot(a.ctx, userInput)
	if err != nil {
		errMsg := fmt.Sprintf("app StartChatbot error: %v", err)
		a.log.Error(errMsg)
		res.ErrorMsg = errMsg
		return res
	}
	res.Response = resp
	return res
}

func (a *App) HasInitOllama() *model.HasInitOllamaResult {
	res := &model.HasInitOllamaResult{}
	has, isInit, isSetModel, err := a.chatbot.FindOllamaContainer(a.ctx)
	if err != nil {
		errMsg := fmt.Sprintf("app HasInitOllama error: %v", err)
		a.log.Error(errMsg)
		res.ErrorMsg = errMsg
		return res
	}
	res.Has = has
	res.IsInit = isInit
	res.IsSetModel = isSetModel
	return res
}

func (a *App) InitOllama() *model.InitOllamaResult {
	res := &model.InitOllamaResult{}
	a.chatbot.InitOllama(a.ctx)
	return res
}

func (a *App) GetInitOllamaProgress() *model.InitOllamaProgressResult {
	res := &model.InitOllamaProgressResult{}
	completed, total, exists := a.chatbot.GetInitOllamaProgress(a.ctx)
	if exists {
		res.Exists = exists
		res.Completed = int(completed)
		res.Total = int(total)
	}
	return res
}

func (a *App) InitSetOllamaModelTask() *model.InitSetOllamaModelResult {
	res := &model.InitSetOllamaModelResult{}
	err := a.chatbot.InitSetOllamaModelTask(a.ctx)
	if err != nil {
		res.ErrorMsg = err.Error()
	}
	return res
}

func (a *App) SetOllamaModel(modelName string) *model.SetOllamaModelResult {
	res := &model.SetOllamaModelResult{}
	a.chatbot.SetOllamaModel(a.ctx, modelName)
	return res
}

func (a *App) GetCurrentUseModel() *model.GetCurrentUseModelResult {
	res := &model.GetCurrentUseModelResult{}
	res.Model = a.chatbot.GetCurrentUseModel(a.ctx)
	return res
}

func (a *App) GetSelectModelList() *model.GetSelectModelListResult {
	res := &model.GetSelectModelListResult{}
	res.Models = a.chatbot.GetSelectModelList(a.ctx)
	return res
}

func (a *App) GetSetOllamaModelProgress() *model.GetSetOllamaModelProgressResult {
	res := &model.GetSetOllamaModelProgressResult{}
	completed, total, exists := a.chatbot.GetSetOllamaModelProgress(a.ctx)
	if exists {
		res.Exists = exists
		res.Completed = int(completed)
		res.Total = int(total)
	}
	return res
}

func (a *App) DeepSeekChat(inputText, deepseekApiKey string) *model.StartChatbotResult {
	res := &model.StartChatbotResult{}
	resp, err := a.chatbot.CreateChatCompletion(a.ctx, inputText, deepseekApiKey)
	if err != nil {
		errMsg := fmt.Sprintf("app DeepSeekChat error: %v", err)
		a.log.Error(errMsg)
		res.ErrorMsg = errMsg
		return res
	}
	res.Response = resp
	return res
}
