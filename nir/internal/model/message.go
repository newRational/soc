package model

import "time"

type Message struct {
	ID          uint64     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Level       string     `db:"level"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
