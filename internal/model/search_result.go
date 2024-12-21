package model

type SearchResult struct {
	Url           string `json:"url"`
	BookName      string `json:"bookName"`
	Author        string `json:"author"`
	Intro         string `json:"intro"`
	LatestChapter string `json:"latestChapter"`
	LatestUpdate  string `json:"latestUpdate"`
}

type CrawlResult struct {
	OutputPath string
	TakeTime   int64
}
