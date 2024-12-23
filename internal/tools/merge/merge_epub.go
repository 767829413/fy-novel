package merge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fy-novel/internal/model"
	"fy-novel/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/go-shiori/go-epub"
)

func epubMergeHandler(book *model.Book, dirPath string) (string, error) {
	var (
		err       error
		filePaths []string
		attempts  = 7
	)
	utils.SpinWaitMaxRetryAttempts(func() bool {
		filePaths, err = utils.GetSortedFilePaths(dirPath)
		if len(filePaths) == 0 || err != nil {
			return false
		}
		return true
	}, attempts)
	if len(filePaths) == 0 {
		return "", fmt.Errorf("epubMergeHandler error getting sorted file paths: %v", err)
	}

	// 等待文件系统更新索引
	epubIns, err := epub.NewEpub(book.BookName)
	if err != nil {
		return "", fmt.Errorf("epubMergeHandler error creating epub instance: %v", err)
	}
	epubIns.SetAuthor(book.Author)
	epubIns.SetDescription(book.Intro)
	epubIns.SetLang("zh")

	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("epubMergeHandler error reading file: %v", err)
		}
		// 获取文件名
		fileName := filepath.Base(filePath)

		// 从文件名中提取标题
		title := strings.SplitN(fileName, "_", 2)[1]
		title = strings.TrimSuffix(title, filepath.Ext(title))
		_, err = epubIns.AddSection(string(content), title, "", "")
		if err != nil {
			return "", fmt.Errorf("epubMergeHandler error adding section: %v", err)
		}
	}

	// 下载封面
	if false && len(book.CoverURL) != 0 {
		client := resty.New()
		resp, err := client.R().Get(book.CoverURL)
		if err == nil {
			coverPath := filepath.Join(dirPath, "cover.jpg")
			if err := os.WriteFile(coverPath, resp.Body(), 0644); err == nil {
				epubIns.AddImage(coverPath, "cover.jpg")
				epubIns.SetCover("cover.jpg", "")
			}
		}
	}
	// 保存 EPUB 文件
	savePath := filepath.Join(filepath.Dir(dirPath), book.BookName+".epub")
	err = epubIns.Write(savePath)
	if err != nil {
		return "", fmt.Errorf("epubMergeHandler error writing EPUB file: %v", err)
	}
	// 处理 EPUB 格式的临时 HTML , 删除文件
	err = os.RemoveAll(dirPath)
	if err != nil {
		return "", fmt.Errorf("epubMergeHandler Error removing temporary HTML files: %v", err)
	}
	return savePath, nil
}
