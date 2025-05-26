package dto

import (
	"net/http"

	"github.com/kwinso/kubsu-web-api/internal/server/dto/validation"
)

type DTO interface {
	Validate() validation.ValidationResult

	ParseJSON(w http.ResponseWriter, r *http.Request) error

	ParseFormData(w http.ResponseWriter, r *http.Request) error
}
