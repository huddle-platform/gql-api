package model

type User struct {
	ID         int  `json:"id"`
	Username   string  `json:"username"`
	Email      *string `json:"email"`
	ProjectIDs *[]string
}
