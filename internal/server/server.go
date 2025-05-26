package server

import (
	"log"
	"net/http"
)

func Serve() error {
	mux := http.NewServeMux()

	RegisterRoutes(mux)

	port := "8080"
	log.Println("Starting server on port", port)
	return http.ListenAndServe(":"+port, mux)
}
