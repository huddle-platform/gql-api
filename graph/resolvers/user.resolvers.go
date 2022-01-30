package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"gitlab.lrz.de/projecthub/gql-api/auth"
	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
	"gitlab.lrz.de/projecthub/gql-api/sqlc"
)

func (r *mutationResolver) SetMyUsername(ctx context.Context, username string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	existingUser, err := r.queries.GetUserByUsername(context.Background(), username)
	if err == nil && existingUser.ID.String() != me.Id {
		return false, fmt.Errorf("username %s already taken", username)
	}

	err = r.queries.SetUserName(context.Background(), sqlc.SetUserNameParams{Username: username, ID: uuid.MustParse(me.Id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) SetMyName(ctx context.Context, name *string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	nameToSet := sql.NullString{
		Valid: false,
	}
	if name != nil {
		nameToSet.String = *name
		nameToSet.Valid = true
	}

	err = r.queries.SetName(context.Background(), sqlc.SetNameParams{Name: nameToSet, ID: uuid.MustParse(me.Id)})
	return err == nil, err
}

func (r *mutationResolver) SetMyDescription(ctx context.Context, description string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.SetDescription(context.Background(), sqlc.SetDescriptionParams{Description: description, ID: uuid.MustParse(me.Id)})
	return err == nil, err
}

func (r *mutationResolver) SetMyProfileImage(ctx context.Context, profileImage *string) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.SetProfileImage(context.Background(), sqlc.SetProfileImageParams{ProfileImage: sql.NullString{
		String: *profileImage,
		Valid:  profileImage != nil,
	}, ID: uuid.MustParse(me.Id)})
	return err == nil, err
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err == nil {
		email, _ := me.GetTrait("email")
		user, err := r.queries.GetUserByID(context.Background(), uuid.MustParse(me.Id))
		if err != nil {
			// create user in database if not exists
			user := sqlc.CreateUserParams{
				ID: uuid.MustParse(me.Id),
			}
			r.queries.CreateUser(context.Background(), user)
			return &model.User{
				ID:    me.Id,
				Email: &email,
			}, nil
		}
		res := model.DBUserToUser(user)
		res.Email = &email
		return res, nil
	}
	return nil, err
}

func (r *queryResolver) MeIfLoggedIn(ctx context.Context) (*model.User, error) {
	res, _ := r.Me(ctx)
	return res, nil
}

func (r *queryResolver) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return model.DBUserToUser(user), nil
}

func (r *userResolver) Participations(ctx context.Context, obj *model.User) ([]*model.Project, error) {
	dbResults, err := r.queries.GetUserParticipations(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	results := model.ProjectsFromDBProjects(dbResults)
	return results, nil
}

func (r *userResolver) CreatedProjects(ctx context.Context, obj *model.User) ([]*model.Project, error) {
	dbResults, err := r.queries.GetProjectsByUserID(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	results := model.ProjectsFromDBProjects(dbResults)
	return results, nil
}

func (r *userResolver) Description(ctx context.Context, obj *model.User) (string, error) {
	user, err := r.queries.GetUserByID(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return "", err
	}
	return user.Description, nil
}

func (r *userResolver) ProfileImage(ctx context.Context, obj *model.User) (*string, error) {
	user, err := r.queries.GetUserByID(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	if user.ProfileImage.Valid {
		return &user.ProfileImage.String, nil
	}
	return nil, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
