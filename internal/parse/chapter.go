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
}

func NewChapterParser(sourceID int) *ChapterParser {
	return &ChapterParser{
		rule: source.GetRuleBySourceID(sourceID),
	}
}

func (b *ChapterParser) Parse(
	chapter *model.Chapter,
	res *model.SearchResult,
	book *model.Book,
	bookDir string,
) (err error) {
	conf := config.GetConf()
	downOk := false
	attemptStart := 1
	attempt := conf.Retry.MaxAttempts
	if conf.Base.SourceID == 3 {
		attempt = 0
	}
	// Fetch content
	utils.SpinWaitMaxRetryAttempts(
		func() bool {
			var err error
			var errTemp = "==> Retrying to download failed chapter content: [%s], Attempt: %d/%d, Reason: %s\n"
			// Prevent duplicate fetching
			if !downOk {
				chapter.Content, err = b.crawl(chapter.URL, attempt)
				if err != nil {
					// Attempt retry
					fmt.Printf(
						errTemp,
						chapter.Title,
						attemptStart,
						attempt,
						err.Error(),
					)
					attemptStart++
					return false
				} else {
					downOk = true
				}
			}
			err = chapterTool.ConvertChapter(chapter, conf.Base.Extname, b.rule)
			if err != nil {
				// Attempt retry
				fmt.Printf(
					errTemp,
					chapter.Title,
					attemptStart,
					attempt,
					err.Error(),
				)
				attemptStart++
				return false
			}
			return true
		},
		attempt,
	)
	return nil
}

func (b *ChapterParser) crawl(url string, retry int) (string, error) {
	nextUrl := url
	sb := bytes.NewBufferString("")

	for {
		collector := getCollector(nil, retry)
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
