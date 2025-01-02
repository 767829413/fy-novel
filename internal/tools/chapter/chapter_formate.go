package chapter

import (
	"regexp"
	"strings"

	"fy-novel/internal/model"
)

func formatForChapter(content string, rule model.Rule) string {
	if rule.Chapter.ParagraphTagClosed {
		if rule.Chapter.ParagraphTag == "p" {
			return content
		} else {
			// 非 <p> 闭合标签, 替换为 <p>
			re := regexp.MustCompile(`<([a-zA-Z][^>]*)>(.*?)</([a-zA-Z][^>]*)>`)
			return re.ReplaceAllStringFunc(content, func(match string) string {
				if strings.HasPrefix(match, "<p>") {
					return match
				}
				submatches := re.FindStringSubmatch(match)
				if len(submatches) == 4 && submatches[1] == submatches[3] {
					return "<p>" + submatches[2] + "</p>"
				}
				return match
			})
		}
	}

	// 标签不闭合, 用某个标签分隔, 例如：段落1<br><br>段落2
	tag := rule.Chapter.ParagraphTag
	parts := strings.Split(content, tag)
	var sb strings.Builder

	for _, s := range parts {
		if strings.TrimSpace(s) != "" {
			sb.WriteString("<p>")
			sb.WriteString(s)
			sb.WriteString("</p>")
		}
	}

	return sb.String()
}
