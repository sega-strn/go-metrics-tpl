package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sega-strn/go-metrics-tpl/internal/handlers"
	"github.com/sega-strn/go-metrics-tpl/internal/storage"
)

func main() {
	// Создаем хранилище метрик
	memStorage := storage.NewMemStorage()

	// Создаем обработчик метрик
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Создаем новый роутер
	r := mux.NewRouter()

	// Эндпоинт для обновления метрик
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", metricsHandler.HandleUpdateMetric).Methods("POST")

	// Эндпоинт для получения значения метрики
	r.HandleFunc("/value/{metricType}/{metricName}", metricsHandler.HandleGetMetric).Methods("GET")

	// Эндпоинт для получения списка всех метрик
	r.HandleFunc("/", metricsHandler.HandleListMetrics).Methods("GET")

	// Запускаем сервер
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
