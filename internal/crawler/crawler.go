package crawler

import (
	"fmt"
	"fy-novel/internal/config"
	"fy-novel/internal/model"
	"fy-novel/internal/parse"
	chapterTool "fy-novel/internal/tools/chapter"
	concurrencyTool "fy-novel/internal/tools/concurrency"
	mergeTool "fy-novel/internal/tools/merge"
	progressTool "fy-novel/internal/tools/progress"
	"os"
	"path/filepath"
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
	// 限制并发处理
	var wg sync.WaitGroup
	threads := concurrencyTool.GetConcurrencyNum(conf.Crawl.Threads)
	semaphore := make(chan struct{}, threads)
	var nowCatalogsCount = int32(0)
	// 总的完成任务数 = 抓取文章数量 + 1(合成任务)
	progressTool.InitTask(res.Url, int64(len(catalogs)+1))
	for _, chapter := range catalogs {
		wg.Add(1)
		go func(chapter *model.Chapter, bookDir string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			// 下载逻辑
			parse.NewChapterParser(conf.Base.SourceID).Parse(chapter, res, book, bookDir)
			chapterTool.CreateFileForChapter(chapter, bookDir)
			defer func() {
				atomic.AddInt32(&nowCatalogsCount, 1)
				progressTool.UpdateProgress(res.Url, int64(nowCatalogsCount))
			}()
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
	// 任务完成
	defer func() {
		atomic.AddInt32(&nowCatalogsCount, 1)
		progressTool.UpdateProgress(res.Url, int64(nowCatalogsCount))
	}()

	return &model.CrawlResult{
		OutputPath: outputPath,
		TakeTime:   int64(time.Since(startTime).Seconds()),
	}, nil
}
