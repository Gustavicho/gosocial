package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// Model
type Post struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
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

func (s *PostStore) GetByID(ctx context.Context, id uint64) (*Post, error) {
	var post Post
	err := s.db.QueryRowContext(ctx, "SELECT * FROM posts WHERE id = $1", id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, id uint64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	return err
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts SET title = $1, content = $2, tags = $3, updated_at = now() WHERE id = $4;
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
	)

	return err
}
