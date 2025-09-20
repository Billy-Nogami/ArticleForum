package postgres

import (
	"ArticleForum/internal/domain"
	"ArticleForum/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dataSourceName string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return &PostgresStorage{db: db}, nil
}

func createTables(db *sql.DB) error {
	postsTable := `
		CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			comments_enabled BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`

	commentsTable := `
		CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			parent_id TEXT REFERENCES comments(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`

	indexes := `
		CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
		CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
	`

	if _, err := db.Exec(postsTable); err != nil {
		return err
	}
	if _, err := db.Exec(commentsTable); err != nil {
		return err
	}
	if _, err := db.Exec(indexes); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) CreatePost(ctx context.Context, title, content string, commentsEnabled bool) (*domain.Post, error) {
	id := uuid.New().String()
	createdAt := time.Now()
	query := `INSERT INTO posts (id, title, content, comments_enabled, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query, id, title, content, commentsEnabled, createdAt)
	if err != nil {
		return nil, err
	}
	return &domain.Post{
		ID:              id,
		Title:           title,
		Content:         content,
		CommentsEnabled: commentsEnabled,
		CreatedAt:       createdAt,
	}, nil
}

func (s *PostgresStorage) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	query := `SELECT id, title, content, comments_enabled, created_at FROM posts WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)
	var post domain.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostgresStorage) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	query := `SELECT id, title, content, comments_enabled, created_at FROM posts ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsEnabled, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, postID string, parentID *string, content string) (*domain.Comment, error) {
	// Проверяем, существует ли пост и разрешены ли комментарии
	post, err := s.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil || !post.CommentsEnabled {
		return nil, nil
	}

	id := uuid.New().String()
	createdAt := time.Now()
	var query string
	if parentID == nil {
		query = `INSERT INTO comments (id, post_id, content, created_at) VALUES ($1, $2, $3, $4)`
		_, err = s.db.ExecContext(ctx, query, id, postID, content, createdAt)
	} else {
		query = `INSERT INTO comments (id, post_id, parent_id, content, created_at) VALUES ($1, $2, $3, $4, $5)`
		_, err = s.db.ExecContext(ctx, query, id, postID, *parentID, content, createdAt)
	}
	if err != nil {
		return nil, err
	}
	return &domain.Comment{
		ID:        id,
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: createdAt,
	}, nil
}

func (s *PostgresStorage) GetComments(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
	query := `SELECT id, post_id, parent_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`
	rows, err := s.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var comment domain.Comment
		var parentID sql.NullString
		if err := rows.Scan(&comment.ID, &comment.PostID, &parentID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			comment.ParentID = &parentID.String
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

var _ storage.Storage = (*PostgresStorage)(nil)
