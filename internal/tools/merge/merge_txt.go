package merge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fy-novel/internal/model"
	"fy-novel/pkg/utils"
)

func txtMergeHandler(book *model.Book, dirPath string) (string, error) {
	outputPath := filepath.Join(
		filepath.Dir(dirPath),
		fmt.Sprintf("%s（%s）.txt", book.BookName, book.Author),
	)
	homePageFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("txtMergeHandler error creating output file: %v", err)
	}
	defer homePageFile.Close()
	// 首页添加书籍信息
	bookInfo := []string{
		fmt.Sprintf("书名：%s", book.BookName),
		fmt.Sprintf("作者：%s", book.Author),
		fmt.Sprintf("简介：%s", book.Intro),
		strings.Repeat("\u3000", 2),
	}

	for _, line := range bookInfo {
		if _, err := homePageFile.WriteString(line + "\n"); err != nil {
			return "", fmt.Errorf("txtMergeHandler error writing book info: %v", err)
		}
	}

	// 获取需要合并的章节文件
	filePaths, err := utils.GetSortedFilePaths(dirPath)
	if err != nil {
		return "", fmt.Errorf("txtMergeHandler error getting sorted file paths: %v", err)
	}

	for _, f := range filePaths {
		fh, err := os.Open(f)
		if err != nil {
			return "", fmt.Errorf("txtMergeHandler error reading file: %v", err)
		}
		_, err = io.Copy(homePageFile, fh)
		if err != nil {
			return "", fmt.Errorf("txtMergeHandler error copying from %s: %v", f, err)
		}
	}

	// 删除 txt 格式的临时文件
	err = os.RemoveAll(dirPath)
	if err != nil {
		return "", fmt.Errorf("epubMergeHandler Error removing temporary HTML files: %v", err)
	}
	return outputPath, nil
}
