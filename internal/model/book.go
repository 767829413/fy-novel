package model

// Book represents the structure of a book

type Book struct {
	URL           string `json:"url"`
	BookName      string `json:"bookName"`
	Author        string `json:"author"`
	Intro         string `json:"intro"`
	Category      string `json:"category"`
	CoverURL      string `json:"coverUrl"`
	LatestChapter string `json:"latestChapter"`
	LatestUpdate  string `json:"latestUpdate"`
	IsEnd         string `json:"isEnd"`
	Catalog       string `json:"catalog"`
}
