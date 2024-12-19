package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mu             sync.RWMutex
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

func (m *Metrics) UpdateGaugeMetric(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gaugeMetrics[name] = value
}

func (m *Metrics) UpdateCounterMetric(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counterMetrics[name] += value
}

func (m *Metrics) CollectRuntimeMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.UpdateGaugeMetric("Alloc", float64(rtm.Alloc))
	m.UpdateGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys))
	m.UpdateGaugeMetric("Frees", float64(rtm.Frees))
	m.UpdateGaugeMetric("GCCPUFraction", rtm.GCCPUFraction)
	m.UpdateGaugeMetric("GCSys", float64(rtm.GCSys))
	m.UpdateGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc))
	m.UpdateGaugeMetric("HeapIdle", float64(rtm.HeapIdle))
	m.UpdateGaugeMetric("HeapInuse", float64(rtm.HeapInuse))
	m.UpdateGaugeMetric("HeapObjects", float64(rtm.HeapObjects))
	m.UpdateGaugeMetric("HeapReleased", float64(rtm.HeapReleased))
	m.UpdateGaugeMetric("HeapSys", float64(rtm.HeapSys))
	m.UpdateGaugeMetric("LastGC", float64(rtm.LastGC))
	m.UpdateGaugeMetric("Lookups", float64(rtm.Lookups))
	m.UpdateGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse))
	m.UpdateGaugeMetric("MCacheSys", float64(rtm.MCacheSys))
	m.UpdateGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse))
	m.UpdateGaugeMetric("MSpanSys", float64(rtm.MSpanSys))
	m.UpdateGaugeMetric("Mallocs", float64(rtm.Mallocs))
	m.UpdateGaugeMetric("NextGC", float64(rtm.NextGC))
	m.UpdateGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC))
	m.UpdateGaugeMetric("NumGC", float64(rtm.NumGC))
	m.UpdateGaugeMetric("OtherSys", float64(rtm.OtherSys))
	m.UpdateGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs))
	m.UpdateGaugeMetric("StackInuse", float64(rtm.StackInuse))
	m.UpdateGaugeMetric("StackSys", float64(rtm.StackSys))
	m.UpdateGaugeMetric("Sys", float64(rtm.Sys))
	m.UpdateGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc))

	// Дополнительные метрики
	m.UpdateGaugeMetric("RandomValue", rand.Float64())
	m.UpdateCounterMetric("PollCount", 1)
}

func (m *Metrics) SendMetrics() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, value := range m.gaugeMetrics {
		err := sendMetric("gauge", name, fmt.Sprintf("%f", value))
		if err != nil {
			log.Printf("Error sending gauge metric %s: %v", name, err)
			return err
		}
	}

	for name, value := range m.counterMetrics {
		err := sendMetric("counter", name, fmt.Sprintf("%d", value))
		if err != nil {
			log.Printf("Error sending counter metric %s: %v", name, err)
			return err
		}
	}

	return nil
}

var sendMetricFunc = func(metricType, name, value string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, name, value)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func sendMetric(metricType, name, value string) error {
	return sendMetricFunc(metricType, name, value)
}

func main() {
	metrics := NewMetrics()

	pollTicker := time.NewTicker(2 * time.Second)
	reportTicker := time.NewTicker(10 * time.Second)

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			metrics.CollectRuntimeMetrics()
		case <-reportTicker.C:
			if err := metrics.SendMetrics(); err != nil {
				log.Printf("Error sending metrics: %v", err)
			}
		}
	}
}
