package parse

import (
	"net/http"
	"sync"

	"fy-novel/internal/model"
	"fy-novel/internal/source"
	"fy-novel/pkg/utils"

	"github.com/gocolly/colly/v2"
)

type SearchResultParser struct {
	rule model.Rule
}

func NewSearchResultParser(sourceID int) *SearchResultParser {
	return &SearchResultParser{
		rule: source.GetRuleBySourceID(sourceID),
	}
}

func (p *SearchResultParser) Parse(keyword string) ([]*model.SearchResult, error) {
	search := p.rule.Search
	isPaging := search.Pagination

	collector := getCollector(
		p.rule.Search.Cookies,
	)

	firstPageResults, err := p.getSearchResults(
		collector,
		p.rule.Search.URL,
		utils.BuildMethod(p.rule.Search.Method),
		keyword,
	)
	if err != nil {
		return nil, err
	}

	if len(firstPageResults) == 0 {
		return nil, nil
	}

	if !isPaging {
		return firstPageResults, nil
	}

	urls := make(map[string]struct{})

	collector.OnHTML(p.rule.Search.NextPage, func(e *colly.HTMLElement) {
		href := e.Attr("href")
		url := utils.NormalizeURL(href, p.rule.URL)
		urls[url] = struct{}{}
	})

	var wg sync.WaitGroup
	resultChan := make(chan []*model.SearchResult, len(urls))
	errorChan := make(chan error, len(urls))
	semaphore := make(chan struct{}, 20) // Limit concurrency to 20

	for url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			results, err := p.getSearchResults(nil, url, http.MethodGet, "")
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- results
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	var additionalResults []*model.SearchResult
	for results := range resultChan {
		additionalResults = append(additionalResults, results...)
	}

	// Check for errors
	for err := range errorChan {
		return nil, err
	}

	return append(firstPageResults, additionalResults...), nil
}

func (p *SearchResultParser) getSearchResults(
	collector *colly.Collector,
	url, method string,
	keyword string,
) ([]*model.SearchResult, error) {
	if collector == nil {
		collector = getCollector(
			p.rule.Search.Cookies,
		)
	}
	var results []*model.SearchResult
	collector.OnHTML(p.rule.Search.Result, func(e *colly.HTMLElement) {
		href := e.ChildAttr(p.rule.Search.BookName, "href")
		bookName := e.ChildText(p.rule.Search.BookName)
		latestChapter := e.ChildText(p.rule.Search.LatestChapter)
		author := e.ChildText(p.rule.Search.Author)
		update := e.ChildText(p.rule.Search.Update)

		if len(bookName) == 0 {
			return
		}

		result := &model.SearchResult{
			Url:           utils.NormalizeURL(href, p.rule.URL),
			BookName:      bookName,
			LatestChapter: latestChapter,
			Author:        author,
			LatestUpdate:  update,
		}

		results = append(results, result)
	})

	// collector.OnResponse(func(r *colly.Response) {
	// 	fmt.Println(string(r.Body))
	// })

	if method == http.MethodGet {
		err := collector.Visit(url)
		if err != nil {
			return nil, err
		}
	} else {
		bodyParams := utils.BuildParams(p.rule.Search.Body, keyword, "kw")
		err := collector.Post(url, bodyParams)
		if err != nil {
			return nil, err
		}
	}
	collector.Wait()
	return results, nil
}
