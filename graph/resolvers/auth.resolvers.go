package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"

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
	body:=fmt.Sprintf("Here you have your verification code you can do nothing with: %d", rand.Intn(1000000))
	_,err:=http.Post(fmt.Sprintf("http://mail-service:8080/send?to=%s&subject=verification_code",email), "text/plain", bytes.NewReader([]byte(body)))
	if err!=nil{
		return &model.EmailRequestResponse{
			Sent: false,
			Message: err.Error(),
		}, nil
	}
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
