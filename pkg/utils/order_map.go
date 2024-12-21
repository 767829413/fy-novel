package utils

import (
	"sync"
)

type orderedMap[K comparable, V any] struct {
	mu   sync.RWMutex
	keys []K
	data map[K]V
}

func NewOrderedMap[K comparable, V any]() *orderedMap[K, V] {
	return &orderedMap[K, V]{
		keys: make([]K, 0),
		data: make(map[K]V),
	}
}

func (om *orderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, exists := om.data[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.data[key] = value
}

func (om *orderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	val, exists := om.data[key]
	return val, exists
}

func (om *orderedMap[K, V]) Delete(key K) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, exists := om.data[key]; exists {
		delete(om.data, key)
		for i, k := range om.keys {
			if k == key {
				om.keys = append(om.keys[:i], om.keys[i+1:]...)
				break
			}
		}
	}
}

func (om *orderedMap[K, V]) Keys() []K {
	om.mu.RLock()
	defer om.mu.RUnlock()

	keys := make([]K, len(om.keys))
	copy(keys, om.keys)
	return keys
}

func (om *orderedMap[K, V]) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return len(om.keys)
}

func (om *orderedMap[K, V]) Range(f func(key K, value V) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for _, key := range om.keys {
		if !f(key, om.data[key]) {
			break
		}
	}
}
