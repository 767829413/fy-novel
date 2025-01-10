package parse

import (
	"bytes"
	"fmt"
	"strings"

	"fy-novel/internal/config"
	"fy-novel/internal/model"
	"fy-novel/internal/source"
	chapterTool "fy-novel/internal/tools/chapter"
	"fy-novel/pkg/utils"
	"github.com/gocolly/colly/v2"
)

type ChapterParser struct {
	rule model.Rule
	conf config.Info
}

func NewChapterParser(conf config.Info) *ChapterParser {
	return &ChapterParser{
		rule: source.GetRuleBySourceID(conf.Base.SourceID),
		conf: conf,
	}
}

func (b *ChapterParser) Parse(
	chapter *model.Chapter,
	res *model.SearchResult,
	book *model.Book,
	bookDir string,
) (err error) {
	// Prevent duplicate fetching
	chapter.Content, err = b.crawl(chapter.URL)
	if err != nil {
		// Attempt retry
		return err
	}
	err = chapterTool.ConvertChapter(chapter, b.conf.Base.Extname, b.rule)
	if err != nil {
		return err
	}
	return nil
}

func (b *ChapterParser) crawl(url string) (string, error) {
	nextUrl := url
	sb := bytes.NewBufferString("")

	for {
		collector := getCollector(nil, b.conf.Retry.MaxAttempts, b.conf.GetRandomDelay())
		collector.OnHTML(b.rule.Chapter.Content, func(e *colly.HTMLElement) {
			html, err := e.DOM.Html()
			if err == nil {
				sb.WriteString(html)
			} else {
				// Print error
				fmt.Printf("ChapterParser crawl Error parsing HTML: %v\n", err)
			}
		})
		if !b.rule.Chapter.Pagination {
			err := collector.Visit(nextUrl)
			if err != nil {
				return "", err
			}
			collector.Wait()
			return sb.String(), nil
		} else {
			collector.OnHTML(b.rule.Chapter.NextPage, func(e *colly.HTMLElement) {
				if strings.Contains(e.Text, "Next Chapter") {
					nextUrl = ""
					return
				}
				href := e.Attr("href")
				nextUrl = utils.NormalizeURL(href, b.rule.URL)
			})
			err := collector.Visit(nextUrl)
			if err != nil {
				return "", err
			}
			collector.Wait()
		}
		if nextUrl == "" {
			break
		}
	}
	return sb.String(), nil
}
