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
	"gitlab.lrz.de/projecthub/gql-api/sqlc"
)

func (r *mutationResolver) CreateProject(ctx context.Context, project *model.NewProjectInput) (*model.Project, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	projectID, err := r.queries.CreateProject(ctx, sqlc.CreateProjectParams{Name: project.Name, Description: project.Description, Creator: uuid.MustParse(user.Id)})
	if err != nil {
		return nil, err
	}
	fetchedProject, err := r.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return model.ProjectFromDBProject(fetchedProject), nil
}

func (r *mutationResolver) AddSavedProject(ctx context.Context, id string) (bool, error) {
	user, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = r.queries.SaveProjectForUser(context.Background(), sqlc.SaveProjectForUserParams{UserID: uuid.MustParse(user.Id), ProjectID: uuid.MustParse(id)})
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
	err = r.queries.UnsaveProjectForUser(context.Background(), sqlc.UnsaveProjectForUserParams{UserID: uuid.MustParse(user.Id), ProjectID: uuid.MustParse(id)})
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
		participants[i] = model.DBUserToUser(p)
	}
	return participants, nil
}

func (r *projectResolver) Creator(ctx context.Context, obj *model.Project) (*model.User, error) {
	dbUser, err := r.queries.GetUserByID(context.Background(), uuid.MustParse(obj.CreatorID))
	if err != nil {
		return nil, err
	}
	return model.DBUserToUser(dbUser), nil
}

func (r *projectResolver) Images(ctx context.Context, obj *model.Project) ([]*model.Image, error) {
	images, err := r.queries.GetImagesOfProject(ctx, uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	toReturn := make([]*model.Image, len(images))
	for i, img := range images {
		description := img.Description.String
		toReturn[i] = &model.Image{
			ID:          img.ID.String(),
			URL:         img.Url,
			CreatedAt:   img.CreatedAt.Time,
			Description: &description,
			Priority:    float64(img.Priority),
		}
	}
	return toReturn, nil
}

func (r *projectResolver) Saved(ctx context.Context, obj *model.Project) (bool, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return false, nil
	}
	savedProjects, err := r.queries.GetSavedProjectsForUser(context.Background(), uuid.MustParse(me.Id))
	if err != nil {
		return false, err
	}
	for _, savedProject := range savedProjects {
		if savedProject.ID.String() == obj.ID {
			return true, nil
		}
	}
	return false, nil
}

func (r *projectResolver) Tags(ctx context.Context, obj *model.Project) ([]string, error) {
	tags, err := r.queries.GetProjectTags(ctx, uuid.MustParse(obj.ID))
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *projectMutationResolver) AddParticipant(ctx context.Context, obj *model.ProjectMutation, id string) (bool, error) {
	err := r.queries.AddParticipantToProject(context.Background(), sqlc.AddParticipantToProjectParams{ProjectID: uuid.MustParse(obj.ID), UserID: uuid.MustParse(id)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *projectMutationResolver) RemoveParticipant(ctx context.Context, obj *model.ProjectMutation, id string) (bool, error) {
	err := r.queries.RemoveParticipantFromProject(context.Background(), sqlc.RemoveParticipantFromProjectParams{ProjectID: uuid.MustParse(obj.ID), UserID: uuid.MustParse(id)})
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
	err := r.queries.UpdateProjectDescription(context.Background(), sqlc.UpdateProjectDescriptionParams{ID: uuid.MustParse(obj.ID), Description: newDescription})
	return err == nil, err
}

func (r *projectMutationResolver) UpdateName(ctx context.Context, obj *model.ProjectMutation, newName string) (bool, error) {
	err := r.queries.UpdateProjectName(context.Background(), sqlc.UpdateProjectNameParams{ID: uuid.MustParse(obj.ID), Name: newName})
	return err == nil, err
}

func (r *projectMutationResolver) AddImage(ctx context.Context, obj *model.ProjectMutation, newImage model.NewImageInput) (bool, error) {
	priority := 0.0
	if newImage.Priority != nil {
		priority = *newImage.Priority
	}
	description := sqlPackage.NullString{}
	if newImage.Description != nil {
		description = sqlPackage.NullString{String: *newImage.Description, Valid: true}
	}
	err := r.queries.AddImageToProject(context.Background(), sqlc.AddImageToProjectParams{
		Project:     uuid.MustParse(obj.ID),
		Url:         newImage.URL,
		Description: description,
		Priority:    float32(priority),
	})
	return err == nil, err
}

func (r *projectMutationResolver) RemoveImage(ctx context.Context, obj *model.ProjectMutation, id string) (bool, error) {
	err := r.queries.DeleteImage(context.Background(), uuid.MustParse(id))
	return err == nil, err
}

func (r *projectMutationResolver) UpdateImageDescription(ctx context.Context, obj *model.ProjectMutation, id string, newDescription *string) (bool, error) {
	err := r.queries.UpdateImageDescription(context.Background(), sqlc.UpdateImageDescriptionParams{ID: uuid.MustParse(id), Description: sqlPackage.NullString{
		String: *newDescription,
		Valid:  newDescription != nil,
	}})
	return err == nil, err
}

func (r *projectMutationResolver) UpdateImagePriority(ctx context.Context, obj *model.ProjectMutation, id string, newPriority float64) (bool, error) {
	err := r.queries.UpdateImagePriority(context.Background(), sqlc.UpdateImagePriorityParams{ID: uuid.MustParse(id), Priority: float32(newPriority)})
	return err == nil, err
}

func (r *projectMutationResolver) AddTag(ctx context.Context, obj *model.ProjectMutation, tag string) (bool, error) {
	err := r.queries.TagProject(context.Background(), sqlc.TagProjectParams{
		Name:      tag,
		ProjectID: uuid.MustParse(obj.ID),
	})
	return err == nil, err
}

func (r *projectMutationResolver) RemoveTag(ctx context.Context, obj *model.ProjectMutation, tag string) (bool, error) {
	err := r.queries.UntagProject(context.Background(), sqlc.UntagProjectParams{
		Name:      tag,
		ProjectID: uuid.MustParse(obj.ID),
	})
	return err == nil, err
}

func (r *queryResolver) SearchProjects(ctx context.Context, searchString string, options model.SearchOptions, offset int, countLimit int) ([]*model.Project, error) {
	//dbResults, err := r.queries.GetProjects(context.Background(), sqlc.GetProjectsParams{Limit: int32(countLimit), Offset: int32(offset)})
	var dbResults []sqlc.Project
	var err error
	if options.Tag != nil {
		// Tag name comes first, then search string name
		dbResults, err = r.queries.SearchProjectsWithTag(context.Background(), sqlc.SearchProjectsWithTagParams{
			Limit:  int32(countLimit),
			Offset: int32(offset),
			Name:   *options.Tag,
			Name_2: "%" + searchString + "%",
		})
	} else {
		dbResults, err = r.queries.SearchProjects(context.Background(), sqlc.SearchProjectsParams{
			Limit:  int32(countLimit),
			Offset: int32(offset),
			Name:   "%" + searchString + "%",
		})
	}
	if err != nil {
		return nil, err
	}
	results := model.ProjectsFromDBProjects(dbResults)
	return results, nil
}

func (r *queryResolver) GetProject(ctx context.Context, id string) (*model.Project, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	dbProject, err := r.queries.GetProjectByID(context.Background(), uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return model.ProjectFromDBProject(dbProject), nil
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
	results := model.ProjectsFromDBProjects(dbResults)
	return results, nil
}

func (r *queryResolver) AvailableTags(ctx context.Context, limit int, offset int) ([]*model.Tag, error) {
	dbRes, err := r.queries.GetTagsByCount(context.Background(), sqlc.GetTagsByCountParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	res := make([]*model.Tag, len(dbRes))
	for i, dbTag := range dbRes {
		res[i] = &model.Tag{
			Name:  dbTag.Name,
			Count: int(dbTag.Count),
		}
	}
	return res, nil
}

// Project returns generated.ProjectResolver implementation.
func (r *Resolver) Project() generated.ProjectResolver { return &projectResolver{r} }

// ProjectMutation returns generated.ProjectMutationResolver implementation.
func (r *Resolver) ProjectMutation() generated.ProjectMutationResolver {
	return &projectMutationResolver{r}
}

type projectResolver struct{ *Resolver }
type projectMutationResolver struct{ *Resolver }
