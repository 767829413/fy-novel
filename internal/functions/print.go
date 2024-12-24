package functions

import (
	"fmt"

	"fy-novel/internal/config"
	"fy-novel/internal/model"
	"fy-novel/internal/source"
	"fy-novel/internal/version"

	"github.com/sirupsen/logrus"
)

type ConfHandler struct {
	log *logrus.Logger
}

func NewGetConf(l *logrus.Logger) *ConfHandler {
	return &ConfHandler{log: l}
}

func (p *ConfHandler) GetConfig() config.Info {
	return config.GetConf()
}

func (p *ConfHandler) SetConfig(conf string) error {
	return config.SetConf(conf)
}

type GetHint struct {
	log *logrus.Logger
}

func NewGetHint(l *logrus.Logger) *GetHint {
	return &GetHint{log: l}
}

func (p *GetHint) GetUsageInfo() *model.GetUsageInfoResult {
	cfg := config.GetConf()
	rule := source.GetRuleBySourceID(cfg.Base.SourceID)
	res := &model.GetUsageInfoResult{
		Title: "fy-novel 使用须知",
		VersionInfo: fmt.Sprintf(
			"fy-novel %s (commit %s, built at %s)\n",
			version.Version,
			version.Commit,
			version.Date,
		),
		Address:           "https://github.com/767829413/fy-novel",
		CurrentBookSource: fmt.Sprintf("%s (ID: %d)", rule.URL, cfg.Base.SourceID),
		ExportFormat:      cfg.Base.Extname,
	}

	return res
}
