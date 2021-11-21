package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	dbsql"database/sql"

	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
	"gitlab.lrz.de/projecthub/gql-api/sql"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project *model.NewProjectInput) (*model.Project, error) {
	projectID, err := r.queries.CreateProject(ctx, sql.CreateProjectParams{Name: project.Name, Description: project.Description})
	if err != nil {
		return nil, err
	}
	roleId, err := r.queries.AddRole(ctx, sql.AddRoleParams{Type: "project-admin",ProjectID: dbsql.NullInt32{Int32:projectID}})
	if err != nil {
		return nil, err
	}
	err=r.queries.GrantRoleToUser(ctx,sql.GrantRoleToUserParams{UserID: 1,RoleID: roleId})
	if err != nil {
		return nil, err
	}
	fetchedProject, err := r.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &model.Project{
		ID:          int(fetchedProject.ID),
		Name:        fetchedProject.Name,
		Description: fetchedProject.Description,
	}, nil
}

func (r *mutationResolver) AddSavedProject(ctx context.Context, id int) ([]*model.Project, error) {
	return []*model.Project{}, nil
}

func (r *projectResolver) Participants(ctx context.Context, obj *model.Project) ([]*model.User, error) {
	participants := make([]*model.User, len(obj.ParticipantIDs))
	for i, id := range obj.ParticipantIDs {
		participants[i] = &model.User{
			ID:       id,
			Username: "user-" + string(id),
		}
	}
	return participants, nil
}

func (r *projectResolver) AddParticipant(ctx context.Context, obj *model.Project, id int) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *projectResolver) RemoveParticipant(ctx context.Context, obj *model.Project, id int) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchProjects(ctx context.Context, searchString string, options model.SearchOptions, start int, count int) ([]*model.Project, error) {
	results := make([]*model.Project, count)

	for i := 0; i < count; i++ {
		results[i] = &model.Project{
			ID:             (i),
			Name:           "Project number" + string(i),
			Description:    "Description of project" + string(i),
			Languages:      []string{"DE"},
			Location:       &model.Location{Name: "Location" + string(i)},
			ParticipantIDs: []int{1234, 5345},
		}
	}

	return results, nil
}

func (r *queryResolver) GetProject(ctx context.Context, id int) (*model.Project, error) {
	return &model.Project{
		ID:             (id),
		Name:           "Project number" + string(id),
		Description:    "Description of project" + string(id),
		Languages:      []string{"DE"},
		Location:       &model.Location{Name: "Location" + string(id)},
		ParticipantIDs: []int{1234, 5345}}, nil
}

func (r *queryResolver) SavedProjects(ctx context.Context) ([]*model.Project, error) {
	results := make([]*model.Project, 10)

	for i := 0; i < 10; i++ {
		results[i] = &model.Project{
			ID:             (i),
			Name:           "Project number" + string(i),
			Description:    "Description of project" + string(i),
			Languages:      []string{"DE"},
			Location:       &model.Location{Name: "Location" + string(i)},
			ParticipantIDs: []int{1234, 5345}}
	}
	return results, nil
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

type projectResolver struct{ *Resolver }
