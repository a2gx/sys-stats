package storage

import (
	"errors"
	"sync"
	"time"
)

// Storage управляет несколькими именованными хранилищами одного типа.
type Storage[T any] struct {
	mu       sync.RWMutex
	stores   map[string]*Store[T]
	window   time.Duration
	cleanup  time.Duration
	capacity int
}

func NewStorage[T any](window, cleanup time.Duration, capacity int) *Storage[T] {
	return &Storage[T]{
		stores:   make(map[string]*Store[T]),
		window:   window,
		cleanup:  cleanup,
		capacity: capacity,
	}
}

// GetOrCreate возвращает существующее хранилище или создаёт новое с заданным именем.
func (m *Storage[T]) GetOrCreate(key string) *Store[T] {
	m.mu.Lock()
	defer m.mu.Unlock()

	if store, ok := m.stores[key]; ok {
		return store
	}

	store := NewStore[T](m.window, m.cleanup, m.capacity)
	m.stores[key] = store

	return store
}

func (m *Storage[T]) Add(key string, value T) {
	m.GetOrCreate(key).Add(value)
}

func (m *Storage[T]) GetSince(key string, duration time.Duration) ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	store, ok := m.stores[key]
	if !ok {
		return nil, errors.New("store not found: " + key)
	}

	return store.GetSince(duration), nil
}
