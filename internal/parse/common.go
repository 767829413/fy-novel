package parse

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	// "github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

const timeoutMillis = 15000

func getCollector(cookies map[string]string) *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true),
		// Attach a debugger to the collector
		// colly.Debugger(&debug.LogDebugger{}),
	)
	extensions.RandomUserAgent(c)
	c.SetRequestTimeout(timeoutMillis * time.Millisecond)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		//Delay:      5 * time.Second,
	})

	// 設定錯誤回呼函式
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("/nError: %s: Request URL: %s/n", err, r.Request.URL)
	})

	// Set cookies
	if len(cookies) > 0 {
		for k, v := range cookies {
			c.OnRequest(func(r *colly.Request) {
				r.Headers.Set("Cookie", fmt.Sprintf("%s=%s", k, v))
			})
		}
	}
	return c
}
