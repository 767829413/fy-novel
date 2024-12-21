package model

// Chapter represents the structure of a book chapter
type Chapter struct {
	URL       string `json:"url"`
	ChapterNo int    `json:"chapterNo"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}
