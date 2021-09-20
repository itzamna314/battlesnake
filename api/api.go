package api

import (
	"net/http"
)

func Build() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Index)
	mux.HandleFunc("/start", Start)
	mux.HandleFunc("/move", Move)
	mux.HandleFunc("/end", End)

	return mux
}
