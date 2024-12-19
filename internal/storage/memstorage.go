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

// MemStorage represents in-memory storage for metrics
type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]*GaugeMetric
	counters map[string]*CounterMetric
}

// NewMemStorage creates a new instance of MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]*GaugeMetric),
		counters: make(map[string]*CounterMetric),
	}
}

// UpdateGauge updates or sets a gauge metric
func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.gauges[name]; !exists {
		ms.gauges[name] = &GaugeMetric{}
	}
	ms.gauges[name].Set(value)
}

// UpdateCounter updates or increments a counter metric
func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.counters[name]; !exists {
		ms.counters[name] = &CounterMetric{}
	}
	ms.counters[name].Add(value)
}

// GetGauge retrieves a gauge metric
func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metric, exists := ms.gauges[name]
	if !exists {
		return 0, fmt.Errorf("gauge metric %s not found", name)
	}
	return metric.value, nil
}

// GetCounter retrieves a counter metric
func (ms *MemStorage) GetCounter(name string) (int64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	metric, exists := ms.counters[name]
	if !exists {
		return 0, fmt.Errorf("counter metric %s not found", name)
	}
	return metric.value, nil
}
