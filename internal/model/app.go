package model

import (
	"fy-novel/internal/config"
)

type GetUsageInfoResult struct {
	VersionInfo       string
	Address           string
	CurrentBookSource string
	ExportFormat      string
}

type GetUpdateInfoResult struct {
	ErrorMsg       string
	NeedUpdate     bool
	LatestVersion  string
	CurrentVersion string
	LatestUrl      string
}

type GetConfigResult struct {
	Config config.Info
}

type ProgressResult struct {
	Exists    bool
	Completed int
	Total     int
}

type HasInitOllamaResult struct {
	Has        bool
	IsInit     bool
	IsSetModel bool
	ErrorMsg   string
}

type InitOllamaResult struct {
	ErrorMsg string
}

type InitOllamaProgressResult struct {
	Exists    bool
	Completed int
	Total     int
}

type InitSetOllamaModelResult struct {
	ErrorMsg string
}

type SetOllamaModelResult struct {
	ErrorMsg string
}

type GetCurrentUseModelResult struct {
	Model string
}

type GetSelectModelListResult struct {
	Models []string
}

type GetSetOllamaModelProgressResult struct {
	Exists    bool
	Completed int
	Total     int
}

type StartChatbotResult struct {
	Response string
	ErrorMsg string
}

type YukkuriParams struct {
	ImgPath    string // specify an image file to transcode
	Threshold  int    // set threshold for image grey processing (default 140)
	Filename   string // set output file name
	TmpImgName string // save temporary image
	AscWidth   int    // set ascii max width
	AscHeight  int    // set ascii max height
	Ykr        bool   // true -> 使用AsciiYukkuri处理 ,false -> 使用Ascii
}

type GenerateAsciiImageResult struct {
	Response string
	ErrorMsg string
}
