package parse

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	// "github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
)

const timeoutMillis = 25000
const retryDefault = 3

var urlLock sync.Mutex

var saveErrorUrl = make(map[string]int)

func getCollector(cookies map[string]string, retry int) *colly.Collector {
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

	if retry == 0 {
		retry = retryDefault
	}

	// 设置错误重试
	c.OnError(func(r *colly.Response, err error) {
		// 加入一个自动重试机制
		link := r.Request.URL.String()
		urlLock.Lock()
		if _, ok := saveErrorUrl[link]; !ok {
			saveErrorUrl[link]++
			r.Request.Retry()
		} else if saveErrorUrl[link] < retry {
			saveErrorUrl[link]++
			r.Request.Retry()
		} else {
			fmt.Printf("\nRetry %d Request URL: %s, Error: %v", saveErrorUrl[link], link, err)
		}
		urlLock.Unlock()
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
