package main

import (
	"log"
	"net/http"
	"os"

	"github.com/itzamna314/battlesnake/handlers"
)

// Main Entrypoint

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/start", handlers.Start)
	http.HandleFunc("/move", handlers.Move)
	http.HandleFunc("/end", handlers.End)

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
