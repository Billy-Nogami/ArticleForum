package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"ArticleForum/internal/graph/model"
	"context"
)

type Resolver struct{}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsEnabled bool) (*model.Post, error) {
	panic("not implemented")
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*model.Comment, error) {
	panic("not implemented")
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	panic("not implemented")
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	panic("not implemented")
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, limit int, offset int) ([]*model.Comment, error) {
	panic("not implemented")
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
