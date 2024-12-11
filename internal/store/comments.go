package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	PostID    uint64 `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID uint64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c."content", c.created_at, c.updated_at, u.username, u.id FROM "comments" c 
			JOIN users u 
			ON u.id = c.user_id
			WHERE c.post_id = $1
			ORDER BY c.created_at DESC 
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}

		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.User.Name, &c.User.ID,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
