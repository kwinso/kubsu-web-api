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
	if IsContentType(r, "application/json") {
		err := v.ParseJSON(w, r)
		return v, err
	} else {
		err := v.ParseFormData(w, r)
		return v, err
	}
}

func ExpectsJSON(r *http.Request) bool {
	return IsAccept(r, "application/json")
}

func IsAccept(r *http.Request, typ string) bool {
	headerAccept := r.Header.Get("Accept")
	if headerAccept == "" {
		return false
	}

	return strings.Contains(headerAccept, typ)
}

func IsContentType(r *http.Request, contentType string) bool {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType == "" {
		return false
	}
	return strings.Contains(headerContentType, contentType)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	HttpError(w, r, errors.New("not found"), http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, r *http.Request, body string) {
	http.Error(w, body, http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	HttpError(w, r, errors.New("unauthorized"), http.StatusUnauthorized)
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
