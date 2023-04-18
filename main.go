package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := OpenDatabase()
	if err != nil {
		log.Printf("error connecting to postgresql db: %v", err)
	}

	r := chi.NewRouter()
	// r.Use(middleware.Logger)

	// r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })

	r.Post("/continents", Create)

	r.Get("/", HealthCheck)

	r.Get("/continents", GetAll)
	
	http.ListenAndServe("localhost:5000", r)
}

