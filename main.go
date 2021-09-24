package main

import (
	"log"
	"net/http"
	"os"
)

// Main Entrypoint
func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	handler := Handler()

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
