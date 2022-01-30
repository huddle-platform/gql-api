package model

import "gitlab.lrz.de/projecthub/gql-api/sqlc"

type User struct {
	ID       string  `json:"id"`
	Username string `json:"username"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
}

type FeedbackMutation struct {
	UserId *string
}

func DBUserToUser(dbUser sqlc.User) *User {
	var name *string
	if dbUser.Name.Valid {
		name = &dbUser.Name.String
	}
	return &User{
		ID:       dbUser.ID.String(),
		Username: dbUser.Username,
		Name:     name,
	}
}
