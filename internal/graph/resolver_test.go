package graph

import (
	"ArticleForum/internal/domain"
	"ArticleForum/internal/storage/mock"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolverWithMocks(t *testing.T) {

	mockStorage := new(mock.MockStorage)
	resolver := NewResolver(mockStorage)

	t.Run("CreatePost with mock", func(t *testing.T) {

		expectedPost := &domain.Post{
			ID:              "1",
			Title:           "Test Title",
			Content:         "Test Content",
			CommentsEnabled: true,
			CreatedAt:       time.Now(),
		}
		mockStorage.On("CreatePost", context.Background(), "Test Title", "Test Content", true).Return(expectedPost, nil)

		post, err := resolver.Mutation().CreatePost(
			context.Background(),
			"Test Title",
			"Test Content",
			true,
		)

		require.NoError(t, err)
		assert.Equal(t, "1", post.ID)
		assert.Equal(t, "Test Title", post.Title)
		assert.Equal(t, "Test Content", post.Content)
		assert.True(t, post.CommentsEnabled)

		mockStorage.AssertExpectations(t)
	})

	t.Run("GetPosts with mock", func(t *testing.T) {
		expectedPosts := []*domain.Post{
			{
				ID:              "1",
				Title:           "Test Title 1",
				Content:         "Test Content 1",
				CommentsEnabled: true,
				CreatedAt:       time.Now(),
			},
			{
				ID:              "2",
				Title:           "Test Title 2",
				Content:         "Test Content 2",
				CommentsEnabled: false,
				CreatedAt:       time.Now(),
			},
		}
		mockStorage.On("GetAllPosts", context.Background()).Return(expectedPosts, nil)

		posts, err := resolver.Query().Posts(context.Background())
		require.NoError(t, err)
		assert.Len(t, posts, 2)
		assert.Equal(t, "1", posts[0].ID)
		assert.Equal(t, "Test Title 1", posts[0].Title)

		mockStorage.AssertExpectations(t)
	})

	t.Run("CreateComment with mock", func(t *testing.T) {

		expectedComment := &domain.Comment{
			ID:        "comment-1",
			PostID:    "post-1",
			ParentID:  nil,
			Content:   "Test Comment",
			CreatedAt: time.Now(),
		}
		mockStorage.On("CreateComment", context.Background(), "post-1", (*string)(nil), "Test Comment").Return(expectedComment, nil)

		comment, err := resolver.Mutation().CreateComment(
			context.Background(),
			"post-1",
			nil,
			"Test Comment",
		)

		require.NoError(t, err)
		assert.Equal(t, "comment-1", comment.ID)
		assert.Equal(t, "post-1", comment.PostID)
		assert.Equal(t, "Test Comment", comment.Content)
		assert.Nil(t, comment.ParentID)

		mockStorage.AssertExpectations(t)
	})

	t.Run("GetComments with mock", func(t *testing.T) {
		expectedComments := []*domain.Comment{
			{
				ID:        "comment-1",
				PostID:    "post-1",
				ParentID:  nil,
				Content:   "Comment 1",
				CreatedAt: time.Now(),
			},
			{
				ID:        "comment-2",
				PostID:    "post-1",
				ParentID:  nil,
				Content:   "Comment 2",
				CreatedAt: time.Now(),
			},
		}
		mockStorage.On("GetComments", context.Background(), "post-1", 10, 0).Return(expectedComments, nil)

		comments, err := resolver.Query().Comments(
			context.Background(),
			"post-1",
			nil,
			nil,
		)

		require.NoError(t, err)
		assert.Len(t, comments, 2)
		assert.Equal(t, "comment-1", comments[0].ID)
		assert.Equal(t, "Comment 1", comments[0].Content)

		mockStorage.AssertExpectations(t)
	})
}
