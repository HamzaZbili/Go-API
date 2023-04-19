package main

import (
	"encoding/json"
	"fmt"
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
	Population int `json:"population"`
	Capital string `json:"capital"`
	Continent int `json:"continent"`
}

func CreateNewCountry (w http.ResponseWriter, r *http.Request) {
	var body CreateCountry
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("error decoding request from body into CreateCountry struct: %v", err)
		return
	}
	if err := DB.QueryRow(`
	INSERT INTO Countries (Name, Population, Capital, Continent)
	VALUES ($1, $2, $3, $4);`,
	body.Name, body.Population, body.Capital, body.Continent).Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("StatusInternalServerError %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetAllCountries(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM Countries;")
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error querying all countries: %v", err)
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
			fmt.Printf("Error looping countries row: %v", err)
		}
		countriesSlice = append(countriesSlice, country)
	}
	allCountries, err := json.Marshal(countriesSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error marshalling countries into json %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(allCountries)
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
`, continentName) // INNER JOIN clause combines rows from two or more tables
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error querying countries by continent name: %v", err)
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
			fmt.Printf("Error looping countries row: %v", err)
		}
		countriesSlice = append(countriesSlice, country)
	}
	continentsToSend, err := json.Marshal(countriesSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error marshalling countries into json %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(continentsToSend)
}

func UpdateCountry(w http.ResponseWriter, r *http.Request){
	countryName := r.URL.Query().Get("country")
	var body CreateCountry
	// parse request body to extract fields
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("error decoding request from body into UpdateCountry struct: %v", err)
		return
	}
// DW.Exec() returns values: an sql. Result and an error. When the error is nil
// you can use the Result to get the ID of the last inserted item
// or to retrieve the number of rows affected by the operation
	result, err := DB.Exec(`
		UPDATE Countries
		SET
		Population = $1,
		Capital = $2,
		Continent = $3
		WHERE Name = $4`,
        body.Population, body.Capital, body.Continent, countryName)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Printf("error updating country: %v", err)
        return
    }
// check how many rows were affected by the update
	rowsAffected, err := result.RowsAffected()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Printf("error checking rows affected by country update: %v", err)
        return
    }

	if rowsAffected == 0 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Country with name %s not found", countryName)
        return
    }
	w.WriteHeader(http.StatusOK)
    fmt.Printf("country updated: %v", countryName)
}

func DeleteCountry(w http.ResponseWriter, r *http.Request){
	id := r.URL.Query().Get("id")
	    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Println("query parameter id missing")
        return
    }
	result, err := DB.Exec(
		`DELETE FROM Countries
		WHERE Country_id = $1;`, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("StatusInternalServerError %v", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Printf("error getting rows affected by country delete: %v", err)
        return
    }

	if rowsAffected == 0 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Printf("Country with id %s not found", id)
        return
    }

	w.WriteHeader(http.StatusOK)
	fmt.Printf("country deleted: %v", id)
}