package server

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/handlers"
)

type Route struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

var Routes = []Route{
	{
		Method:  "GET",
		Path:    "/",
		Handler: handlers.Index,
	},
}

func RegisterRoutes(mux *http.ServeMux) {

	for _, route := range Routes {
		log.Println("Registering route:", route.Method, route.Path)
		mux.HandleFunc(route.Path, route.Handler)
	}

}
