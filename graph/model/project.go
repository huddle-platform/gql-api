package model

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Languages   []string  `json:"languages"`
	Location    *Location `json:"location"`
	CreatedAt   *time.Time `json:"createdAt"`
	CreatorID   string
}

type ProjectMutation struct {
	ID string
}
