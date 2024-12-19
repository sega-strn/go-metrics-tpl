package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsCollection(t *testing.T) {
	metrics := NewMetrics()

	t.Run("Collect Runtime Metrics", func(t *testing.T) {
		initialPollCount := metrics.counterMetrics["PollCount"]
		metrics.CollectRuntimeMetrics()

		assert.Greater(t, metrics.counterMetrics["PollCount"], initialPollCount)
		assert.NotZero(t, metrics.gaugeMetrics["RandomValue"])
	})
}

func TestSendMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Временно заменим функцию отправки метрик
	originalSendMetricFunc := sendMetricFunc
	defer func() { sendMetricFunc = originalSendMetricFunc }()

	sendMetricFunc = func(metricType, name, value string) error {
		return nil
	}

	metrics := NewMetrics()
	metrics.UpdateGaugeMetric("TestGauge", 42.5)
	metrics.UpdateCounterMetric("TestCounter", 10)

	t.Run("Send Metrics", func(t *testing.T) {
		err := metrics.SendMetrics()
		assert.NoError(t, err)
	})
}
