package server

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/middleware"
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
// 2. When getting main page, it should check cookies for username and password, and render the old submission if found
// 2.1. If not set, should be a POST request to /submissions
// 2.2. If set, should be a POST request to /submissions/{id}
// 3. Submitting via JS. Check for cookies and change the request accordingly

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
