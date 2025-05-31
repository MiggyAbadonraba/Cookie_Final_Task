package main

import (
	"os"
	"testing"
)

func TestGetAddr_Default(t *testing.T) {
	os.Unsetenv("HTTP_ADDR")
	addr := getAddr()
	if addr != ":8080" {
		t.Errorf("getAddr: expected \":8080\" when HTTP_ADDR is not set, got %q", addr)
	}
}

func TestGetAddr_FromEnv(t *testing.T) {
	os.Setenv("HTTP_ADDR", ":9999")
	defer os.Unsetenv("HTTP_ADDR")
	addr := getAddr()
	if addr != ":9999" {
		t.Errorf("getAddr: expected \":9999\" from HTTP_ADDR, got %q", addr)
	}
}
