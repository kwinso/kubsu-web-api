package httputil

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kwinso/kubsu-web-api/internal/server/dto"
	"github.com/kwinso/kubsu-web-api/internal/server/templates"
)

type HttpErrorTemplateData struct {
	Message string
}

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
	HttpError(w, r, "not found", http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, r *http.Request, body string) {
	HttpError(w, r, body, http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	HttpError(w, r, "unauthorized", http.StatusUnauthorized)
}

func HttpError(w http.ResponseWriter, r *http.Request, body string, code int) {
	if r.Header.Get("Accept") == "application/json" {
		resp := make(map[string]string)
		resp["error"] = body

		jsonResp, jsonErr := json.MarshalIndent(resp, "", "  ")
		if jsonErr != nil {
			log.Println("Error marshalling error response:", jsonErr)
			Error500(w, r)
			return
		}

		http.Error(w, string(jsonResp), code)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	renderedTemplate, err := templates.Render(strconv.Itoa(code), HttpErrorTemplateData{body})
	if err != nil {
		log.Println("Error rendering template:", err)
		Error500(w, r)
		return
	}

	w.WriteHeader(code)
	w.Write(renderedTemplate)
}

func Error500(w http.ResponseWriter, r *http.Request) {
	HttpError(w, r, "Internal Server Error", http.StatusInternalServerError)
}

func WriteJSON(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func WriteTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := templates.Render(name, data)
	if err != nil {
		log.Println("Error rendering template:", err)
		Error500(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(tmpl)
}
