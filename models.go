package main

import "time"

//Наши данные, которые будут отправляться в таблицу

type Cookie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Ingredients string    `json:"ingredients"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
