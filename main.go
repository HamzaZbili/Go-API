package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5" // router
	// "github.com/go-chi/chi/v5/middleware"
)

type Notification struct {
	Message string
}

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

	r.Get("/", HealthCheck)

	r.Post("/continents", CreateNewContinent)
	r.Get("/continents/all", GetAllContinents)
	r.Get("/continents", GetOneContinent)

	r.Post("/countries", CreateNewCountry)
	r.Get("/countries/all", GetAllCountries)
	r.Get("/countries", GetCountriesInContinent)
	
	http.ListenAndServe("localhost:5000", r)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// convention
	// w - response writer 
	// r - request
	var rootMessage Notification
	rootMessage.Message = "all good here!"
	response, err := json.Marshal(rootMessage.Message)
	// Marshal() takes a Go value as an input and
	// returns a byte slice and an error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling welcome message into json %v", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}