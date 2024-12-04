package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gustavicho/gosocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	UserID  uint64   `json:"user_id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Add validation for payload

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
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the url
	postIDParam := chi.URLParam(r, "id")
	postID, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	// Get the post
	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusOK, post)
}
