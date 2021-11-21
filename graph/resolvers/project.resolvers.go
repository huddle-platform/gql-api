package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project *model.NewProjectInput) (*model.Project, error) {
	return &model.Project{
		ID:             "1",
		Name:           "test",
		Description:    "Example Project. We are doing great stuff.",
		Languages:      []string{"DE", "EN"},
		Location:       &model.Location{Name: "Online"},
		ParticipantIDs: []string{"1234", "5353"},
	}, nil
}

func (r *mutationResolver) AddSavedProject(ctx context.Context, id string) ([]*model.Project, error) {
	return []*model.Project{}, nil
}

func (r *projectResolver) Participants(ctx context.Context, obj *model.Project) ([]*model.User, error) {
	participants := make([]*model.User, len(obj.ParticipantIDs))
	for i, id := range obj.ParticipantIDs {
		participants[i] = &model.User{
			ID:       id,
			Username: "user-" + id,
		}
	}
	return participants, nil
}

func (r *projectResolver) AddParticipant(ctx context.Context, obj *model.Project, id string) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectResolver) RemoveParticipant(ctx context.Context, obj *model.Project, id string) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchProjects(ctx context.Context, searchString string, options model.SearchOptions, start int, count int) ([]*model.Project, error) {
	results := make([]*model.Project, count)

	for i := 0; i < count; i++ {
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

func (r *queryResolver) GetProject(ctx context.Context, id string) (*model.Project, error) {
	return &model.Project{
		ID:             string(id),
		Name:           "Project number" + string(id),
		Description:    "Description of project" + string(id),
		Languages:      []string{"DE"},
		Location:       &model.Location{Name: "Location" + string(id)},
		ParticipantIDs: []string{"1234", "5345"}}, nil
}

func (r *queryResolver) SavedProjects(ctx context.Context) ([]*model.Project, error) {
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

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

type projectResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) SaveProject(ctx context.Context, id string) (*model.Project, error) {
	panic(fmt.Errorf("not implemented"))
}
