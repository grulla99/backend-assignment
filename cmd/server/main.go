package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"backed-assignment-aro/internal/handler"
	"backed-assignment-aro/internal/repository"
	"backed-assignment-aro/internal/service"
)

func main() {
	repo := repository.NewMemoryRepository()
	svc := service.NewIssueService(repo)
	h := handler.NewHTTPHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	h.Register(r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
