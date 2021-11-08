package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *authResolver) Login(ctx context.Context, obj *model.Auth, token string) (*model.LoginResponse, error) {
	fmt.Println("login")
	jwt := "asdfoiafvboaiduf"
	return &model.LoginResponse{
		Error: nil,
		Jwt:   &jwt,
	}, nil
}

func (r *authResolver) RequestLoginEmail(ctx context.Context, obj *model.Auth, email string) (*model.EmailRequestResponse, error) {
	fmt.Println("would be sending email here...")
	return &model.EmailRequestResponse{
		Sent:    true,
		Message: "Successfully sent email",
	}, nil
}

func (r *mutationResolver) Auth(ctx context.Context) (*model.Auth, error) {
	return &model.Auth{}, nil
}

// Auth returns generated.AuthResolver implementation.
func (r *Resolver) Auth() generated.AuthResolver { return &authResolver{r} }

type authResolver struct{ *Resolver }

