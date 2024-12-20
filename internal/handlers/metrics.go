package handlers

import (
	"fmt"
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

func (h *MetricsHandler) HandleGetMetric(w http.ResponseWriter, r *http.Request) {
	// Разбираем путь запроса
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}

	metricType := parts[2]
	metricName := parts[3]

	switch metricType {
	case "gauge":
		value, err := h.storage.GetGauge(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%.3f", value)))

	case "counter":
		value, err := h.storage.GetCounter(metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%d", value)))

	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}
}

func (h *MetricsHandler) HandleListMetrics(w http.ResponseWriter, r *http.Request) {
	// Get all metrics from storage
	gauges := h.storage.GetAllGauges()
	counters := h.storage.GetAllCounters()

	// Create a response string
	var response strings.Builder
	response.WriteString("Metrics:\n\n")

	response.WriteString("Gauges:\n")
	for name, value := range gauges {
		response.WriteString(fmt.Sprintf("%s: %.3f\n", name, value))
	}

	response.WriteString("\nCounters:\n")
	for name, value := range counters {
		response.WriteString(fmt.Sprintf("%s: %d\n", name, value))
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.String()))
}
