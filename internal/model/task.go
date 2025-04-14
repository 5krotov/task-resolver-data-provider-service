package model

import "time"

type Task struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	Difficulty int       `db:"difficulty"`
	Status     int       `db:"status"`
	LastUpdate time.Time `db:"last_update"`
}
