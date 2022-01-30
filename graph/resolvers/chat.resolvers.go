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

func (r *chatResolver) Messages(ctx context.Context, obj *model.Chat, until *time.Time, count int) ([]*model.Message, error) {
	// The user should be authenticated at this point
	untilDate := time.Now()
	if until != nil {
		untilDate = *until
	}
	if obj.P1.User != nil && obj.P2.User != nil {
		// user user chatd=obj.P2.User.ID

		dbChat, err := r.queries.GetMessagesBetweenUsers(ctx, sqlc.GetMessagesBetweenUsersParams{
			SenderID:   uuid.MustParse(obj.P1.User.ID),
			ReceiverID: uuid.MustParse(obj.P2.User.ID),
			Time:       untilDate,
			Limit:      int32(count)})
		if err != nil {
			return nil, err
		}
		res := make([]*model.Message, len(dbChat))
		for i, dbMessage := range dbChat {
			author := model.MessageAuthorP1
			if dbMessage.SenderID.String() != obj.P1.User.ID {
				author = model.MessageAuthorP2
			}
			res[i] = &model.Message{
				Author:  author,
				Content: dbMessage.Content,
				Time:    dbMessage.Time,
			}
		}
		return res, nil
	} else {
		project := obj.P1.Project
		user := obj.P2.User
		p1IsUser := false
		if project == nil {
			project = obj.P2.Project
			user = obj.P1.User
			p1IsUser = true
			if user == nil {
				return nil, fmt.Errorf("invalid chat")
			}
		}

		dbChat, err := r.queries.GetMessagesBetweenUserAndProject(ctx, sqlc.GetMessagesBetweenUserAndProjectParams{
			UserID:    uuid.MustParse(user.ID),
			ProjectID: uuid.MustParse(project.ID),
			Time:      untilDate,
			Limit:     int32(count)})
		if err != nil {
			return nil, err
		}
		res := make([]*model.Message, len(dbChat))
		for i, dbMessage := range dbChat {
			p1IsSender := dbMessage.Userissender
			if !p1IsUser {
				p1IsSender = !p1IsSender
			}
			author := model.MessageAuthorP1
			if !p1IsSender {
				author = model.MessageAuthorP2
			}
			res[i] = &model.Message{
				Author:  author,
				Content: dbMessage.Content,
				Time:    dbMessage.Time,
			}
		}
		return res, nil

	}
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
	return err == nil, err
}

func (r *mutationResolver) WriteMessageToProject(ctx context.Context, projectID string, content string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.WriteProjectUserMessage(context.Background(), sqlc.WriteProjectUserMessageParams{
		UserID:       uuid.MustParse(me.Id),
		ProjectID:    uuid.MustParse(projectID),
		Userissender: true,
		Content:      content,
	})
	return err == nil, err
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
			P1: &model.ChatParticipant{
				Project: obj,
			},
			P2: &model.ChatParticipant{
				User: model.DBUserToUser(dbChat),
			},
			Me: model.MessageAuthorP1,
		}
	}
	return res, nil
}

func (r *projectResolver) GetChatByUserID(ctx context.Context, obj *model.Project, withUserID string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if me.Id != obj.CreatorID {
		return nil, fmt.Errorf("you are not allowed to read chats of other projects")
	}
	otherUser, err := r.UserFromID(ctx, withUserID)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		P1: &model.ChatParticipant{
			Project: obj,
		},
		P2: &model.ChatParticipant{
			User: otherUser,
		},
		Me: model.MessageAuthorP1,
	}, nil
}

func (r *projectMutationResolver) WriteMessageToUser(ctx context.Context, obj *model.ProjectMutation, userID string, content string) (bool, error) {
	err := r.queries.WriteProjectUserMessage(context.Background(), sqlc.WriteProjectUserMessageParams{
		ProjectID:    uuid.MustParse(obj.ID),
		UserID:       uuid.MustParse(userID),
		Content:      content,
		Userissender: false,
	})
	return err == nil, err
}

func (r *queryResolver) Chats(ctx context.Context) ([]*model.Chat, error) {
	me, err := r.Me(ctx)
	if err != nil {
		return nil, err
	}
	projectUserChats, err := r.queries.GetChatsWithCreatedProjects(context.Background(), uuid.MustParse(me.ID))
	if err != nil {
		return nil, err
	}
	userUserChats, err := r.queries.GetChatsWithUser(context.Background(), uuid.MustParse(me.ID))
	if err != nil {
		return nil, err
	}

	pingedProjects, err := r.queries.GetProjectChatsWithUser(context.Background(), uuid.MustParse(me.ID))
	if err != nil {
		return nil, err
	}
	res := make([]*model.Chat, len(userUserChats)+len(projectUserChats)+len(pingedProjects))
	for i, chatPartner := range userUserChats {
		res[i] = &model.Chat{
			P1: &model.ChatParticipant{
				User: me,
			},
			P2: &model.ChatParticipant{
				User: model.DBUserToUser(chatPartner),
			},
			Me: model.MessageAuthorP1,
		}
	}
	for i, projectMatch := range projectUserChats {
		// you are the project
		project, _ := r.GetProject(ctx, projectMatch.ID.String())
		user, _ := r.UserFromID(ctx, projectMatch.UserID.String())
		res[i+len(userUserChats)] = &model.Chat{
			P1: &model.ChatParticipant{
				Project: project,
			},
			P2: &model.ChatParticipant{
				User: user,
			},
			Me: model.MessageAuthorP1,
		}
	}
	for i, projectMatch := range pingedProjects {
		// you are the user
		res[i+len(userUserChats)+len(projectUserChats)] = &model.Chat{
			P1: &model.ChatParticipant{
				User: me,
			},
			P2: &model.ChatParticipant{
				Project: model.ProjectFromDBProject(projectMatch),
			},
			Me: model.MessageAuthorP1,
		}
	}
	return res, nil
}

func (r *queryResolver) GetChatByUsername(ctx context.Context, withUsername string) (*model.Chat, error) {
	otherID, err := r.UserIdFromusername(context.Background(), withUsername)
	if err != nil {
		return nil, err
	}
	return r.GetChatByUserID(ctx, otherID)
}

func (r *queryResolver) GetChatByUserID(ctx context.Context, withUserID string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meUser, err := r.UserFromID(ctx, me.Id)
	if err != nil {
		return nil, err
	}
	otherUser, err := r.UserFromID(ctx, withUserID)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		P1: &model.ChatParticipant{
			User: meUser,
		},
		P2: &model.ChatParticipant{
			User: otherUser,
		},
		Me: model.MessageAuthorP1,
	}, nil
}

func (r *queryResolver) GetChatByProjectID(ctx context.Context, withProjectID string) (*model.Chat, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meUser, err := r.UserFromID(ctx, me.Id)
	if err != nil {
		return nil, err
	}
	project, err := r.GetProject(ctx, withProjectID)
	if err != nil {
		return nil, err
	}
	return &model.Chat{
		P1: &model.ChatParticipant{
			User: meUser,
		},
		P2: &model.ChatParticipant{
			Project: project,
		},
		Me: model.MessageAuthorP1,
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
