package main

import (
	"context"
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
	getConf      *functions.GetConf
	getHint      *functions.GetHint
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
	a.getConf = functions.NewGetConf(log)
	a.getHint = functions.NewGetHint(log)
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
	res.Config = a.getConf.GetConfig()
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
