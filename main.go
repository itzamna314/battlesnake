package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	timeout = flag.Int("latency", 200, "maximum time to wait before responding")
)

// Main Entrypoint
func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	handler := Handler()

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
