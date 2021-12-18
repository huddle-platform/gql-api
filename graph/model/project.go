package model

type Project struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Languages      []string  `json:"languages"`
	Location       *Location `json:"location"`
	ParticipantIDs []string
}
