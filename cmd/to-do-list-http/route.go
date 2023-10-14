package main

import (
	"net/http"
	"to-do-list/internal/handler/http/task"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/runtime/middleware"
)

func newRoutes(task *task.Handler) *chi.Mux {
	myRouter := chi.NewRouter()
	myRouter.Get("/api/tasks", task.GetAll)
	myRouter.Post("/api/task", task.Create)
	myRouter.Put("/api/task/{id}", task.Update)
	myRouter.Delete("/api/task/{id}", task.Delete)

	myRouter.Handle("/docs.yaml", http.FileServer(http.Dir("./docs")))
	opts := middleware.SwaggerUIOpts{SpecURL: "docs.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	myRouter.Handle("/docs", sh)
	return myRouter

}
