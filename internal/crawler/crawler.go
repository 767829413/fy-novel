package crawler

import (
	"fmt"
	"fy-novel/internal/config"
	"fy-novel/internal/model"
	"fy-novel/internal/parse"
	chapterTool "fy-novel/internal/tools/chapter"
	mergeTool "fy-novel/internal/tools/merge"
	progressTool "fy-novel/internal/tools/progress"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Crawler interface {
	Search(key string) ([]*model.SearchResult, error)
	Crawl(res *model.SearchResult, start, end int) (*model.CrawlResult, error)
}

type novelCrawler struct {
	log *logrus.Logger
}

func NewNovelCrawler(log *logrus.Logger) Crawler {
	return &novelCrawler{
		log: log,
	}
}

func (nc *novelCrawler) Search(key string) ([]*model.SearchResult, error) {
	// Parse
	res, err := parse.NewSearchResultParser(config.GetConf()).Parse(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (nc *novelCrawler) Crawl(res *model.SearchResult, start, end int) (*model.CrawlResult, error) {
	conf := config.GetConf()
	// Fetch and parse the novel details page
	book, err := parse.NewBookParser(config.GetConf()).Parse(res.Url)
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
	// Get the novel's table of contents
	catalogsParser := parse.NewCatalogsParser(conf)
	catalogs, err := catalogsParser.Parse(res.Url, start, end)
	if err != nil {
		return nil, err
	}
	if len(catalogs) == 0 {
		return nil, nil
	}

	startTime := time.Now()
	// Parse and download content
	// Limit concurrent processing
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, conf.GetConcurrencyNum())
	var nowCatalogsCount = int64(0)
	// Total completed tasks = number of chapters fetched + 1 (merging task)
	progressTool.InitTask(res.Url, int64(len(catalogs)+1))
	for _, chapter := range catalogs {
		wg.Add(1)
		go func(chapter *model.Chapter, bookDir string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			defer func() {
				atomic.AddInt64(&nowCatalogsCount, 1)
				progressTool.UpdateProgress(res.Url, nowCatalogsCount)
			}()
			// Download logic
			err := parse.NewChapterParser(conf).Parse(chapter, res, book, bookDir)
			if err != nil {
				fmt.Printf("parse.NewChapterParser(conf.Base.SourceID).Parse error: %v", err)
				return
			}
			chapterTool.CreateFileForChapter(chapter, bookDir)
		}(
			chapter,
			bookDir,
		)
	}
	wg.Wait()

	// Merge and generate the novel file format
	outputPath, err := mergeTool.MergeSaveHandler(book, dirPath)
	if err != nil {
		return nil, err
	}
	// Task completed
	defer func() {
		atomic.AddInt64(&nowCatalogsCount, 1)
		progressTool.UpdateProgress(res.Url, int64(nowCatalogsCount))
	}()

	return &model.CrawlResult{
		OutputPath: outputPath,
		TakeTime:   int64(time.Since(startTime).Seconds()),
	}, nil
}
