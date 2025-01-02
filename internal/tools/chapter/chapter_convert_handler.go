package chapter

import (
	"fy-novel/internal/definition"
	"fy-novel/internal/model"
)

func ConvertChapter(
	chapter *model.Chapter,
	extName string,
	rule model.Rule,
) error {
	var content string
	var err error
	content = formatForChapter(filterForChapter(chapter, rule), rule)

	switch extName {
	case definition.NovelExtname_TXT:
		content = txtConvert(chapter.Title, content)
	case definition.NovelExtname_EPUB, definition.NovelExtname_HTML:
		content, err = templateConvert(chapter.Title, content, extName)
		if err != nil {
			return err
		}
	}
	chapter.Content = content
	return nil
}
