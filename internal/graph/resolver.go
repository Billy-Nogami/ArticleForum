package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"ArticleForum/internal/graph/model"
	"ArticleForum/internal/storage/memory"
	"context"
)

type Resolver struct {
	storage *memory.MemoryStorage
}

func NewResolver() *Resolver {
	return &Resolver{
		storage: memory.NewMemoryStorage(),
	}
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, commentsEnabled bool) (*model.Post, error) {
	post, err := r.storage.CreatePost(ctx, title, content, commentsEnabled)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
	}, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*model.Comment, error) {
	comment, err := r.storage.CreateComment(ctx, postID, parentID, content)
	if err != nil {
		return nil, err
	}

	if comment == nil {
		return nil, nil // Пост не найден или комментарии запрещены
	}

	return &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.storage.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Post
	for _, post := range posts {
		result = append(result, &model.Post{
			ID:              post.ID,
			Title:           post.Title,
			Content:         post.Content,
			CommentsEnabled: post.CommentsEnabled,
			CreatedAt:       post.CreatedAt,
		})
	}
	return result, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	post, err := r.storage.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, nil
	}

	return &model.Post{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
	}, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) ([]*model.Comment, error) {
	actualLimit := 10
	if limit != nil {
		actualLimit = *limit
	}

	actualOffset := 0
	if offset != nil {
		actualOffset = *offset
	}

	comments, err := r.storage.GetComments(ctx, postID, actualLimit, actualOffset)
	if err != nil {
		return nil, err
	}

	var result []*model.Comment
	for _, comment := range comments {
		result = append(result, &model.Comment{
			ID:        comment.ID,
			PostID:    comment.PostID,
			ParentID:  comment.ParentID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		})
	}
	return result, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	ch := make(chan *model.Comment, 1)
	return ch, nil
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct {
	storage *memory.MemoryStorage
}
func NewResolver() *Resolver {
	return &Resolver{
		storage: memory.NewMemoryStorage(),
	}
}
*/
