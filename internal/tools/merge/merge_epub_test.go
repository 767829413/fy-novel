package merge

import (
	"encoding/json"
	"testing"

	"fy-novel/internal/model"
)

var (
	bookJsonStr = ``
	dirPath     = ``
)

func TestEpubMergeHandler(t *testing.T) {
	book := model.Book{}
	json.Unmarshal([]byte(bookJsonStr), &book)
	_, err := epubMergeHandler(&book, dirPath)
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
}
