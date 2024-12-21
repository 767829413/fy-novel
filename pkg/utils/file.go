package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// SortFilesByName 文件排序，按文件名升序
func SortFilesByName(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		name1 := files[i].Name()
		name2 := files[j].Name()

		no1, err1 := extractNumber(name1)
		no2, err2 := extractNumber(name2)

		if err1 != nil || err2 != nil {
			return name1 < name2 // 如果无法提取数字，则按字符串比较
		}

		return no1 < no2
	})

	return files, nil
}

// extractNumber 从文件名中提取数字
func extractNumber(name string) (int, error) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) < 2 {
		return 0, nil // 如果文件名不包含下划线，返回0
	}
	return strconv.Atoi(parts[0])
}

// GetSortedFilePaths 获取排序后的文件路径列表
func GetSortedFilePaths(dir string) ([]string, error) {
	files, err := SortFilesByName(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}
