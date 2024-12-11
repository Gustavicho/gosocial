package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Gustavicho/gosocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type CtxKey string

const PostCtxKey CtxKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=255"`
	Content string   `json:"content" validate:"required,max=255"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=255"`
	Content *string   `json:"content" validate:"omitempty,max=255"`
	Tags    *[]string `json:"tags" validate:"omitempty"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		// TODO: Get user dynamically
		UserID:  1,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	err := app.store.Posts.Create(ctx, post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromCtx(ctx)

	coments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = coments

	if err = jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromCtx(ctx)

	err := app.store.Posts.Delete(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusCreated, map[string]string{
		"message": "Deleted with success!",
	}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromCtx(ctx)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// if nil, the field will not be updated
	// in other words. if the field is not present
	// in the request, it will not be updated
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	err := app.store.Posts.Update(ctx, post)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err = jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getPostFromCtx(ctx context.Context) *store.Post {
	return ctx.Value(PostCtxKey).(*store.Post)
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDParam := chi.URLParam(r, "id")
		postID, err := strconv.ParseUint(postIDParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(r.Context(), postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, PostCtxKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
