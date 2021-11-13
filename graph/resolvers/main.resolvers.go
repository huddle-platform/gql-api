package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, user model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return &model.User{
		ID:   "1",
		Name: "John Doe",
	}, nil
}

func (r *queryResolver) Messages(ctx context.Context) ([]*model.Mesage, error) {
	res, err := r.pool.Query(ctx, "SELECT * FROM messages")
	if err != nil {
		return nil, err
	}
	var messages []*model.Mesage
	for res.Next() {
		var msg model.Mesage
		err = res.Scan(&msg.ID, &msg.Message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
