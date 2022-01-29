package resolvers

import (
	"context"

	"github.com/google/uuid"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *Resolver) UserFromID(ctx context.Context, id string) (*model.User, error) {
	user, err := r.queries.GetUserByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID.String(),
		Username: &user.Username,
		
	}, nil
}

func (r *Resolver)UserIdFromusername(ctx context.Context, username string) (string, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return user.ID.String(), nil
}