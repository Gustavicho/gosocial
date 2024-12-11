package store

import (
	"context"
	"database/sql"
	"errors"
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

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

func (s *UserStore) GetByID(ctx context.Context, id uint64) (*User, error) {
	row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user := new(User)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatetAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
