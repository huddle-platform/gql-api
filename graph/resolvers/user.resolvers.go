package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	email := "peter.schlonz@gmail.com"
	return &model.User{
		ID:       "1",
		Username: "Testuser",
		Email:    &email,
	}, nil
}

func (r *userResolver) Projects(ctx context.Context, obj *model.User) ([]*model.Project, error) {
	results := make([]*model.Project, 10)

	for i := 0; i < 10; i++ {
		results[i] = &model.Project{
			ID:             string(i),
			Name:           "Project number" + string(i),
			Description:    "Description of project" + string(i),
			Languages:      []string{"DE"},
			Location:       &model.Location{Name: "Location" + string(i)},
			ParticipantIDs: []string{"1234", "5345"},
		}
	}

	return results, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
