package mock

import (
	"ArticleForum/internal/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreatePost(ctx context.Context, title, content string, commentsEnabled bool) (*domain.Post, error) {
	args := m.Called(ctx, title, content, commentsEnabled)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockStorage) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockStorage) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Post), args.Error(1)
}

func (m *MockStorage) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*domain.Comment, error) {
	args := m.Called(ctx, postID, parentID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockStorage) GetComments(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
	args := m.Called(ctx, postID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Comment), args.Error(1)
}
