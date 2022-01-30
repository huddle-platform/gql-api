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
	"gitlab.lrz.de/projecthub/gql-api/sqlc"
)

func (r *chatResolver) Me(ctx context.Context, obj *model.Chat) (*model.User, error) {
	return r.UserFromID(ctx, obj.Me_id)
}

func (r *chatResolver) OtherUser(ctx context.Context, obj *model.Chat) (*model.User, error) {
	if obj.ChatType != model.ChatTypeUser {
		return nil, nil
	}
	return r.UserFromID(ctx, obj.Other_id)
}

func (r *chatResolver) OtherProject(ctx context.Context, obj *model.Chat) (*model.Project, error) {
	if obj.ChatType != model.ChatTypeProject {
		return nil, nil
	}
	dbProject, err := r.queries.GetProjectByID(context.Background(), uuid.MustParse(obj.Other_id))
	if err != nil {
		return nil, err
	}
	return model.ProjectFromDBProject(dbProject), nil
}

func (r *chatResolver) Messages(ctx context.Context, obj *model.Chat, until *time.Time, count int) ([]*model.Message, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if obj.Me_id != me.Id {
		return nil, fmt.Errorf("you are not allowed to read messages of other users")
	}
	untilDate := time.Now()
	if until != nil {
		untilDate = *until
	}
	dbChat, err := r.queries.GetMessagesBetweenUsers(ctx, sqlc.GetMessagesBetweenUsersParams{
		SenderID:   uuid.MustParse(obj.Me_id),
		ReceiverID: uuid.MustParse(obj.Other_id),
		Time:       untilDate,
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
	err = r.queries.WriteMessage(context.Background(), sqlc.WriteMessageParams{
		SenderID:   uuid.MustParse(me.Id),
		ReceiverID: uuid.MustParse(userID),
		Content:    content,
	})
	return err == nil, nil
}

func (r *mutationResolver) WriteMessageToProject(ctx context.Context, projectID string, content string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.WriteUserMessageToProject(context.Background(), sqlc.WriteUserMessageToProjectParams{
		UserID:    uuid.MustParse(me.Id),
		ProjectID: uuid.MustParse(projectID),
		Content:   content,
	})
	return err == nil, nil
}

func (r *projectResolver) Chats(ctx context.Context, obj *model.Project) ([]*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if me.Id != obj.CreatorID {
		return nil, fmt.Errorf("you are not allowed to read chats of other projects")
	}
	dbChats, err := r.queries.GetChatsWithProject(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	res := make([]*model.Chat, len(dbChats))
	for i, dbChat := range dbChats {
		res[i] = &model.Chat{
			Me_id:    obj.ID,
			Other_id: dbChat.String(),
			ChatType: model.ChatTypeProject,
		}
	}
	return res, nil
}

func (r *projectMutationResolver) WriteMessageToUser(ctx context.Context, obj *model.ProjectMutation, userID string, content string) (bool, error) {
	err := r.queries.WriteProjectMessageToUser(context.Background(), sqlc.WriteProjectMessageToUserParams{
		ProjectID: uuid.MustParse(obj.ID),
		UserID:    uuid.MustParse(userID),
		Content:   content,
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
	projectUserChats, err := r.queries.GetChatsWithCreatedProjects(context.Background(), uuid.MustParse(me.Id))
	if err != nil {
		return nil, err
	}
	userUserChats, err := r.queries.GetChatsWithUser(context.Background(), uuid.MustParse(me.Id))
	if err != nil {
		return nil, err
	}
	res := make([]*model.Chat, len(userUserChats)+len(projectUserChats))
	for i, chatPartner := range userUserChats {
		res[i] = &model.Chat{
			Me_id:    me.Id,
			Other_id: chatPartner.String(),
			ChatType: model.ChatTypeUser,
		}
	}
	for i, chatPartner := range projectUserChats {
		res[i+len(userUserChats)] = &model.Chat{
			Me_id:    me.Id,
			Other_id: chatPartner.String(),
			ChatType: model.ChatTypeProject,
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
