package progress

import (
	"sync"
)

var progressMap sync.Map

// Progress 结构体用于跟踪任务进度
type progress struct {
	total   int64
	current int64
	mu      sync.RWMutex
}

func InitProgress(uniqueKe string, total int64) {
	progressMap.Store(uniqueKe, &progress{total: total})
}

func SetProgress(uniqueKe string, current int64) {
	if v, ok := progressMap.Load(uniqueKe); ok {
		p := v.(*progress)
		p.mu.Lock()
		p.current = current
		p.mu.Unlock()
	}
}

func GetProgress(uniqueKe string) float64 {
	if v, ok := progressMap.Load(uniqueKe); ok {
		return float64(v.(*progress).total) / float64(v.(*progress).total)
	}
	return 0
}
