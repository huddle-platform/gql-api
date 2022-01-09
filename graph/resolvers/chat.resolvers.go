package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

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

func (r *mutationResolver) WriteMessageToUser(ctx context.Context, userID string, content string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.WriteMessage(context.Background(), sql.WriteMessageParams{
		SenderID:   uuid.MustParse(me.Id),
		ReceiverID: uuid.MustParse(userID),
		Content:    content,
	})
	return err == nil, nil
}

func (r *queryResolver) Chat(ctx context.Context, with string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		Me_id:    me.Id,
		Other_id: with,
	}, nil
}

func (r *queryResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
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

func (r *queryResolver) GetChatByUsername(ctx context.Context, withUsername string) (*model.Chat, error) {
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

func (r *queryResolver) GetChatByID(ctx context.Context, withUserID string) (*model.Chat, error) {
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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type chatResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
