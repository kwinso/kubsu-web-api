package server

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/handlers"
	"github.com/kwinso/kubsu-web-api/internal/server/middleware"
)

type Route struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

// TODO:
// 1. Route to serve main HTML page
// 2. Route to handle POST request to create a new submission
// 2.1 Should supoprt JSON & form data
// 2.2 Should validate the request
// 2.3 Should set username and password to cookies
// 3. When getting main page, it should check cookies for username and password, and render the old submission if found
// 4. Support updates with POST by ID (for form data)
// 5. Support updates with PUT by ID (for JSON)

var routes = []Route{
	{
		Method:  "GET",
		Path:    "/",
		Handler: handlers.Index,
	},
	{
		Method:  "POST",
		Path:    "/submissions",
		Handler: handlers.CreateSubmission,
	},
}

func bootstrapRouting(mux *http.ServeMux) http.Handler {
	for _, route := range routes {
		log.Println("Registering route:", route.Method, route.Path)
		mux.HandleFunc(route.Path, route.Handler)
	}

	return middleware.NewMimeTypeEnforce(mux, "application/json", "multipart/form-data", "text/html")
}
