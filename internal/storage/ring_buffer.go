// Кольцевой буфер для хранения временных меток и значений циклическая структура данных
// фиксированного размера, которая перезаписывает старые данные, когда буфер заполнен.
//
// https://medium.com/checker-engineering/a-practical-guide-to-implementing-a-generic-ring-buffer-in-go-866d27ec1a05

package storage

import (
	"time"
)

type entry[T any] struct {
	Value     T
	Timestamp time.Time
}

// ringBuffer реализует кольцевой буфер для хранения Entry[T].
type ringBuffer[T any] struct {
	data     []entry[T]
	head     int
	size     int
	capacity int
}

// newRingBuffer создает новый кольцевой буфер с заданной емкостью.
func newRingBuffer[T any](capacity int) *ringBuffer[T] {
	return &ringBuffer[T]{
		data:     make([]entry[T], capacity),
		capacity: capacity,
	}
}

// Add добавляет новое значение в кольцевой буфер с временной меткой.
func (r *ringBuffer[T]) add(value T, ts time.Time) {
	pos := (r.head + r.size) % r.capacity
	r.data[pos] = entry[T]{Value: value, Timestamp: ts}

	if r.size < r.capacity {
		r.size++
	} else {
		r.head = (r.head + 1) % r.capacity
	}
}

// getSince возвращает все значения, добавленные за последние duration.
func (r *ringBuffer[T]) getSince(cutoff time.Time) []T {
	var result []T

	for i := 0; i < r.size; i++ {
		idx := (r.head + i) % r.capacity
		e := r.data[idx]

		if e.Timestamp.After(cutoff) {
			result = append(result, e.Value)
		}
	}

	return result
}

// removeOlderThan удаляет старые записи по времени
func (r *ringBuffer[T]) removeOlderThan(cutoff time.Time) {
	for r.size > 0 {
		if r.data[r.head].Timestamp.After(cutoff) {
			break
		}

		r.head = (r.head + 1) % r.capacity
		r.size--
	}
}
