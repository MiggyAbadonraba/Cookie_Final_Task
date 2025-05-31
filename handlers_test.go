package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShutdownHandler_MethodNotAllowed(t *testing.T) {
	srv := &http.Server{}
	handler := ShutdownHandler(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/shutdown", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHomeHandler_GET(t *testing.T) {

	tmpl := template.Must(template.New("index").Parse("<html><body>Hello, world!</body></html>"))
	handler := HomeHandler(nil, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "<html>") {
		t.Errorf("template not rendered correctly")
	}

}

func TestCookiesHandler_GET(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "ingredients", "description", "created_at"}).
		AddRow(1, "Choco", "Sugar", "Yummy", time.Now())
	mock.ExpectQuery("SELECT id,name,ingredients,description,created_at FROM cookies").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/api/cookies", nil)
	w := httptest.NewRecorder()
	handler := CookiesHandler(db)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	var result []Cookie
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil || len(result) != 1 {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestCookieByIDHandler_GET(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name", "ingredients", "description", "created_at"}).
		AddRow(1, "Choco", "Flour", "Tasty", time.Now())
	mock.ExpectQuery("SELECT id,name,ingredients,description,created_at FROM cookies WHERE id=\\$1").
		WithArgs(1).WillReturnRows(row)

	req := httptest.NewRequest(http.MethodGet, "/api/cookies/1", nil)
	w := httptest.NewRecorder()
	handler := CookiesHandler(db)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestCookieByIDHandler_PUT(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("UPDATE cookies SET name=\\$1,ingredients=\\$2,description=\\$3,created_at=\\$4 WHERE id=\\$5").
		WithArgs("Choco", "Flour", "Good", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := `{"name":"Choco","ingredients":"Flour","description":"Good"}`
	req := httptest.NewRequest(http.MethodPut, "/api/cookies/1", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler := CookieByIDHandler(db)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestCookieByIDHandler_DELETE(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("DELETE FROM cookies WHERE id=\\$1").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/api/cookies/1", nil)
	w := httptest.NewRecorder()
	handler := CookieByIDHandler(db)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}
