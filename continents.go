package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http" // provides an HTTP client and server implementation
)

type Continent struct {
	Continent_id int `json:"id"`
	Name string `json:"name"`
}

type CreateContinent struct {
	Name string `json:"name"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Print("all good here!")
}

func Create(w http.ResponseWriter, r *http.Request) {
	// w is an instance of http.Request
	// it allows handler function to write
	// HTTP response back to client

	// r is instance of the http.Request struct
	// it represents the client's HTTP request
	// to the server
	var body CreateContinent
	err := json.NewDecoder(r.Body).Decode(&body)
	// NewDecoder returns a new decoder that reads from r
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request from body into CreateContinent struct: %v", err)
	}
	// below assigns err to result from query
	// returns error if var is nil
	if err := DB.QueryRow("INSERT INTO Continents (name) VALUES ($1)", body.Name).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("StatusInternalServerError %v", err)
		return
	}
	// WriteHeader used for building response to client
	w.WriteHeader(http.StatusCreated)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	var allContinents []Continent
	if err := DB.QueryRowContext(context.Background(),"SELECT * FROM Continents").Err(); err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error querying all continents: %v", err)
	}

	response, err := json.Marshal(allContinents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling continents into json %v", err)
		return
	}

	w.Write(response)
}



