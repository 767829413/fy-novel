package functions

import (
	"fmt"

	"fy-novel/internal/config"
	"fy-novel/internal/source"
	"fy-novel/internal/version"
	"github.com/sirupsen/logrus"
)

type GetConf struct {
	log *logrus.Logger
}

func NewGetConf(l *logrus.Logger) *GetConf {
	return &GetConf{log: l}
}

func (p *GetConf) GetConfigString() (string, error) {
	conf, err := config.GetConf().ToJSON()
	return conf, err
}

type GetHint struct {
	log *logrus.Logger
}

func NewGetHint(l *logrus.Logger) *GetHint {
	return &GetHint{log: l}
}

func (p *GetHint) GetUsageInfo() []string {
	cfg := config.GetConf()
	rule := source.GetRuleBySourceID(cfg.Base.SourceID)
	res := []string{
		"使用须知",
		fmt.Sprintf(
			"easy-novel %s (commit %s, built at %s)\n",
			version.Version,
			version.Commit,
			version.Date,
		),
		fmt.Sprintf("官方地址：%s", "https://fy-novel"),
		fmt.Sprintf("当前书源：%s (ID: %d)", rule.URL, cfg.Base.SourceID),
		fmt.Sprintf("导出格式：%s", cfg.Base.Extname),
		"请务必阅读 readme.txt",
	}
	return res
}
