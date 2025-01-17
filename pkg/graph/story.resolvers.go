package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"

	"github.com/google/uuid"
	"github.com/navikt/nada-backend/pkg/auth"
	"github.com/navikt/nada-backend/pkg/graph/generated"
	"github.com/navikt/nada-backend/pkg/graph/models"
)

// CreateStory is the resolver for the createStory field.
func (r *mutationResolver) CreateStory(ctx context.Context, files []*models.UploadFile, input models.NewStory) (*models.Story, error) {
	story, err := r.repo.CreateStory(ctx, auth.GetUser(ctx).Email, input)
	if err != nil {
		return nil, err
	}

	if err = WriteFilesToBucket(ctx, story.ID.String(), files); err != nil {
		return nil, err
	}

	// Create a new File object with the uploaded file's public URL
	return story, nil
}

// UpdateStoryMetadata is the resolver for the updateStoryMetadata field.
func (r *mutationResolver) UpdateStoryMetadata(ctx context.Context, id uuid.UUID, name string, description string, keywords []string, teamkatalogenURL *string, productAreaID *string, teamID *string, group string) (*models.Story, error) {
	existing, err := r.repo.GetStory(ctx, id)
	if err != nil {
		return nil, err
	}

	user := auth.GetUser(ctx)
	if !user.GoogleGroups.Contains(existing.Group) {
		return nil, ErrUnauthorized
	}

	story, err := r.repo.UpdateStoryMetadata(ctx, id, name, description, keywords, teamkatalogenURL,
		productAreaID, teamID, group)
	if err != nil {
		return nil, err
	}

	return story, nil
}

// DeleteStory is the resolver for the deleteStory field.
func (r *mutationResolver) DeleteStory(ctx context.Context, id uuid.UUID) (bool, error) {
	s, err := r.repo.GetStory(ctx, id)
	if err != nil {
		return false, err
	}

	user := auth.GetUser(ctx)
	if !user.GoogleGroups.Contains(s.Group) {
		return false, ErrUnauthorized
	}

	if err = r.repo.DeleteStory(ctx, id); err != nil {
		return false, err
	}

	if err := deleteStoryFolder(ctx, id.String()); err != nil {
		r.log.WithError(err).
			Errorf("Data story %v metadata deleted but failed to delete story files in GCP",
				id)
		return false, err
	}

	return true, nil
}

// DataStory is the resolver for the dataStory field.
func (r *queryResolver) DataStory(ctx context.Context, id uuid.UUID) (*models.Story, error) {
	return r.repo.GetStory(ctx, id)
}

// ProductAreaID is the resolver for the productAreaID field.
func (r *storyResolver) ProductAreaID(ctx context.Context, obj *models.Story) (*string, error) {
	if teamID := ptrToString(obj.TeamID); teamID != "" {
		team, err := r.teamkatalogen.GetTeam(ctx, teamID)
		if err != nil {
			return nil, err
		}
		return &team.ProductAreaID, nil
	}

	return nil, nil
}

// Story returns generated.StoryResolver implementation.
func (r *Resolver) Story() generated.StoryResolver { return &storyResolver{r} }

type storyResolver struct{ *Resolver }
