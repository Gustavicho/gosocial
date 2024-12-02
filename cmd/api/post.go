package main

import (
	"net/http"

	"github.com/Gustavicho/gosocial/internal/store"
)

type CreatePostPayload struct {
	UserID  uint64   `json:"user_id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
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
