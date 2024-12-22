package model

type GetConfigResult struct {
	Config string
	Error  string
}

type ProgressResult struct {
	Exists    bool
	Completed int
	Total     int
}
