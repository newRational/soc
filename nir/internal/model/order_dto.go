package model

import "time"

type AddOrderRequest struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Level       string     `json:"level"`
	Description string     `json:"description"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
}

type AddOrderResponse struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Level       string     `json:"level"`
	Description string     `json:"description"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
}

type Response struct {
	Status int
	Data   []byte
}
