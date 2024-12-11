package main

import (
	"net/http"

	"github.com/Gustavicho/gosocial/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromCtx(ctx)

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		// TODO: Get user dynamically
		UserID:  1,
		PostID:  post.ID,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: Redirect to post
	if err := jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
