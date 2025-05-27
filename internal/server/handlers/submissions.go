package handlers

import (
	"log"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/db"
	"github.com/kwinso/kubsu-web-api/internal/server/dto"
	"github.com/kwinso/kubsu-web-api/internal/server/httputil"
	"github.com/kwinso/kubsu-web-api/internal/util"
)

type SubmissionController struct {
}

func (c *SubmissionController) Boot(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /submissions", c.POST)
	return mux
}

func (c *SubmissionController) POST(w http.ResponseWriter, r *http.Request) {
	submission, err := httputil.ParseBody(w, r, &dto.CreateSubmissionDTO{})
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

		// Redirect back to main page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
