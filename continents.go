package main

import (
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
	// convention
	// w - response writer 
	// r - for request
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
		return
	}
	// below assigns err to result from query
	// returns error if var is nil
	if err := DB.QueryRow("INSERT INTO Continents (name) VALUES ($1);", body.Name).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("StatusInternalServerError %v", err)
		return
	}
	// WriteHeader used for building response to client
	w.WriteHeader(http.StatusCreated)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM Continents;")
	// queries all rows in continents
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error querying all continents: %v", err)
		return
	}
	defer rows.Close() // releases any resources held by the rows
	// rows returns continents as slice

	var continents []Continent

	for rows.Next() {
	// Next() iterates over rows
		var continent Continent
		if err := rows.Scan(&continent.Continent_id, &continent.Name); err != nil {
		// Scan() reads values of current row
		// receives pointers to where data will be passed to
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error loopin continents row: %v", err)
		}
		// appends to continents array
		continents = append(continents, continent)
	}
	response, err := json.Marshal(continents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling continents into json %v", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}



