package routes

import (
	"github.com/devjoemedia/scrumpilot-go-api/handlers"
	"github.com/go-chi/chi/v5"
)

func CommentRoute() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.GetComments)
	r.Get("/{id}", handlers.GetCommentByID)
	r.Post("/", handlers.CreateComment)
	r.Patch("/{id}", handlers.UpdateComment)
	r.Delete("/{id}", handlers.DeleteComment)

	return r
}
