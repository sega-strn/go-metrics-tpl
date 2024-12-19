package main

import (
	"log"
	"net/http"

	"github.com/sega-strn/go-metrics-tpl/internal/handlers"
	"github.com/sega-strn/go-metrics-tpl/internal/storage"
)

func main() {
	// Создаем хранилище метрик
	memStorage := storage.NewMemStorage()

	// Создаем обработчик метрик
	metricsHandler := handlers.NewMetricsHandler(memStorage)

	// Настраиваем роутинг
	http.HandleFunc("/update/", metricsHandler.HandleUpdateMetric)

	// Запускаем сервер
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
