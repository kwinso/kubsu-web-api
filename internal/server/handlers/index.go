package handlers

import (
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/httputil"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Get username and password from cookies
	// username, _ := r.Cookie("username")
	// password, _ := r.Cookie("password")

	// msg := "Hello, world!"
	// if username != nil && password != nil {
	// 	msg = "Hello, " + username.Value + "! Your password is " + password.Value
	// }

	// w.Write([]byte(msg))
	httputil.WriteTemplate(w, r, "main", nil)
}
