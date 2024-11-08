package main

import (
	"log"
	"synth.com/file_converter/internal/router"
)

func main() {
	r := router.NewRouter() // Initialize router
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Unable to start the server:", err)
	}
}
