package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/dto"
	"github.com/kwinso/kubsu-web-api/internal/server/httputil"
)

type SubmissionController struct {
}

func (c *SubmissionController) Boot(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("POST /submissions", c.CreateSubmission)
	return mux
}

func (c *SubmissionController) CreateSubmission(w http.ResponseWriter, r *http.Request) {
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

	// back to json
	s, err := json.Marshal(submission)
	if err != nil {
		panic(err) // Best effort.
	}

	w.Write(s)
}
