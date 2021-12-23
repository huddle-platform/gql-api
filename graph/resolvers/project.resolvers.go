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

func (r *mutationResolver) CreateProject(ctx context.Context, project *model.NewProjectInput) (*model.Project, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	projectID, err := r.queries.CreateProject(ctx, sql.CreateProjectParams{Name: project.Name, Description: project.Description, Creator: uuid.MustParse(user.Id)})
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

func (r *mutationResolver) AddSavedProject(ctx context.Context, id string) (bool, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.SaveProjectForUser(context.Background(), sql.SaveProjectForUserParams{UserID: uuid.MustParse(user.Id), ProjectID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RemoveSavedProject(ctx context.Context, id string) (bool, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.UnsaveProjectForUser(context.Background(), sql.UnsaveProjectForUserParams{UserID: uuid.MustParse(user.Id), ProjectID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *projectResolver) Participants(ctx context.Context, obj *model.Project) ([]*model.User, error) {
	dbParticipants, err := r.queries.GetParticipantsOfProject(context.Background(), uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	participants := make([]*model.User, len(dbParticipants))
	for i, p := range dbParticipants {
		participants[i] = &model.User{
			ID:       p.ID.String(),
			Username: &p.Username,
			Email:    &p.Email,
		}
	}
	return participants, nil
}

func (r *projectResolver) Creator(ctx context.Context, obj *model.Project) (*model.User, error) {
	dbUser, err := r.queries.GetUserByID(context.Background(), uuid.MustParse(obj.CreatorID))
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       dbUser.ID.String(),
		Username: &dbUser.Username,
		Email:    &dbUser.Email,
	}, nil
}

func (r *projectResolver) AddParticipant(ctx context.Context, obj *model.Project, id string) (bool, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	if user.Id != obj.CreatorID {
		return false, fmt.Errorf("only the creator can add participants")
	}
	err = r.queries.AddParticipantToProject(context.Background(), sql.AddParticipantToProjectParams{ProjectID: uuid.MustParse(obj.ID), UserID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *projectResolver) RemoveParticipant(ctx context.Context, obj *model.Project, id string) (bool, error) {
	err := r.queries.RemoveParticipantFromProject(context.Background(), sql.RemoveParticipantFromProjectParams{ProjectID: uuid.MustParse(obj.ID), UserID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
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
			ID:          dbProject.ID.String(),
			Name:        dbProject.Name,
			Description: dbProject.Description,
			Languages:   []string{"DE"},
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
		ID:          dbProject.ID.String(),
		Name:        dbProject.Name,
		Description: dbProject.Description,
	}, nil
}

func (r *queryResolver) SavedProjects(ctx context.Context) ([]*model.Project, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	dbResults, err := r.queries.GetSavedProjectsForUser(context.Background(), uuid.MustParse(user.Id))
	if err != nil {
		return nil, err
	}
	results := make([]*model.Project, len(dbResults))
	for i, p := range dbResults {
		results[i] = &model.Project{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Languages:   []string{"DE"},
		}
	}
	return results, nil
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

type projectResolver struct{ *Resolver }
