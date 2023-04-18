package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Continent struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

type CreateContinent struct {
	Name string `json:"name"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	var body CreateContinent
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request from body into CreateContinent struct: %v", err)
	}

	if err := DB.QueryRow("INSERT INTO continents (name) VALUES ($1)", body.Name).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("StatusInternalServerError %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

