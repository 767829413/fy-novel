package main

import (
	"context"
	"fy-novel/internal/functions"
	"fy-novel/internal/model"
	"os"

	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx          context.Context
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
}

func (a *App) SerachNovel(name string) ([]*model.SearchResult, error) {
	return a.downloader.Serach(name)
}

func (a *App) DownLoadNovel(sr *model.SearchResult) (*model.CrawlResult, error) {
	return a.downloader.DownLoad(sr)
}

func (a *App) CheckUpdate() string {
	return a.checkUpdater.CheckUpdate()
}

func (a *App) GetUsageInfo() []string {
	return a.getHint.GetUsageInfo()
}

func (a *App) GetConfig() (string, error) {
	return a.getConf.GetConfigString()
}

func (a *App) GetDownloadProgress() int {
	return 0
}
