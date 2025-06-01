package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/kwinso/kubsu-web-api/internal/db"
	"github.com/kwinso/kubsu-web-api/internal/server/dto"
	"github.com/kwinso/kubsu-web-api/internal/server/dto/validation"
	"github.com/kwinso/kubsu-web-api/internal/server/httputil"
	"github.com/kwinso/kubsu-web-api/internal/server/templates"
	"github.com/kwinso/kubsu-web-api/internal/util"
)

type SubmissionController struct {
}

func (c *SubmissionController) checkValidation(w http.ResponseWriter, r *http.Request, vr validation.ValidationResult, submission *dto.CreateOrUpdateSubmissionDTO, existingSubmissionId int32) bool {
	if !vr.HasErrors() {
		return false
	}
	if httputil.ExpectsJSON(r) {
		httputil.BadRequest(w, r, vr.Errors)
		return true
	}

	// Ugly, but works
	langs := make([]int32, len(submission.Languages))
	for i, language := range submission.Languages {
		langs[i] = int32(language)
	}

	allLanguages, err := db.Query.GetAllLanguages(r.Context())
	if err != nil {
		log.Println("Error getting languages:", err)
		httputil.Error500(w, r)
		return true
	}

	httputil.WriteTemplate(w, r, "main", templates.IndexContext{
		Submission: templates.IndexContextSubmission{
			ID:        existingSubmissionId,
			Name:      submission.Name,
			Phone:     submission.Phone,
			Email:     submission.Email,
			BirthDate: submission.BirthDate,
			Bio:       submission.Bio,
			Sex:       int(submission.Sex),
			Languages: langs,
		},
		Errors:    vr.Errors,
		Languages: allLanguages,
	})

	return true
}

// TODO: Add abiilty to show templates (for errors and for the main page)
func (c *SubmissionController) Boot(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /submissions", c.CreateSubmission)

	// General PUT handler for JSON clients
	mux.HandleFunc("PUT /submissions/{id}", c.UpdateSubmision)
	// Fallback for js-disabled clients who will submit via form data
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
		httputil.BadRequest(w, r, err.Error())
		return
	}

	vr := submission.Validate()
	failed := c.checkValidation(w, r, vr, submission, 0)
	if failed {
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
			Name:  "submission_id",
			Value: strconv.Itoa(int(createdSubmission.ID)),
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "password",
			Value: password,
		})

		http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
	}
}

func (c *SubmissionController) UpdateSubmision(w http.ResponseWriter, r *http.Request) {
	submission, err := httputil.ParseBody(w, r, &dto.CreateOrUpdateSubmissionDTO{})
	if err != nil {
		httputil.HttpError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	pathId := r.PathValue("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		httputil.BadRequest(w, r, "Invalid id provided")
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

	vr := submission.Validate()
	failed := c.checkValidation(w, r, vr, submission, submissionId)
	if failed {
		return
	}

	tx, err := db.DB.Begin(r.Context())
	if err != nil {
		log.Println("Error starting transaction:", err)
		httputil.Error500(w, r)
		return
	}
	defer tx.Rollback(r.Context())

	qtx := db.Query.WithTx(tx)

	err = qtx.UpdateSubmission(r.Context(), db.UpdateSubmissionParams{
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

	// Update submission languages
	err = qtx.DeleteSubmissionLanguages(r.Context(), submissionId)
	if err != nil {
		log.Println("Error deleting submission languages:", err)
		httputil.Error500(w, r)
		return
	}

	for _, language := range submission.Languages {
		err = qtx.AddLanguageToSubmission(r.Context(), db.AddLanguageToSubmissionParams{
			SubmissionID: submissionId,
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
		httputil.WriteJSON(w, r, submission)
	} else {
		http.Redirect(w, r, "/?success=true", http.StatusSeeOther)
	}
}
