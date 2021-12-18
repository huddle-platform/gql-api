package model
type Chat struct {
	ID           int     `json:"id"`
	With         *User      `json:"with"`
	Messages     []*Message `json:"messages"`
}