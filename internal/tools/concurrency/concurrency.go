package concurrency

import (
	"runtime"
	"sync"
)

const (
	maxThreads = 6 // Concurrency greater than this is prone to failure
)

var once sync.Once
var defaultConcurrency int

func GetConcurrencyNum(target, sourceID int) int {
	once.Do(func() {
		defaultConcurrency = runtime.NumCPU()
	})
	if target == -1 || target == 0 {
		target = defaultConcurrency * 2
	}
	if target >= maxThreads {
		target = maxThreads
	}
	// 书源3有防爬机制，不能并发太高
	if sourceID == 3 {
		return 1
	}
	return target
}
