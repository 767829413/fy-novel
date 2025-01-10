package parse

import (
	"fy-novel/internal/config"
	"fy-novel/internal/model"
	"fy-novel/internal/source"
	"fy-novel/pkg/utils"

	"github.com/gocolly/colly/v2"
	// "github.com/gocolly/colly/v2/debug"
)

type BookParser struct {
	rule model.Rule
	conf config.Info
}

func NewBookParser(conf config.Info) *BookParser {
	return &BookParser{
		rule: source.GetRuleBySourceID(conf.Base.SourceID),
		conf: conf,
	}
}

func (b *BookParser) Parse(bookUrl string) (*model.Book, error) {
	book := &model.Book{}
	collector := getCollector(nil, b.conf.Retry.MaxAttempts, b.conf.GetRandomDelay())
	// 抓取书名
	collector.OnHTML(b.rule.Book.BookName, func(e *colly.HTMLElement) {
		bookName := e.Attr("content")
		book.BookName = bookName
	})
	// 抓取作者
	collector.OnHTML(b.rule.Book.Author, func(e *colly.HTMLElement) {
		author := e.Attr("content")
		book.Author = author
	})
	// 抓取介绍
	collector.OnHTML(b.rule.Book.Intro, func(e *colly.HTMLElement) {
		intro := e.Attr("content")
		book.Intro = utils.CleanBlank(intro)
	})
	// 抓取封面url地址
	collector.OnHTML(b.rule.Book.CoverURL, func(e *colly.HTMLElement) {
		coverUrl := utils.NormalizeURL(e.Attr("src"), b.rule.URL)
		book.CoverURL = coverUrl
	})
	err := collector.Visit(bookUrl)
	if err != nil {
		return nil, err
	}
	collector.Wait()
	return book, nil
}
