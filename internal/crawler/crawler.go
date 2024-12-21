package crawler

import (
	"fmt"
	"fy-novel/internal/config"
	"fy-novel/internal/model"
	"fy-novel/internal/parse"
	chapterTool "fy-novel/internal/tools/chapter"
	mergeTool "fy-novel/internal/tools/merge"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Crawler interface {
	Search(key string) ([]*model.SearchResult, error)
	Crawl(res *model.SearchResult, start, end int) (*model.CrawlResult, error)
}

type novelCrawler struct{}

func NewNovelCrawler() Crawler {
	return &novelCrawler{}
}

func (nc *novelCrawler) Search(key string) ([]*model.SearchResult, error) {
	conf := config.GetConf()
	// 解析
	res, err := parse.NewSearchResultParser(conf.Base.SourceID).Parse(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (nc *novelCrawler) Crawl(res *model.SearchResult, start, end int) (*model.CrawlResult, error) {
	conf := config.GetConf()
	// 小说详情页抓取解析
	book, err := parse.NewBookParser(conf.Base.SourceID).Parse(res.Url)
	if err != nil {
		return nil, err
	}

	// Format the directory name as "BookName (Author)"
	bookDir := fmt.Sprintf("%s (%s)", book.BookName, book.Author)
	dirPath := filepath.Join(conf.Base.DownloadPath, bookDir)
	// Create the directory
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	// 获取小说目录
	catalogsParser := parse.NewCatalogsParser(conf.Base.SourceID)
	catalogs, err := catalogsParser.Parse(res.Url, start, end)
	if err != nil {
		return nil, err
	}
	if len(catalogs) == 0 {
		return nil, nil
	}

	startTime := time.Now()
	// 解析下载内容
	var wg sync.WaitGroup
	cpuNum := runtime.NumCPU()
	threads := conf.Crawl.Threads
	if threads == -1 {
		threads = cpuNum * 2
	}
	semaphore := make(chan struct{}, threads)
	var nowCatalogsCount int32 = int32(len(catalogs))
	// process.SetSaveProgress(res.Url, &process.Pressing{All: nowCatalogsCount})
	for _, chapter := range catalogs {
		wg.Add(1)
		go func(chapter *model.Chapter, bookDir string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			// 下载逻辑
			parse.NewChapterParser(conf.Base.SourceID).Parse(chapter, res, book, bookDir)
			chapterTool.CreateFileForChapter(chapter, bookDir)
			atomic.AddInt32(&nowCatalogsCount, -1)
		}(
			chapter,
			bookDir,
		)
	}
	wg.Wait()

	// 合并生成小说文件格式
	outputPath, err := mergeTool.MergeSaveHandler(book, dirPath)
	if err != nil {
		return nil, err
	}
	return &model.CrawlResult{
		OutputPath: outputPath,
		TakeTime:   int64(time.Since(startTime).Seconds()),
	}, nil
}
