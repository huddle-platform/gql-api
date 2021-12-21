package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
	"gitlab.lrz.de/projecthub/gql-api/sql"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project *model.NewProjectInput) (*model.Project, error) {
	projectID, err := r.queries.CreateProject(ctx, sql.CreateProjectParams{Name: project.Name, Description: project.Description})
	if err != nil {
		return nil, err
	}
	fetchedProject, err := r.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &model.Project{
		ID:          fetchedProject.ID.String(),
		Name:        fetchedProject.Name,
		Description: fetchedProject.Description,
	}, nil
}

func (r *mutationResolver) AddSavedProject(ctx context.Context, id string) ([]*model.Project, error) {
	return []*model.Project{}, nil
}

func (r *projectResolver) Participants(ctx context.Context, obj *model.Project) ([]*model.User, error) {
	participants := make([]*model.User, len(obj.ParticipantIDs))
	for i, id := range obj.ParticipantIDs {
		participants[i] = &model.User{
			ID: id,
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

func (r *queryResolver) SearchProjects(ctx context.Context, searchString string, options model.SearchOptions, offset int, countLimit int) ([]*model.Project, error) {
	dbResults, err := r.queries.GetProjects(context.Background())
	if err != nil {
		return nil, err
	}
	results := []*model.Project{}

	for i := 0; i < countLimit && i+offset < len(dbResults); i++ {
		dbProject := dbResults[i+offset]
		results = append(results, &model.Project{
			ID:             dbProject.ID.String(),
			Name:           dbProject.Name,
			Description:    dbProject.Description,
			Languages:      []string{"DE"},
			ParticipantIDs: []string{dbProject.Creator.String()},
		})
	}
	return results, nil
}

func (r *queryResolver) GetProject(ctx context.Context, id string) (*model.Project, error) {
	dbProject, err := r.queries.GetProjectByID(context.Background(), uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return &model.Project{
		ID:             dbProject.ID.String(),
		Name:           dbProject.Name,
		Description:    dbProject.Description,
		ParticipantIDs: []string{dbProject.Creator.String()},
	}, nil
}

func (r *queryResolver) SavedProjects(ctx context.Context) ([]*model.Project, error) {
	results := make([]*model.Project, 10)

	for i := 0; i < 10; i++ {
		results[i] = &model.Project{
			ID:             "ndpifp",
			Name:           "Project number" + string(i),
			Description:    "Description of project" + string(i),
			Languages:      []string{"DE"},
			Location:       &model.Location{Name: "Location" + string(i)},
			ParticipantIDs: []string{"1234", "5345"}}
	}
	return results, nil
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

type projectResolver struct{ *Resolver }
