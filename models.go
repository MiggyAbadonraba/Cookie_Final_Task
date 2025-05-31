package main

import "time"

type Cookie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Ingredients string    `json:"ingredients"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
