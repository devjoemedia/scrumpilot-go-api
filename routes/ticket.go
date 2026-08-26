package routes

import (
	"github.com/devjoemedia/scrumpilot-go-api/handlers"
	"github.com/go-chi/chi/v5"
)

func TicketRoute() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handlers.GetTickets)
	r.Get("/{id}", handlers.GetTicketByID)
	r.Get("/{id}/comments", handlers.GetComments)
	r.Post("/", handlers.CreateTicket)
	r.Patch("/{id}/assign", handlers.AssignTicket)
	r.Patch("/{id}", handlers.UpdateTicket)
	r.Delete("/{id}", handlers.DeleteTicket)

	return r
}
