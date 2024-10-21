package main

import (
	"log"
)

func main() {
	log.Println("Starting the server...")

	if err := Init(); err != nil {
		log.Fatalf("Error when initiating the server: %v\n", err)
	}

	log.Println("Server has started successfully")

	Run()
}
