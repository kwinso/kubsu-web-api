package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/kwinso/kubsu-web-api/internal/server/dto"
)

func ParseBody[T dto.DTO](w http.ResponseWriter, r *http.Request, v T) (T, error) {
	if ExpectsJSON(r) {
		err := v.ParseJSON(w, r)
		return v, err
	} else {
		err := v.ParseFormData(w, r)
		return v, err
	}

	return v, errors.New("unsupported content type")
}

func ExpectsJSON(r *http.Request) bool {
	return IsContentType(r, "application/json")
}

func IsContentType(r *http.Request, contentType string) bool {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType == "" {
		return false
	}
	return strings.Contains(headerContentType, contentType)
}

func BadRequest(w http.ResponseWriter, r *http.Request, body string) {
	http.Error(w, body, http.StatusBadRequest)
}

func HttpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	if r.Header.Get("Accept") == "application/json" {
		resp := make(map[string]string)
		resp["error"] = err.Error()

		jsonResp, jsonErr := json.MarshalIndent(resp, "", "  ")
		if jsonErr != nil {
			log.Println("Error marshalling error response:", jsonErr)
			Error500(w, r)
			return
		}

		http.Error(w, string(jsonResp), http.StatusBadRequest)

		return
	}

	http.Error(w, err.Error(), code)
}

func Error500(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func WriteJSON(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
