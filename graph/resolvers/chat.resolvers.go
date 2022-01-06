package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.lrz.de/projecthub/gql-api/auth"
	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
	"gitlab.lrz.de/projecthub/gql-api/sql"
)

func (r *chatResolver) Me(ctx context.Context, obj *model.Chat) (*model.User, error) {
	return r.UserFromID(ctx, obj.Me_id)
}

func (r *chatResolver) Other(ctx context.Context, obj *model.Chat) (*model.User, error) {
	return r.UserFromID(ctx, obj.Other_id)
}

func (r *chatResolver) Messages(ctx context.Context, obj *model.Chat, start int, count int) ([]*model.Message, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if obj.Me_id != me.Id {
		return nil, fmt.Errorf("you are not allowed to read messages of other users")
	}
	dbChat, err := r.queries.GetMessagesBetweenUsers(ctx, sql.GetMessagesBetweenUsersParams{
		SenderID:   uuid.MustParse(obj.Me_id),
		ReceiverID: uuid.MustParse(obj.Other_id),
		Offset:     int32(start),
		Limit:      int32(count)})
	if err != nil {
		return nil, err
	}
	res := make([]*model.Message, len(dbChat))
	for i, dbMessage := range dbChat {
		author := model.MessageAuthorMe
		if dbMessage.SenderID.String() != obj.Me_id {
			author = model.MessageAuthorOther
		}
		res[i] = &model.Message{
			Author:  author,
			Content: dbMessage.Content,
			Time:    dbMessage.Time,
		}
	}
	return res, nil
}

func (r *chatResolver) WriteMessage(ctx context.Context, obj *model.Chat, content string) (*model.Message, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if obj.Me_id != me.Id {
		return nil, fmt.Errorf("you are not allowed to write messages for other users")
	}
	err = r.queries.WriteMessage(context.Background(), sql.WriteMessageParams{
		SenderID:   uuid.MustParse(obj.Me_id),
		ReceiverID: uuid.MustParse(obj.Other_id),
		Content:    content,
	})
	if err != nil {
		return nil, err
	}
	return &model.Message{
		Content: content,
		Author:  model.MessageAuthorMe,
		Time:    time.Now(),
	}, nil
}

func (r *mutationResolver) Chat(ctx context.Context, with string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		Me_id:    me.Id,
		Other_id: with,
	}, nil
}

func (r *mutationResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	chatPartners, err := r.queries.GetChatsWithUser(context.Background(), uuid.MustParse(me.Id))
	if err != nil {
		return nil, err
	}
	res := make([]*model.Chat, len(chatPartners))
	for i, chatPartner := range chatPartners {
		res[i] = &model.Chat{
			Me_id:    me.Id,
			Other_id: chatPartner.String(),
		}
	}
	return res, nil
}

func (r *mutationResolver) GetChatByUsername(ctx context.Context, withUsername string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	otherID, err := r.UserIdFromusername(context.Background(), withUsername)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		Me_id:    me.Id,
		Other_id: otherID,
	}, nil
}

func (r *mutationResolver) GetChatByID(ctx context.Context, withUserID string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		Me_id:    me.Id,
		Other_id: withUserID,
	}, nil
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type chatResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) NewChat(ctx context.Context, withUsername string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	otherID, err := r.UserIdFromusername(context.Background(), withUsername)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		Me_id:    me.Id,
		Other_id: otherID,
	}, nil
}
