package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devjoemedia/scrumpilot-go-api/config"
	"github.com/devjoemedia/scrumpilot-go-api/database"
	_ "github.com/devjoemedia/scrumpilot-go-api/docs"
	"github.com/devjoemedia/scrumpilot-go-api/middleware"
	"github.com/devjoemedia/scrumpilot-go-api/routes"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Example API
// @version 1.0
// @description API documentation

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load Envs
	config.Load()

	// Connect to DB
	database.Connect()

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))
	r.Mount("/api/v1/auth", routes.AuthRoute())

	// Protected routes — middleware applied at mount point
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Mount("/api/v1/todos", routes.TodoRoute())
		r.Mount("/api/v1/users", routes.UserRoutes())
		r.Mount("/api/v1/tickets", routes.TicketRoute())
		r.Mount("/api/v1/comments", routes.CommentRoute())
	})

	port := ":" + config.AppConfig.AppPort
	fmt.Printf("🚀 Server running on port %s\n", port)
	fmt.Println("📚 Docs:: http://localhost:8000/swagger/index.html")

	// This line blocks the program from exiting
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
