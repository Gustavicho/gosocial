package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Resource not found")
)

type Store struct {
	Posts interface {
		GetByID(context.Context, uint64) (*Post, error)
		Create(context.Context, *Post) error
		Delete(context.Context, uint64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		GetByID(context.Context, uint64) (*User, error)
		Create(context.Context, *User) error
	}
	Comments interface {
		GetByPostID(context.Context, uint64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Store {
	return Store{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
