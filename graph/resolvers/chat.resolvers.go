package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *chatResolver) WriteMessage(ctx context.Context, obj *model.Chat, content string) (*model.Message, error) {
	// Write message here
	return &model.Message{
		Content: content,
		Author:  model.MessageAuthorMe,
		Time:    time.Now(),
	}, nil
}

func (r *mutationResolver) Chat(ctx context.Context, id string) (*model.Chat, error) {
	return &model.Chat{
		ID: "1",
		With: &model.User{
			ID:       "1234",
			Username: "peterschlonz42",
		},
	}, nil
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type chatResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

func (r *queryResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
	return []*model.Chat{
		{
			ID: "1",
			With: &model.User{
				ID:       "1234",
				Username: "peterschlonz42",
			},
			Messages: []*model.Message{
				{
					Author:  model.MessageAuthorMe,
					Time:    time.Now(),
					Content: "Hello World",
				},
				{
					Author:  model.MessageAuthorOther,
					Time:    time.Now(),
					Content: "You tooo",
				}},
		}}, nil
}
