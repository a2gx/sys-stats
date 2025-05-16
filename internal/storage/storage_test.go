package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorage_GetOrCreate(t *testing.T) {
	storage := NewStorage[int](10*time.Second, 5*time.Second, 5)

	// Создаём новое хранилище
	store1 := storage.GetOrCreate("test")
	assert.NotNil(t, store1)

	// Получаем уже существующее хранилище
	store2 := storage.GetOrCreate("test")
	assert.Equal(t, store1, store2)
}

func TestStorage_AddAndGetSince(t *testing.T) {
	storage := NewStorage[int](10*time.Second, 5*time.Second, 5)

	// Добавляем значения в хранилище
	storage.Add("test", 1)
	time.Sleep(1 * time.Second)
	storage.Add("test", 2)
	time.Sleep(1 * time.Second)
	storage.Add("test", 3)

	// Получаем значения за последние 2 секунды
	result, err := storage.GetSince("test", 2*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 3}, result)

	// Получаем все значения за последние 10 секунд
	result, err = storage.GetSince("test", 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestStorage_GetSince_StoreNotFound(t *testing.T) {
	storage := NewStorage[int](10*time.Second, 5*time.Second, 5)

	// Пытаемся получить данные из несуществующего хранилища
	result, err := storage.GetSince("nonexistent", 5*time.Second)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "store not found: nonexistent", err.Error())
}
