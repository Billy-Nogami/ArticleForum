package memory

import (
	"ArticleForum/internal/domain"
	"ArticleForum/internal/storage"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	posts    map[string]*domain.Post
	comments map[string]*domain.Comment
	mu       sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    make(map[string]*domain.Post),
		comments: make(map[string]*domain.Comment),
	}
}

func (s *MemoryStorage) CreatePost(ctx context.Context, title, content string, commentsEnabled bool) (*domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := &domain.Post{
		ID:              uuid.New().String(),
		Title:           title,
		Content:         content,
		CommentsEnabled: commentsEnabled,
		CreatedAt:       time.Now(),
	}
	s.posts[post.ID] = post
	return post, nil
}

func (s *MemoryStorage) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[id]
	if !exists {
		return nil, nil
	}
	return post, nil
}

func (s *MemoryStorage) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*domain.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *MemoryStorage) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*domain.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[postID]
	if !exists {
		return nil, nil
	}

	if !post.CommentsEnabled {
		return nil, nil
	}

	comment := &domain.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.comments[comment.ID] = comment
	return comment, nil
}

func (s *MemoryStorage) GetComments(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var comments []*domain.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}

	if offset >= len(comments) {
		return []*domain.Comment{}, nil
	}

	end := offset + limit
	if end > len(comments) {
		end = len(comments)
	}

	return comments[offset:end], nil
}

var _ storage.Storage = (*MemoryStorage)(nil)
