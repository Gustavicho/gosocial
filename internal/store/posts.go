package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

// Model
type Post struct {
	ID        uint64   `json:"id"`
	UserID    uint64   `json:"user_id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (user_id, title, content, tags) VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at;
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
