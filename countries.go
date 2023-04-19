package main

import (
	"encoding/json"
	"log"
	"net/http" // provides an HTTP client and server implementation
)

type Country struct {
	Country_id int `json:"id"`
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

func GetAllCountries(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM Countries;")
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error querying all countries: %v", err)
		return
	}
	defer rows.Close()

	var countriesSlice []Country

	for rows.Next() {
		var country Country
		if err := rows.Scan(&country.Country_id,
			&country.Name, &country.Poplation, &country.Capital,
			&country.Continent); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error looping countries row: %v", err)
		}
		countriesSlice = append(countriesSlice, country)
	}
	continentsToSend, err := json.Marshal(countriesSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling countries into json %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(continentsToSend)
}

func GetCountriesInContinent(w http.ResponseWriter, r *http.Request) {
	continentName := r.URL.Query().Get("name")
	rows, err := DB.Query(
	// c and ct in the SELECT statement are used an alias for the Countries and Continents tables
	// using * here can make SQL queries less efficient and harder to maintain.
		`SELECT c.Country_id, c.Name, c.Population, c.Capital, ct.Continent_id
		FROM Countries c
		INNER JOIN Continents ct ON c.Continent = ct.Continent_id
		WHERE ct.Name = $1
`, continentName)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error querying countries by continent name: %v", err)
		return
	}
	defer rows.Close()

	var countriesSlice []Country

	for rows.Next() {
		var country Country
		if err := rows.Scan(&country.Country_id,
			&country.Name, &country.Poplation, &country.Capital,
			&country.Continent); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error looping countries row: %v", err)
		}
		countriesSlice = append(countriesSlice, country)
	}
	continentsToSend, err := json.Marshal(countriesSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling countries into json %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(continentsToSend)
}

