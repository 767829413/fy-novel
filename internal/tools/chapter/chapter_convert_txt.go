package chapter

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func txtConvert(title, content string) string {
	// 全角空格，用于首行缩进
	indent := strings.Repeat("\u3000", 2)

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return fmt.Sprintf("Error parsing HTML: %v", err)
	}

	var result strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				result.WriteString(indent)
				result.WriteString(html.UnescapeString(text))
				result.WriteString("\n")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return fmt.Sprintf("%s\n\n%s", title, result.String())
}
