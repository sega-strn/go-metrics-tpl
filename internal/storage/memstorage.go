package storage

import (
	"fmt"
	"sync"
)

// MetricType определяет тип метрики
type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

// Metric представляет собой интерфейс для работы с метриками
type Metric interface {
	Type() MetricType
	Value() interface{}
}

// GaugeMetric представляет метрику типа gauge
type GaugeMetric struct {
	value float64
}

func (g *GaugeMetric) Type() MetricType {
	return GaugeType
}

func (g *GaugeMetric) Value() interface{} {
	return g.value
}

func (g *GaugeMetric) Set(value float64) {
	g.value = value
}

// CounterMetric представляет метрику типа counter
type CounterMetric struct {
	value int64
}

func (c *CounterMetric) Type() MetricType {
	return CounterType
}

func (c *CounterMetric) Value() interface{} {
	return c.value
}

func (c *CounterMetric) Add(value int64) {
	c.value += value
}

// MemStorage представляет собой потокобезопасное хранилище метрик в памяти
type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]*GaugeMetric
	counters map[string]*CounterMetric
}

// NewMemStorage создает новое хранилище метрик в памяти
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]*GaugeMetric),
		counters: make(map[string]*CounterMetric),
	}
}

// UpdateGauge обновляет или устанавливает метрику типа gauge
func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.gauges[name]; !exists {
		ms.gauges[name] = &GaugeMetric{}
	}
	ms.gauges[name].Set(value)
}

// UpdateCounter обновляет или увеличивает метрику типа counter
func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.counters[name]; !exists {
		ms.counters[name] = &CounterMetric{}
	}
	ms.counters[name].Add(value)
}

// GetGauge возвращает метрику типа gauge
func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metric, ok := ms.gauges[name]
	if !ok {
		return 0, fmt.Errorf("gauge metric %s not found", name)
	}

	return metric.value, nil
}

// GetCounter возвращает метрику типа counter
func (ms *MemStorage) GetCounter(name string) (int64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metric, exists := ms.counters[name]
	if !exists {
		return 0, fmt.Errorf("counter metric %s not found", name)
	}
	return metric.value, nil
}

// RLock provides a read lock for thread-safe access to metrics
func (ms *MemStorage) RLock() {
	ms.mu.RLock()
}

// RUnlock releases the read lock
func (ms *MemStorage) RUnlock() {
	ms.mu.RUnlock()
}

// GetAllGauges returns a map of all gauge metrics
func (ms *MemStorage) GetAllGauges() map[string]float64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	gauges := make(map[string]float64)
	for name, metric := range ms.gauges {
		gauges[name] = metric.value
	}
	return gauges
}

// GetAllCounters returns a map of all counter metrics
func (ms *MemStorage) GetAllCounters() map[string]int64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	counters := make(map[string]int64)
	for name, metric := range ms.counters {
		counters[name] = metric.value
	}
	return counters
}
