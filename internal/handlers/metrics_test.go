package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sega-strn/go-metrics-tpl/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandler(t *testing.T) {
	memStorage := storage.NewMemStorage()
	handler := NewMetricsHandler(memStorage)

	t.Run("Update Gauge Metric", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/update/gauge/testGauge/42.5", nil)
		w := httptest.NewRecorder()
		handler.HandleUpdateMetric(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		value, err := memStorage.GetGauge("testGauge")
		assert.NoError(t, err)
		assert.Equal(t, 42.5, value)
	})

	t.Run("Update Counter Metric", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/update/counter/testCounter/10", nil)
		w := httptest.NewRecorder()
		handler.HandleUpdateMetric(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		value, err := memStorage.GetCounter("testCounter")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), value)
	})

	t.Run("Invalid Metric Type", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/update/invalid/testMetric/42", nil)
		w := httptest.NewRecorder()
		handler.HandleUpdateMetric(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
