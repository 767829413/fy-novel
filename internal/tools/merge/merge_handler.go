package merge

import (
	"fmt"

	"fy-novel/internal/config"
	"fy-novel/internal/definition"
	"fy-novel/internal/model"
)

func MergeSaveHandler(book *model.Book, dirPath string) (string, error) {
	conf := config.GetConf()
	switch conf.Base.Extname {
	case definition.NovelExtname_TXT:
		return txtMergeHandler(book, dirPath)
	case definition.NovelExtname_EPUB:
		return epubMergeHandler(book, dirPath)
	default:
		return "", fmt.Errorf("unsupported extension: %s", conf.Base.Extname)
	}
}
