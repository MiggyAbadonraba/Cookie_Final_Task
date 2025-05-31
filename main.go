package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Инициализация подключения к БД
	db, err := InitPostgres()
	if err != nil {
		log.Fatalf("[ERROR] connect to Postgres: %v", err)
	}
	defer db.Close()
	log.Println("[INFO] Connected to Postgres")

	// 2. Создание таблицы, если её нет
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS cookies (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            ingredients TEXT NOT NULL,
            description TEXT,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		log.Fatalf("[ERROR] ensure table exists: %v", err)
	}
	log.Println("[INFO] Table ready")

	// 3. Парсинг HTML-шаблона
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("[ERROR] parse template: %v", err)
	}

	// 4. Настройка HTTP-сервера с возможностью graceful shutdown
	srv := &http.Server{Addr: getAddr(), Handler: nil}

	// 5. Регистрация маршрутов
	http.HandleFunc("/", HomeHandler(db, tmpl))
	http.HandleFunc("/api/cookies", CookiesHandler(db))     // GET, POST
	http.HandleFunc("/api/cookies/", CookieByIDHandler(db)) // GET, PUT, DELETE по /api/cookies/{id}
	http.HandleFunc("/api/shutdown", ShutdownHandler(srv))  // POST для остановки сервера
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 6. Перехват системных сигналов для корректной остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[INFO] Shutdown requested via OS signal")
		srv.Shutdown(context.Background())
	}()

	// 7. Запуск сервера
	log.Printf("[INFO] Server listening on %s (DEBUG=%s)", srv.Addr, os.Getenv("DEBUG"))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("[ERROR] ListenAndServe: %v", err)
	}
	log.Println("[INFO] Server stopped")
}
