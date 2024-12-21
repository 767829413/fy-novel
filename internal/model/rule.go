package model

// Rule represents the main structure for rules
type Rule struct {
	ID      string  `json:"id"`
	URL     string  `json:"url"`
	Name    string  `json:"name"`
	Comment string  `json:"comment"`
	Type    string  `json:"type"`
	Search  search  `json:"search"`
	Book    book    `json:"book"`
	Chapter chapter `json:"chapter"`
}

// Search represents the search rules
type search struct {
	URL           string            `json:"url"`
	Method        string            `json:"method"`
	Param         string            `json:"param"`
	Body          map[string]string `json:"body"`
	Cookies       map[string]string `json:"cookies"`
	Pagination    bool              `json:"pagination"`
	NextPage      string            `json:"nextPage"`
	Result        string            `json:"result"`
	BookName      string            `json:"bookName"`
	LatestChapter string            `json:"latestChapter"`
	Author        string            `json:"author"`
	Update        string            `json:"update"`
}

// Book represents the book rules
type book struct {
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
	CatalogOffset int    `json:"catalogOffset"`
}

// Chapter represents the chapter rules
type chapter struct {
	URL                string `json:"url"`
	Pagination         bool   `json:"pagination"`
	NextPage           string `json:"nextPage"`
	ChapterNo          int    `json:"chapterNo"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	ParagraphTagClosed bool   `json:"paragraphTagClosed"`
	ParagraphTag       string `json:"paragraphTag"`
	FilterTxt          string `json:"filterTxt"`
	FilterTag          string `json:"filterTag"`
}
