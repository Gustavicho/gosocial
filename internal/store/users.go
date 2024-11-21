package store

import (
	"context"
	"database/sql"
)

// model
type User struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatetAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (name, email, password) VALUES ($1, $2, $3)
	RETURNING id, created_at;
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatetAt,
	)
	if err != nil {
		return err
	}

	return nil
}
