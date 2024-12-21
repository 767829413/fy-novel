package parse

import (
	"sort"

	"fy-novel/internal/model"
	"fy-novel/internal/source"
	"fy-novel/pkg/utils"
	"github.com/gocolly/colly/v2"
	// "github.com/gocolly/colly/v2/debug"
)

type CatalogsParser struct {
	rule *model.Rule
}

func NewCatalogsParser(sourceID int) *CatalogsParser {
	return &CatalogsParser{
		rule: source.GetRuleBySourceID(sourceID),
	}
}

func (b *CatalogsParser) Parse(bookUrl string, start, end int) ([]*model.Chapter, error) {
	collector := getCollector(nil)

	var chapters = make(map[string]*model.Chapter)

	collector.OnHTML(b.rule.Book.Catalog, func(e *colly.HTMLElement) {
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
