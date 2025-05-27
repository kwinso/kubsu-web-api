package server

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/handlers"
)

var routes = []Route{
	{
		Method:  "GET",
		Path:    "/",
		Handler: handlers.Index,
	},
	{
		Method:     "POST",
		Path:       "/submissions",
		Controller: &handlers.SubmissionController{},
	},
}

func Serve() error {
	mux := http.NewServeMux()

	bootstrapedMux := bootstrapRouting(mux, routes)

	port := "8080"
	log.Println("Starting server on port", port)
	return http.ListenAndServe(":"+port, bootstrapedMux)
}
