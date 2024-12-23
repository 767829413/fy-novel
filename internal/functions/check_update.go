package functions

import (
	"encoding/json"
	"fmt"
	"time"

	"fy-novel/internal/model"
	"fy-novel/internal/version"

	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/sirupsen/logrus"
)

const (
	RELEASE_URL = "https://api.github.com/repos/767829413/easy-novel/releases"
)

type CheckUpdater struct {
	log          *logrus.Logger
	timeoutMills int
}

func NewCheckUpdate(l *logrus.Logger, timeoutMills int) *CheckUpdater {
	return &CheckUpdater{log: l, timeoutMills: timeoutMills}
}

func (c *CheckUpdater) CheckUpdate() *model.GetUpdateInfoResult {
	var res = "无法获取版本信息"
	c.log.Info("Starting check update")
	// 实现检查更新的逻辑
	collector := colly.NewCollector(
		colly.Async(true),
	)
	extensions.RandomUserAgent(collector)
	collector.SetRequestTimeout(time.Duration(c.timeoutMills) * time.Millisecond)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
	})
	var latestVersion, latestUrl string
	collector.OnResponse(func(r *colly.Response) {
		var releases []map[string]interface{}
		err := json.Unmarshal(r.Body, &releases)
		if err != nil || len(releases) == 0 {
			return
		}

		latest := releases[0]
		currentVersion := version.Version
		latestVersion = latest["tag_name"].(string)
		latestUrl = latest["html_url"].(string)

		v1, err := semver.NewVersion(currentVersion)
		if err != nil {
			return
		}

		v2, err := semver.NewVersion(latestVersion)
		if err != nil {
			return
		}

		if v2.GreaterThan(v1) {
			res = fmt.Sprintf("发现新版本: %s, 当前版本: %s", latestVersion, currentVersion)
		} else {
			res = fmt.Sprintf("%s 已是最新版本！", latestVersion)
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		res = fmt.Sprintf("检查失败, 当前网络环境暂时无法访问 GitHub, 请稍后再试 (%s)", err.Error())
	})

	collector.Visit(RELEASE_URL)
	collector.Wait()
	return &model.GetUpdateInfoResult{
		UpdateInfo: res,
		LatestUrl:  latestUrl,
	}
}
