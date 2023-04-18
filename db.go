package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func OpenDatabase() error {
	var err error
	DB, err = sql.Open("postgres", "user=hamzazbili dbname=postgresql_practice sslmode=disable")
	if err != nil {
		return err
	}
	fmt.Println("connection to postgresql db successful")
	return nil
}

func CloseDatabase() error {
	return DB.Close()
}