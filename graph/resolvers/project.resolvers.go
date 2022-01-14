package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	sqlPackage "database/sql"
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
		CreatorID:   fetchedProject.Creator.String(),
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

func (r *mutationResolver) ProjectMutation(ctx context.Context, id string) (*model.ProjectMutation, error) {
	// check here if user is allowed to modify project
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("authenticate please")
	}

	fetchedProject, err := r.queries.GetProjectByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	if fetchedProject.Creator.String() != me.Id {
		return nil, fmt.Errorf("only the creator can modify a project")
	}
	return &model.ProjectMutation{ID: id}, nil
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
	}, nil
}

func (r *projectResolver) Images(ctx context.Context, obj *model.Project) ([]*model.Image, error) {
	images, err := r.queries.GetImagesOfProject(ctx, uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	toReturn := make([]*model.Image, len(images))
	for i, img := range images {
		toReturn[i] = &model.Image{
			ID:          img.ID.String(),
			URL:         img.Url,
			CreatedAt:   img.CreatedAt.Time,
			Description: &img.Description.String,
			Priority:    float64(img.Priority),
		}
	}
	return toReturn, nil
}

func (r *projectMutationResolver) AddParticipant(ctx context.Context, obj *model.ProjectMutation, id string) (bool, error) {
	err := r.queries.AddParticipantToProject(context.Background(), sql.AddParticipantToProjectParams{ProjectID: uuid.MustParse(obj.ID), UserID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *projectMutationResolver) RemoveParticipant(ctx context.Context, obj *model.ProjectMutation, id string) (bool, error) {
	err := r.queries.RemoveParticipantFromProject(context.Background(), sql.RemoveParticipantFromProjectParams{ProjectID: uuid.MustParse(obj.ID), UserID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *projectMutationResolver) DeleteProject(ctx context.Context, obj *model.ProjectMutation) (bool, error) {
	err := r.queries.DeleteProject(context.Background(), uuid.MustParse(obj.ID))
	return err == nil, err
}

func (r *projectMutationResolver) UpdateDescription(ctx context.Context, obj *model.ProjectMutation, newDescription string) (bool, error) {
	err := r.queries.UpdateProjectDescription(context.Background(), sql.UpdateProjectDescriptionParams{ID: uuid.MustParse(obj.ID), Description: newDescription})
	return err == nil, err
}

func (r *projectMutationResolver) UpdateName(ctx context.Context, obj *model.ProjectMutation, newName string) (bool, error) {
	err := r.queries.UpdateProjectName(context.Background(), sql.UpdateProjectNameParams{ID: uuid.MustParse(obj.ID), Name: newName})
	return err == nil, err
}

func (r *projectMutationResolver) AddImage(ctx context.Context, obj *model.ProjectMutation, newImage model.NewImageInput) (bool, error) {
	priority := 0.0
	if newImage.Priority != nil {
		priority = *newImage.Priority
	}
	err := r.queries.AddImageToProject(context.Background(), sql.AddImageToProjectParams{
		Project:     uuid.MustParse(obj.ID),
		Url:         newImage.URL,
		Description: sqlPackage.NullString{String: *newImage.Description, Valid: newImage.Description != nil},
		Priority:    float32(priority),
	})
	return err == nil, err
}

func (r *projectMutationResolver) RemoveImage(ctx context.Context, obj *model.ProjectMutation, id string) (bool, error) {
	err := r.queries.DeleteImage(context.Background(), uuid.MustParse(id))
	return err == nil, err
}

func (r *projectMutationResolver) UpdateImageDescription(ctx context.Context, obj *model.ProjectMutation, id string, newDescription *string) (bool, error) {
	err := r.queries.UpdateImageDescription(context.Background(), sql.UpdateImageDescriptionParams{ID: uuid.MustParse(id), Description: sqlPackage.NullString{
		String: *newDescription,
		Valid:  newDescription != nil,
	}})
	return err == nil, err
}

func (r *projectMutationResolver) UpdateImagePriority(ctx context.Context, obj *model.ProjectMutation, id string, newPriority float64) (bool, error) {
	err := r.queries.UpdateImagePriority(context.Background(), sql.UpdateImagePriorityParams{ID: uuid.MustParse(id), Priority: float32(newPriority)})
	return err == nil, err
}

func (r *queryResolver) SearchProjects(ctx context.Context, searchString string, options model.SearchOptions, offset int, countLimit int) ([]*model.Project, error) {
	dbResults, err := r.queries.GetProjects(context.Background(), sql.GetProjectsParams{Limit: int32(countLimit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Project, len(dbResults))

	for i, dbProject := range dbResults {
		results[i] = &model.Project{
			ID:          dbProject.ID.String(),
			Name:        dbProject.Name,
			Description: dbProject.Description,
			CreatorID:   dbProject.Creator.String(),
			CreatedAt:   &dbProject.CreatedAt.Time,
			Languages:   []string{},
		}
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
		CreatorID:   dbProject.Creator.String(),
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
			CreatorID:   p.Creator.String(),
		}
	}
	return results, nil
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

// ProjectMutation returns generated.ProjectMutationResolver implementation.
func (r *Resolver) ProjectMutation() generated.ProjectMutationResolver {
	return &projectMutationResolver{r}
}

type projectResolver struct{ *Resolver }
type projectMutationResolver struct{ *Resolver }
