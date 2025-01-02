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

func GetConcurrencyNum(target int) int {
	once.Do(func() {
		defaultConcurrency = runtime.NumCPU()
	})
	if target == -1 || target == 0 {
		target = defaultConcurrency * 2
	}
	if target >= maxThreads {
		target = maxThreads
	}
	return target
}
