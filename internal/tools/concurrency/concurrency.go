package concurrency

import (
	"runtime"
	"sync"
)

const (
	maxThreads = 8 // 并发数大于这个容易失败
)

var once sync.Once
var defaultConcurrency int

func GetConcurrencyNum(target int) int {
	once.Do(func() {
		defaultConcurrency = runtime.NumCPU()
	})
	if target == -1 {
		target = defaultConcurrency * 2
	}
	if target >= maxThreads {
		target = maxThreads
	}
	return target
}
