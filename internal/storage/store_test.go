package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore_AddAndGetSince(t *testing.T) {
	store := NewStore[int](10*time.Second, 5*time.Second, 5)

	//now := time.Now()
	store.Add(1)
	time.Sleep(1 * time.Second)
	store.Add(2)
	time.Sleep(1 * time.Second)
	store.Add(3)

	// Get values added in the last 2 seconds
	result := store.GetSince(2 * time.Second)
	assert.Equal(t, []int{2, 3}, result)

	// Get all values added in the last 10 seconds
	result = store.GetSince(10 * time.Second)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStore_CleanupOldEntries(t *testing.T) {
	store := NewStore[int](3*time.Second, 1*time.Second, 5)

	//now := time.Now()
	store.Add(1)
	store.Add(2)
	time.Sleep(2 * time.Second)
	store.Add(3)
	time.Sleep(2 * time.Second)

	// Только последняя запись должна остаться после очистки
	result := store.GetSince(10 * time.Second)
	assert.Equal(t, []int{3}, result)
}

func TestStore_Empty(t *testing.T) {
	store := NewStore[int](10*time.Second, 5*time.Second, 5)

	// Вернём пустое значение
	result := store.GetSince(1 * time.Second)
	assert.Equal(t, 0, len(result))

	// Добавим значение и проверим, что оно есть
	store.Add(1)
	result = store.GetSince(10 * time.Second)
	assert.Equal(t, []int{1}, result)
}
