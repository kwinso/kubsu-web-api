package server

import (
	"embed"
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/handlers"
)

//go:embed all:static
var staticFiles embed.FS

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

	// file, err := staticFiles.Open(".ss/main.min.css")
	// if err != nil {
	// 	return err
	// }

	bootstrapedMux := bootstrapRouting(mux, routes)
	mux.Handle("GET /static/{file...}", http.FileServer(http.FS(staticFiles)))
	// mux.Handle("GET /static/{file...}", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.URL.Path)
	// 	// file, _ := staticFs.Open(r.URL.Path[1:])
	// 	w.Write([]byte("Hello, world!"))
	// })))

	port := "8080"
	log.Println("Starting server on port", port)
	return http.ListenAndServe(":"+port, bootstrapedMux)

}
