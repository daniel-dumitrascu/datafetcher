package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 || len(args[1]) == 0 {
		fmt.Println("Please provide a valid path to the directory where the mapping files are stored")
		return
	}

	log.Println("Starting the server...")

	if err := Init(args[1]); err != nil {
		log.Fatalf("Error when initiating the server: %v\n", err)
	}

	log.Println("Server has started successfully")

	Run()
}
