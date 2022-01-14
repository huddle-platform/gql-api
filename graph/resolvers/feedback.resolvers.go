package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.lrz.de/projecthub/gql-api/auth"
	"gitlab.lrz.de/projecthub/gql-api/graph/generated"
	"gitlab.lrz.de/projecthub/gql-api/graph/model"
)

func (r *feedbackMutationResolver) SubmitImageNotWorking(ctx context.Context, obj *model.FeedbackMutation, imageID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Feedback(ctx context.Context) (*model.FeedbackMutation, error) {
	me, err := auth.IdentityFromContext(ctx)
	if err != nil {
		return &model.FeedbackMutation{}, nil
	}
	return &model.FeedbackMutation{
		UserId: &me.Id,
	}, nil
}

// FeedbackMutation returns generated.FeedbackMutationResolver implementation.
func (r *Resolver) FeedbackMutation() generated.FeedbackMutationResolver {
	return &feedbackMutationResolver{r}
}

type feedbackMutationResolver struct{ *Resolver }
