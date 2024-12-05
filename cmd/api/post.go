package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gustavicho/gosocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=255"`
	Content string   `json:"content" validate:"required,max=255"`
	Tags    []string `json:"tags"`
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

	writeJSON(w, http.StatusCreated, post)
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the url
	postIDParam := chi.URLParam(r, "id")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	// Get the post
	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	writeJSON(w, http.StatusOK, post)
}
