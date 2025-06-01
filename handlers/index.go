package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/kwinso/kubsu-web-api/db"
	"github.com/kwinso/kubsu-web-api/httputil"
	"github.com/kwinso/kubsu-web-api/templates"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Get username and password from cookies
	submissionId, _ := r.Cookie("submission_id")
	username, _ := r.Cookie("username")
	password, _ := r.Cookie("password")

	var submission db.Submission
	submissionLanguages := make([]int32, 0)
	if submissionId != nil && username != nil && password != nil {
		submissionId, err := strconv.Atoi(submissionId.Value)
		if err != nil {
			httputil.BadRequest(w, r, "Invalid auth cookie")
			return
		}
		submission, err = db.Query.GetFullSubmissionByIdAndCredentials(r.Context(), db.GetFullSubmissionByIdAndCredentialsParams{
			ID:       int32(submissionId),
			Username: username.Value,
			Password: password.Value,
		})

		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				httputil.Error500(w, r)
				return
			}
		}

		submissionLanguages, err = db.Query.GetSubmissionLanguages(r.Context(), submission.ID)
		if err != nil {
			log.Println("Error getting submission languages:", err)
			httputil.Error500(w, r)
			return
		}
	}

	allLanguages, err := db.Query.GetAllLanguages(r.Context())
	if err != nil {
		log.Println("Error getting languages:", err)
		httputil.Error500(w, r)
		return
	}

	success := r.URL.Query().Get("success")

	// w.Write([]byte(msg))
	httputil.WriteTemplate(w, r, "main", templates.IndexContext{
		Success: success == "true",
		Submission: templates.IndexContextSubmission{
			ID:        submission.ID,
			Name:      submission.Name,
			Phone:     submission.Phone,
			Email:     submission.Email,
			BirthDate: submission.BirthDate,
			Bio:       submission.Bio,
			Sex:       int(submission.Sex),
			Languages: submissionLanguages,
		},
		Errors:    map[string]string{},
		Languages: allLanguages,
	})
}
