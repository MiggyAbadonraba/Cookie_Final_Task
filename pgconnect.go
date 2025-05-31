package main

import (
	"database/sql"
	"fmt"
	"os"
)

// InitPostgres читает настройки из переменных окружения и открывает подключение
func InitPostgres() (*sql.DB, error) {
	host := getEnv("PG_HOST", "localhost")
	port := getEnv("PG_PORT", "5432")
	user := getEnv("PG_USER", "postgres")
	pass := getEnv("PG_PASSWORD", "postgres")
	dbname := getEnv("PG_DB", "cookiesdb")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)
	return sql.Open("postgres", dsn)
}

// getEnv возвращает значение переменной окружения или дефолт
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// getAddr возвращает адрес для ListenAndServe
func getAddr() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}
	return ":8080"
}
