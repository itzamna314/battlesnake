package main

import (
	"net/http"
)

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/start", Start)
	mux.HandleFunc("/move", Move)
	mux.HandleFunc("/end", End)

	return mux
}
