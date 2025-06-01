package server

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/handlers"
	"github.com/kwinso/kubsu-web-api/static"
)

var routes = []Route{
	{
		Method:  "GET",
		Path:    "/",
		Handler: handlers.Index,
	},
	// Allow users that used POST to be redirected to index route without changing the method
	{
		Method:  "POST",
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
	mux.Handle("GET /static/{file...}", http.StripPrefix("/static", http.FileServer(http.FS(static.StaticFiles))))

	port := "8080"
	log.Println("Starting server on port", port)
	return http.ListenAndServe(":"+port, bootstrapedMux)

}
