package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresStorageIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:12345678@localhost/articleforum_test?sslmode=disable"
	}

	storage, err := NewPostgresStorage(dsn)
	require.NoError(t, err)

	err = clearDatabase(storage.db)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Create and get post", func(t *testing.T) {
		post, err := storage.CreatePost(ctx, "Integration Title", "Integration Content", true)
		require.NoError(t, err)
		require.NotEmpty(t, post.ID)

		retrievedPost, err := storage.GetPost(ctx, post.ID)
		require.NoError(t, err)
		assert.Equal(t, post.ID, retrievedPost.ID)
		assert.Equal(t, "Integration Title", retrievedPost.Title)
		assert.Equal(t, "Integration Content", retrievedPost.Content)
		assert.True(t, retrievedPost.CommentsEnabled)
	})

	t.Run("Get all posts", func(t *testing.T) {
		_, err := storage.CreatePost(ctx, "Post 1", "Content 1", true)
		require.NoError(t, err)
		_, err = storage.CreatePost(ctx, "Post 2", "Content 2", false)
		require.NoError(t, err)

		posts, err := storage.GetAllPosts(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(posts), 2)
	})

	t.Run("Create comment", func(t *testing.T) {
		post, err := storage.CreatePost(ctx, "For Comment", "Content", true)
		require.NoError(t, err)

		comment, err := storage.CreateComment(ctx, post.ID, nil, "Test Comment")
		require.NoError(t, err)
		require.NotNil(t, comment)
		assert.Equal(t, post.ID, comment.PostID)
		assert.Nil(t, comment.ParentID)
		assert.Equal(t, "Test Comment", comment.Content)

		nonExistentComment, err := storage.CreateComment(ctx, "non-existent", nil, "Comment")
		require.NoError(t, err)
		assert.Nil(t, nonExistentComment)
	})

	t.Run("Get comments with pagination", func(t *testing.T) {
		post, err := storage.CreatePost(ctx, "For Pagination", "Content", true)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			_, err := storage.CreateComment(ctx, post.ID, nil, "Comment")
			require.NoError(t, err)
		}

		comments, err := storage.GetComments(ctx, post.ID, 3, 0)
		require.NoError(t, err)
		assert.Len(t, comments, 3)

		commentsPage2, err := storage.GetComments(ctx, post.ID, 3, 3)
		require.NoError(t, err)
		assert.Len(t, commentsPage2, 2)
	})

	t.Run("Create comment to post with disabled comments", func(t *testing.T) {
		post, err := storage.CreatePost(ctx, "No Comments", "Content", false)
		require.NoError(t, err)

		comment, err := storage.CreateComment(ctx, post.ID, nil, "Should not work")
		require.NoError(t, err)
		assert.Nil(t, comment)
	})

	defer clearDatabase(storage.db)
}

func clearDatabase(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM comments")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM posts")
	return err
}
