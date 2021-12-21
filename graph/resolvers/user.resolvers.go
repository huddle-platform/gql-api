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

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	me, exists := auth.IdentityFromContext(ctx)
	if exists {
		email, _ := me.GetTrait("email")
		user, err := r.queries.GetUserByID(context.Background(), uuid.MustParse(me.Id))
		if err != nil {
			// create user in database if not exists
			user := sql.CreateUserParams{
				ID:    uuid.MustParse(me.Id),
				Email: email,
			}
			r.queries.CreateUser(context.Background(), user)
			return &model.User{
				ID:    me.Id,
				Email: &email,
			}, nil
		}
		return &model.User{
			ID:       me.Id,
			Email:    &email,
			Username: &user.Username,
		}, nil
	}
	return nil, fmt.Errorf("not logged in")
}

func (r *userResolver) Projects(ctx context.Context, obj *model.User) ([]*model.Project, error) {
	dbResults, err := r.queries.GetProjectsByUserID(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	results := make([]*model.Project, len(dbResults))
	for i, res := range dbResults {
		{
			results[i] = &model.Project{
				ID:             res.ID.String(),
				Name:           res.Name,
				Description:    res.Description,
				ParticipantIDs: []string{obj.ID},
			}
		}

	}
	return results, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
