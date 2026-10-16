package domain

import "time"

type Task struct {
	ID       int
	Title    string
	Status   string // "planned", "in_progress", "cancelled", "done"
	Deadline time.Time
}
