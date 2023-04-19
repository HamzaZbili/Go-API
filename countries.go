package main

import (
	"encoding/json"
	"log"
	"net/http" // provides an HTTP client and server implementation
)

type Country struct {
	Country int `json:"id"`
	Name string `json:"name"`
	Poplation int `json:"population"`
	Capital string `json:"capital"`
	Continent int `json:"continent_id"`
}

type CreateCountry struct {
	Name string `json:"name"`
	Poplation int `json:"population"`
	Capital string `json:"capital"`
	Continent int `json:"continent"`
}

func CreateNewCountry (w http.ResponseWriter, r *http.Request) {
	var body CreateCountry
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request from body into CreateCountry struct: %v", err)
		return
	}
	if err := DB.QueryRow("INSERT INTO Countries (Name, Population, Capital, Continent) VALUES ($1, $2, $3, $4);",
	body.Name, body.Poplation, body.Capital, body.Continent).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("StatusInternalServerError %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}