package concurrency

import (
	"runtime"
	"sync"
)

var once sync.Once
var defaultConcurrency int

func GetConcurrencyNumBySourceID(target int) int {
	once.Do(func() {
		defaultConcurrency = runtime.NumCPU()
	})
	if target == -1 || target == 0 {
		target = defaultConcurrency * 2
	}
	return target
}
