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
