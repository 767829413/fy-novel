package chapter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fy-novel/internal/config"
	"fy-novel/internal/definition"
	"fy-novel/internal/model"
)

// createChapterFile saves the chapter content to a file
func CreateFileForChapter(chapter *model.Chapter, bookDir string) error {
	if chapter == nil {
		return nil
	}

	path, err := generatePath(chapter, bookDir)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, []byte(chapter.Content), 0644)
	if err != nil {
		return fmt.Errorf("Error writing file: %v\n", err)
	}
	return nil
}

// generatePath generates the file path for the chapter
func generatePath(chapter *model.Chapter, bookDir string) (string, error) {
	conf := config.GetConf()
	extName := conf.Base.Extname
	if extName == definition.NovelExtname_EPUB {
		extName = definition.NovelExtname_HTML
	}

	parentPath := filepath.Join(conf.Base.DownloadPath, bookDir)
	switch conf.Base.Extname {
	case definition.NovelExtname_HTML:
		return filepath.Join(parentPath, fmt.Sprintf("%d_.%s", chapter.ChapterNo, extName)), nil
	case definition.NovelExtname_EPUB, definition.NovelExtname_TXT:
		// Replace illegal characters in Windows file names
		title := strings.ReplaceAll(chapter.Title, "\\", "")
		title = strings.ReplaceAll(title, "/", "")
		title = strings.ReplaceAll(title, ":", "")
		title = strings.ReplaceAll(title, "*", "")
		title = strings.ReplaceAll(title, "?", "")
		title = strings.ReplaceAll(title, "<", "")
		title = strings.ReplaceAll(title, ">", "")
		return filepath.Join(
			parentPath,
			fmt.Sprintf("%d_%s.%s", chapter.ChapterNo, title, extName),
		), nil
	default:
		return "", fmt.Errorf("unsupported extension: %s", conf.Base.Extname)
	}
}
