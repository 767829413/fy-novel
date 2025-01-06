package progress

import (
	"sync"
)

var (
	progressTracker *ProgressTracker
	once            sync.Once
)

type ProgressTracker struct {
	tasks map[string]*TaskProgress
	mu    sync.RWMutex
}

type TaskProgress struct {
	Total     int64
	Completed int64
	mu        sync.RWMutex
}

func init() {
	once.Do(func() {
		progressTracker = &ProgressTracker{
			tasks: make(map[string]*TaskProgress),
		}
	})
}

// InitTask Initialize a new task progress
func InitTask(taskID string, total int64) {
	progressTracker.mu.Lock()
	defer progressTracker.mu.Unlock()
	progressTracker.tasks[taskID] = &TaskProgress{Total: total}
}

// UpdateProgress Updating the progress of a given task
func UpdateProgress(taskID string, completed int64) {
	progressTracker.mu.RLock()
	task, exists := progressTracker.tasks[taskID]
	progressTracker.mu.RUnlock()

	if exists {
		task.mu.Lock()
		task.Completed += completed
		task.mu.Unlock()
	}
}

// GetProgress Get the progress of a given task
func GetProgress(taskID string) (int64, int64, bool) {
	progressTracker.mu.RLock()
	task, exists := progressTracker.tasks[taskID]
	progressTracker.mu.RUnlock()

	if !exists {
		return 0, 0, false
	}

	task.mu.RLock()
	completed := task.Completed
	total := task.Total
	task.mu.RUnlock()

	return completed, total, true
}

func DeleteTask(taskID string) {
	progressTracker.mu.Lock()
	defer progressTracker.mu.Unlock()
	delete(progressTracker.tasks, taskID)
}
