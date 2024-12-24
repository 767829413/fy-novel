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
	RELEASE_URL = "https://api.github.com/repos/767829413/fy-novel/releases"
)

type CheckUpdater struct {
	log          *logrus.Logger
	timeoutMills int
}

func NewCheckUpdate(l *logrus.Logger, timeoutMills int) *CheckUpdater {
	return &CheckUpdater{log: l, timeoutMills: timeoutMills}
}

func (c *CheckUpdater) CheckUpdate() *model.GetUpdateInfoResult {
	c.log.Info("Starting check update")
	res := &model.GetUpdateInfoResult{
		CurrentVersion: version.Version,
	}
	// Implement update check logic
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
		latestVersion = latest["tag_name"].(string)
		latestUrl = latest["html_url"].(string)
		res.LatestVersion = latestVersion
		res.LatestUrl = latestUrl
		v1, err := semver.NewVersion(res.CurrentVersion)
		if err != nil {
			return
		}

		v2, err := semver.NewVersion(latestVersion)
		if err != nil {
			return
		}

		if v2.GreaterThan(v1) {
			res.NeedUpdate = true
		} else {
			res.NeedUpdate = false
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		res.ErrorMsg = fmt.Sprintf(
			"Check failed, current network environment cannot access GitHub, please try again later (%s)",
			err.Error(),
		)
	})

	collector.Visit(RELEASE_URL)
	collector.Wait()
	return res
}
