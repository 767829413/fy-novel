package chapter

import (
	"regexp"
	"strings"

	"fy-novel/internal/model"
	"github.com/PuerkitoBio/goquery"
)

func filterForChapter(chapter *model.Chapter, rule *model.Rule) string {
	return newFilterBuilder(chapter).
		FilterEscape(true).
		FilterAds(true).
		FilterDuplicateTitle(true).
		Build(rule)
}

type filterBuilder struct {
	title                     string
	content                   string
	applyEscapeFilter         bool
	applyAdsFilter            bool
	applyDuplicateTitleFilter bool
}

func newFilterBuilder(chapter *model.Chapter) *filterBuilder {
	return &filterBuilder{
		title:   chapter.Title,
		content: chapter.Content,
	}
}

func (fb *filterBuilder) FilterEscape(apply bool) *filterBuilder {
	fb.applyEscapeFilter = apply
	return fb
}

func (fb *filterBuilder) FilterAds(apply bool) *filterBuilder {
	fb.applyAdsFilter = apply
	return fb
}

func (fb *filterBuilder) FilterDuplicateTitle(apply bool) *filterBuilder {
	fb.applyDuplicateTitleFilter = apply
	return fb
}

func (fb *filterBuilder) Build(rule *model.Rule) string {
	if fb.applyEscapeFilter {
		re := regexp.MustCompile(`&[^;]+;`)
		fb.content = re.ReplaceAllString(fb.content, "")
	}

	if fb.applyAdsFilter {
		re := regexp.MustCompile(rule.Chapter.FilterTxt)
		filteredContent := re.ReplaceAllString(fb.content, "")

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(filteredContent))
		if err == nil {
			for _, tag := range strings.Fields(rule.Chapter.FilterTag) {
				doc.Find(tag).Remove()
			}
			fb.content, _ = doc.Html()
		}
	}

	fb.content = strings.TrimSpace(fb.content)

	if fb.applyDuplicateTitleFilter {
		title2 := strings.TrimSpace(fb.title)
		if strings.HasPrefix(fb.content, fb.title) || strings.HasPrefix(fb.content, title2) {
			re := regexp.MustCompile(
				`^(` + regexp.QuoteMeta(fb.title) + `|` + regexp.QuoteMeta(title2) + `)`,
			)
			fb.content = re.ReplaceAllString(fb.content, "")
		}
	}

	return strings.TrimSpace(fb.content)
}
