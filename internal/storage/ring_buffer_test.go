package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer_AddAndGetSince(t *testing.T) {
	buffer := newRingBuffer[int](3)

	now := time.Now()
	buffer.add(1, now.Add(-30*time.Second))
	buffer.add(2, now.Add(-20*time.Second))
	buffer.add(3, now.Add(-10*time.Second))

	// Проверяем, что все элементы добавлены
	assert.Equal(t, 3, buffer.size)

	// Получаем элементы, добавленные за последние 25 секунд
	result := buffer.getSince(now.Add(-25 * time.Second))
	assert.Equal(t, []int{2, 3}, result)
}

func TestRingBuffer_Overwrite(t *testing.T) {
	buffer := newRingBuffer[int](3)

	now := time.Now()
	buffer.add(1, now.Add(-4*time.Second))
	buffer.add(2, now.Add(-3*time.Second))
	buffer.add(3, now.Add(-2*time.Second))
	buffer.add(4, now.Add(-1*time.Second)) // Перезаписывает первый элемент

	// Проверяем, что размер буфера не превышает capacity
	assert.Equal(t, 3, buffer.size)

	// Проверяем, что элементы перезаписаны корректно
	result := buffer.getSince(now.Add(-5 * time.Second))
	assert.Equal(t, []int{2, 3, 4}, result)
}

func TestRingBuffer_RemoveOlderThan(t *testing.T) {
	buffer := newRingBuffer[int](3)

	now := time.Now()
	buffer.add(1, now.Add(-30*time.Second))
	buffer.add(2, now.Add(-20*time.Second))
	buffer.add(3, now.Add(-10*time.Second))

	// Удаляем элементы старше 22 секунд
	buffer.removeOlderThan(now.Add(-22 * time.Second))

	// Проверяем, что остались только новые элементы
	assert.Equal(t, 2, buffer.size)
	result := buffer.getSince(now.Add(-50 * time.Second))
	assert.Equal(t, []int{2, 3}, result)
}

func TestRingBuffer_Empty(t *testing.T) {
	buffer := newRingBuffer[int](3)

	now := time.Now()

	// Проверяем, что пустой буфер возвращает пустой результат
	result := buffer.getSince(now.Add(-1 * time.Second))
	assert.Empty(t, result)

	// Проверяем, что удаление из пустого буфера не вызывает ошибок
	buffer.removeOlderThan(now)
	assert.Equal(t, 0, buffer.size)
}
