package handlers

import (
	"net/http"
	"strconv"

	"github.com/kwinso/kubsu-web-api/db"
	"github.com/kwinso/kubsu-web-api/httputil"
)

type AuthController struct {
}

func (c *AuthController) Boot(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("GET /login", http.HandlerFunc(c.LoginPage))
	mux.Handle("POST /login", http.HandlerFunc(c.Login))
	mux.Handle("GET /logout", http.HandlerFunc(c.Logout))

	return mux
}

func (c *AuthController) LoginPage(w http.ResponseWriter, r *http.Request) {
	httputil.WriteTemplate(w, r, "login", nil)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Delete cookies
	for _, cookie := range []string{"username", "password", "submission_id"} {
		http.SetCookie(w, &http.Cookie{
			Name:     cookie,
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if username == "" || password == "" {
		httputil.Unauthorized(w, r)
		return
	}

	submission, err := db.Query.GetSubmissionByCredentials(r.Context(), db.GetSubmissionByCredentialsParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		httputil.Unauthorized(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "password",
		Value:    password,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "submission_id",
		Value:    strconv.Itoa(int(submission.ID)),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
