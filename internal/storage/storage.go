package storage

import (
	"ArticleForum/internal/domain"
	"context"
)

type Storage interface {
	CreatePost(ctx context.Context, title, content string, commentsEnabled bool) (*domain.Post, error)
	GetPost(ctx context.Context, id string) (*domain.Post, error)
	GetAllPosts(ctx context.Context) ([]*domain.Post, error)
	CreateComment(ctx context.Context, postID string, parentID *string, content string) (*domain.Comment, error)
	GetComments(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error)
}
