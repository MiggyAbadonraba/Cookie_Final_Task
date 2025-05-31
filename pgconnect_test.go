package main

import (
	"os"
	"testing"
)

func TestGetEnv_Default(t *testing.T) {
	os.Unsetenv("SOME_KEY")
	val := getEnv("SOME_KEY", "defaultVal")
	if val != "defaultVal" {
		t.Errorf("getEnv: expected defaultVal, got %q", val)
	}
}

func TestGetEnv_FromEnv(t *testing.T) {
	os.Setenv("SOME_KEY", "value123")
	defer os.Unsetenv("SOME_KEY")
	val := getEnv("SOME_KEY", "defaultVal")
	if val != "value123" {
		t.Errorf("getEnv: expected value123, got %q", val)
	}
}

func TestInitPostgres_Basic(t *testing.T) {
	os.Unsetenv("PG_HOST")
	os.Unsetenv("PG_PORT")
	os.Unsetenv("PG_USER")
	os.Unsetenv("PG_PASSWORD")
	os.Unsetenv("PG_DB")
	db, err := InitPostgres()

	if err != nil {
		t.Fatalf("InitPostgres returned error: %v", err)
	}

	if db == nil {
		t.Errorf("InitPostgres: expected non-nil *sql.DB, got nil")
	}
}
