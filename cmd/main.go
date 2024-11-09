package main

import (
	"log"
	"synth.com/file_converter/internal/router"
)

func main() {
	r := router.NewRouter()

	// Start the server
	err := r.Run(":8080") // Run on port 8080
	if err != nil {
		log.Fatal("Unable to start the server:", err)
	}
}
