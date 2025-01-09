package parse

import (
	"fmt"
	"sort"

	"fy-novel/internal/model"
	"fy-novel/internal/source"
	"fy-novel/pkg/utils"

	"github.com/gocolly/colly/v2"
	// "github.com/gocolly/colly/v2/debug"
)

type CatalogsParser struct {
	rule model.Rule
}

func NewCatalogsParser(sourceID int) *CatalogsParser {
	return &CatalogsParser{
		rule: source.GetRuleBySourceID(sourceID),
	}
}

func (b *CatalogsParser) Parse(bookUrl string, start, end, retry int) ([]*model.Chapter, error) {
	collector := getCollector(nil, retry)

	var chapters = make(map[string]*model.Chapter)

	if len(b.rule.Catalog.URL) > 0 {
		id := utils.GetGroup1(b.rule.Book.URL, bookUrl)
		bookUrl = fmt.Sprintf(b.rule.Catalog.URL, id)
	}

	collector.OnHTML(b.rule.Catalog.Result, func(e *colly.HTMLElement) {
		chapter := &model.Chapter{
			Title: e.Text,
			URL:   utils.NormalizeURL(e.Attr("href"), b.rule.URL),
		}
		chapter.ChapterNo = len(chapters) + 1
		chapters[chapter.Title] = chapter
	})

	err := collector.Visit(bookUrl)
	if err != nil {
		return nil, err
	}

	collector.Wait()

	res := make([]*model.Chapter, 0, len(chapters))
	for _, chapter := range chapters {
		res = append(res, chapter)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ChapterNo < res[j].ChapterNo
	})
	return res, nil
}
