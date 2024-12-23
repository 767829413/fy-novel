package model

import "fy-novel/internal/config"

type GetUsageInfoResult struct {
	Title             string
	VersionInfo       string
	Address           string
	CurrentBookSource string
	ExportFormat      string
}

type GetUpdateInfoResult struct {
	UpdateInfo string
	LatestUrl  string
}

type GetConfigResult struct {
	Config config.Info
}

type ProgressResult struct {
	Exists    bool
	Completed int
	Total     int
}
