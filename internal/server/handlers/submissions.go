package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/kwinso/kubsu-web-api/internal/db"
	"github.com/kwinso/kubsu-web-api/internal/server/dto"
	"github.com/kwinso/kubsu-web-api/internal/server/httputil"
	"github.com/kwinso/kubsu-web-api/internal/util"
)

type SubmissionController struct {
}

// TODO: Add abiilty to show templates (for errors and for the main page)
func (c *SubmissionController) Boot(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /submissions", c.CreateSubmission)

	// General PUT handler for JSON clients
	mux.HandleFunc("PUT /submissions/{id}", c.UpdateSubmision)
	// Fallback for non-js clients who will submit via form data
	mux.HandleFunc("POST /submissions/{id}", c.UpdateSubmision)

	return mux
}

func (c *SubmissionController) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	cookieUsername, _ := r.Cookie("username")
	cookiePass, _ := r.Cookie("password")

	if cookieUsername != nil && cookiePass != nil {
		httputil.BadRequest(w, r, "you are already logged in")
		return
	}

	submission, err := httputil.ParseBody(w, r, &dto.CreateOrUpdateSubmissionDTO{})
	if err != nil {
		httputil.HttpError(w, r, err, http.StatusBadRequest)
		return
	}

	vr := submission.Validate()
	if vr.HasErrors() {
		httputil.BadRequest(w, r, vr.Format(r))
		return
	}

	username, password := util.GenerateRandomCredentials()

	tx, err := db.DB.Begin(r.Context())
	if err != nil {
		log.Println("Error starting transaction:", err)
		httputil.Error500(w, r)
		return
	}
	defer tx.Rollback(r.Context())

	qtx := db.Query.WithTx(tx)

	createdSubmission, err := qtx.CreateSubmission(r.Context(), db.CreateSubmissionParams{
		Name:      submission.Name,
		Phone:     submission.Phone,
		Email:     submission.Email,
		BirthDate: submission.BirthDate,
		Bio:       submission.Bio,
		Sex:       int16(submission.Sex),
		Username:  username,
		Password:  password,
	})

	if err != nil {
		log.Println("Error creating submission:", err)
		httputil.Error500(w, r)
		return
	}

	for _, language := range submission.Languages {
		err := qtx.AddLanguageToSubmission(r.Context(), db.AddLanguageToSubmissionParams{
			SubmissionID: createdSubmission.ID,
			LanguageID:   int32(language),
		})

		if err != nil {
			log.Println("Error adding language to submission:", err)
			httputil.Error500(w, r)
			return
		}
	}

	tx.Commit(r.Context())

	if httputil.ExpectsJSON(r) {
		httputil.WriteJSON(w, r, createdSubmission)
	} else {
		// Set username and password cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    username,
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "password",
			Value:    password,
			HttpOnly: true,
		})

		w.Write([]byte("Here will be your main page"))
	}
}

func (c *SubmissionController) UpdateSubmision(w http.ResponseWriter, r *http.Request) {
	submission, err := httputil.ParseBody(w, r, &dto.CreateOrUpdateSubmissionDTO{})
	if err != nil {
		httputil.HttpError(w, r, err, http.StatusBadRequest)
		return
	}

	vr := submission.Validate()
	if vr.HasErrors() {
		httputil.BadRequest(w, r, vr.Format(r))
		return
	}

	pathId := r.PathValue("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		httputil.BadRequest(w, r, "invalid id provided")
		return
	}

	username, err := r.Cookie("username")
	if err != nil {
		httputil.Unauthorized(w, r)
		return
	}
	password, err := r.Cookie("password")
	if err != nil {
		httputil.Unauthorized(w, r)
		return
	}

	submissionId, err := db.Query.GetSubmissionByIdAndCredentials(r.Context(), db.GetSubmissionByIdAndCredentialsParams{
		ID:       int32(id),
		Username: username.Value,
		Password: password.Value,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httputil.NotFound(w, r)
			return
		}

		log.Println("Error getting submission:", err)
		httputil.Error500(w, r)
		return
	}

	err = db.Query.UpdateSubmission(r.Context(), db.UpdateSubmissionParams{
		Name:      submission.Name,
		Phone:     submission.Phone,
		Email:     submission.Email,
		BirthDate: submission.BirthDate,
		Bio:       submission.Bio,
		Sex:       int16(submission.Sex),
		ID:        submissionId,
	})
	if err != nil {
		log.Println("Error updating submission:", err)
		httputil.Error500(w, r)
		return
	}

	if httputil.ExpectsJSON(r) {
		httputil.WriteJSON(w, r, submission)
	} else {
		// Redirect back to main page
		// TODO: Render a page
		w.Write([]byte("Here will be your main page"))
	}
}
