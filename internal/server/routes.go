package server

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/middleware"
)

type Controller interface {
	Boot(mux *http.ServeMux) *http.ServeMux
}

type Route struct {
	Method     string
	Path       string
	Handler    func(http.ResponseWriter, *http.Request)
	Controller Controller
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

func bootstrapRouting(mux *http.ServeMux, routes []Route) http.Handler {
	for _, route := range routes {
		log.Println("Registering route:", route.Method, route.Path)
		if route.Controller == nil {
			mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
			continue
		}

		mux = route.Controller.Boot(mux)
	}

	return middleware.NewMimeTypeEnforce(mux, "application/json", "multipart/form-data", "text/html")
}
