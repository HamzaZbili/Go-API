package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http" // provides an HTTP client and server implementation
)

type Continent struct {
	Continent_id int `json:"id"`
	Name string `json:"name"`
}

type CreateQuery struct {
	Name string `json:"name"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	// w is an instance of http.Request
	// it allows handler function to write
	// HTTP response back to client

	// r is instance of the http.Request struct
	// it represents the client's HTTP request
	// to the server
	var body CreateQuery
	err := json.NewDecoder(r.Body).Decode(&body)
	// NewDecoder returns a new decoder that reads from r
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request from body into CreateQuery struct: %v", err)
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

	var continentsSlice []Continent

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
		continentsSlice = append(continentsSlice, continent)
	}
	continentsToSend, err := json.Marshal(continentsSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling continents into json %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(continentsToSend)
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	// var continent Continent
	queryResult := DB.QueryRow("SELECT Continent_id, name FROM Continents WHERE Continent_id = $1;", id)
	var continent Continent
	if err := queryResult.Scan(&continent.Continent_id, &continent.Name); err == sql.ErrNoRows {
		// sql.ErrNoRows returned from Scan()
		fmt.Println("No rows found")
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error scanning continent row: %v", err)
	}
	response, err := json.Marshal(continent.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling Continent into json %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}



