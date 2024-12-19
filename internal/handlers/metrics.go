package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/sega-strn/go-metrics-tpl/internal/storage"
)

type MetricsHandler struct {
	storage *storage.MemStorage
}

func NewMetricsHandler(s *storage.MemStorage) *MetricsHandler {
	return &MetricsHandler{storage: s}
}

func (h *MetricsHandler) HandleUpdateMetric(w http.ResponseWriter, r *http.Request) {
	// Разбираем путь запроса
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid gauge value", http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(metricName, value)

	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid counter value", http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(metricName, value)

	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}