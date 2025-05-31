package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ShutdownHandler возвращает handler для graceful shutdown по POST /api/shutdown
func ShutdownHandler(srv *http.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Println("[INFO] Shutdown requested via API")
		w.Write([]byte("shutting down"))

		// Останавливаем сервер через секунду (даём время JS получить ответ)
		go func() {
			time.Sleep(1 * time.Second)
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("[ERROR] shutdown: %v", err)
			}
		}()
	}
}

func HomeHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// … ваш существующий POST-код с log.Printf и INSERT
		}
		// просто отрисовка страницы (без списка по умолчанию)
		data := struct {
			Debug bool
		}{
			Debug: os.Getenv("DEBUG") == "true",
		}
		tmpl.Execute(w, data)
	}
}

func CookiesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем каждый наш запрос в консоли
		log.Printf("[HTTP] %s %s", r.Method, r.URL.Path)

		switch r.Method {
		case http.MethodGet:
			// Возвращаем все записи
			rows, _ := db.Query("SELECT id,name,ingredients,description,created_at FROM cookies ORDER BY created_at DESC")
			defer rows.Close()
			var out []Cookie
			for rows.Next() {
				var c Cookie
				rows.Scan(&c.ID, &c.Name, &c.Ingredients, &c.Description, &c.CreatedAt)
				out = append(out, c)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(out)

		case http.MethodPost:
			// Создаём новую запись
			var c Cookie
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &c)
			db.Exec(
				"INSERT INTO cookies (name, ingredients, description) VALUES ($1,$2,$3)",
				c.Name, c.Ingredients, c.Description,
			)
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func CookieByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем каждый наш запрос в консоли
		log.Printf("[HTTP] %s %s", r.Method, r.URL.Path)

		// Извлечь ID из URL /api/cookies/{id}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		id, _ := strconv.Atoi(parts[len(parts)-1])

		switch r.Method {
		case http.MethodGet:
			var c Cookie
			db.QueryRow(
				"SELECT id,name,ingredients,description,created_at FROM cookies WHERE id=$1", id,
			).Scan(&c.ID, &c.Name, &c.Ingredients, &c.Description, &c.CreatedAt)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)

		case http.MethodPut:
			var c Cookie
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &c)
			db.Exec(
				"UPDATE cookies SET name=$1,ingredients=$2,description=$3,created_at=$4 WHERE id=$5",
				c.Name, c.Ingredients, c.Description, time.Now(), id,
			)
			w.WriteHeader(http.StatusOK)

		case http.MethodDelete:
			db.Exec("DELETE FROM cookies WHERE id=$1", id)
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
