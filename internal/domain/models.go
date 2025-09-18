package domain

import "time"

type Post struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postID"`
	ParentID  *string   `json:"parentID"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
