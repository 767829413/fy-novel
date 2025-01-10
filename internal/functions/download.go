package functions

import (
	"math"

	"fy-novel/internal/config"
	"fy-novel/internal/crawler"
	"fy-novel/internal/model"

	"github.com/sirupsen/logrus"
)

type Downloader struct {
	log     *logrus.Logger
	crawler crawler.Crawler
}

func NewDownload(l *logrus.Logger) *Downloader {
	conf := config.GetConf()
	return &Downloader{log: l, crawler: crawler.NewNovelCrawler(conf)}
}

func (d *Downloader) Serach(name string) ([]*model.SearchResult, error) {
	return d.crawler.Search(name)
}

func (d *Downloader) DownLoad(sr *model.SearchResult) (*model.CrawlResult, error) {
	start, end := 1, math.MaxInt // Max int
	return d.crawler.Crawl(sr, start, end)
}
