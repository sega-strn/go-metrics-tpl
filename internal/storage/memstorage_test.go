package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()

	// Тест обновления и получения gauge метрики
	t.Run("Gauge Metrics", func(t *testing.T) {
		storage.UpdateGauge("testGauge", 42.5)
		value, err := storage.GetGauge("testGauge")
		assert.NoError(t, err)
		assert.Equal(t, 42.5, value)
	})

	// Тест обновления и получения counter метрики
	t.Run("Counter Metrics", func(t *testing.T) {
		storage.UpdateCounter("testCounter", 10)
		storage.UpdateCounter("testCounter", 5)
		value, err := storage.GetCounter("testCounter")
		assert.NoError(t, err)
		assert.Equal(t, int64(15), value)
	})

	// Тест получения несуществующей метрики
	t.Run("Non-existent Metrics", func(t *testing.T) {
		_, err := storage.GetGauge("nonExistentGauge")
		assert.Error(t, err)

		_, err = storage.GetCounter("nonExistentCounter")
		assert.Error(t, err)
	})
}
