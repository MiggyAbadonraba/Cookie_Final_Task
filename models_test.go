package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestCookieJSONTags(t *testing.T) {
	c := Cookie{
		ID:          42,
		Name:        "Test",
		Ingredients: "Flour, Sugar",
		Description: "Tasty",
		CreatedAt:   time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshall returned error: %v", err)
	}

	jsonStr := string(data)

	expectedFields := []string{"\"id\":", "\"name\":", "\"ingredients\":", "\"description\":", "\"created_at\":"}
	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON output missing field %s in %s", field, jsonStr)
		}
	}
}
