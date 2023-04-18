package main

import (
	"log"
)

func main() {
	err := OpenDatabase()
	if err != nil {
		log.Printf("error connecting to postgresql db: %v", err)
	}
	defer CloseDatabase()
}