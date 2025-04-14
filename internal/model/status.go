package model

import "time"

type Status struct {
	ID        int64     `db:"id"`
	Status    int       `db:"status"`
	Timestamp time.Time `db:"timestamp"`
	TaskID    int64     `db:"task_id"`
}
