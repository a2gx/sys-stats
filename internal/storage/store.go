package storage

import (
	"sync"
	"time"
)

// Store хранилище типа in-memory с ограничением по времени хранения.
type Store[T any] struct {
	mu      sync.RWMutex
	buffer  *ringBuffer[T]
	window  time.Duration // Время, в течение которого записи считаются актуальными
	cleanup time.Duration // Время, через которое старые записи удаляются
}

// NewStore создаёт новое хранилище с тайм-аутом и автоматической очисткой.
func NewStore[T any](window, cleanup time.Duration, capacity int) *Store[T] {
	s := &Store[T]{
		buffer:  newRingBuffer[T](capacity),
		window:  window,
		cleanup: cleanup,
	}
	go s.cleanupLoop()
	return s
}

// Add сохраняет новое значение с временной меткой.
func (s *Store[T]) Add(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.buffer.add(value, time.Now())
}

// GetSince возвращает все значения, добавленные за последние duration.
func (s *Store[T]) GetSince(duration time.Duration) []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.buffer.getSince(time.Now().Add(-duration))
}

// cleanupOld удаляем старые записи по таймеру
// private
func (s *Store[T]) cleanupLoop() {
	ticker := time.NewTicker(s.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.buffer.removeOlderThan(time.Now().Add(-s.window))
		}()
	}
}
