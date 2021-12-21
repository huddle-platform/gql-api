package model
type Chat struct {
	ID           string     `json:"id"`
	With         *User      `json:"with"`
	Messages     []*Message `json:"messages"`
}