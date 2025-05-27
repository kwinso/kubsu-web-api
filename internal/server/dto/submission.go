package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	config "github.com/kwinso/kubsu-web-api/internal"
	"github.com/kwinso/kubsu-web-api/internal/server/dto/validation"
)

type CreateSubmissionDTO struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Bio       string `json:"bio"`
	Sex       int    `json:"sex"`
	Languages []int  `json:"languages"`
}

func (s *CreateSubmissionDTO) Validate() validation.ValidationResult {
	rules := []validation.Rule{
		{
			Field: "phone",
			Validators: []validation.Validator{
				validation.ValidPhone(s.Phone),
			},
			Message: "Invalid phone number",
		},
	}
	return validation.RunRules(rules, s)
}

func (s *CreateSubmissionDTO) ParseJSON(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&s)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (s *CreateSubmissionDTO) ParseFormData(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(config.MAX_MEMORY)
	if err != nil {
		return err
	}

	s.Name = r.PostFormValue("name")
	s.Phone = r.PostFormValue("phone")
	s.Email = r.PostFormValue("email")
	s.BirthDate = r.PostFormValue("birth_date")
	s.Bio = r.PostFormValue("bio")

	sex, err := strconv.Atoi(r.FormValue("sex"))
	if err != nil {
		return errors.New("expected sex to be an integer")
	}
	s.Sex = sex

	languages := r.PostFormValue("languages")
	if languages == "" {
		return errors.New("expected languages to be a comma separated list of integers")
	}

	for language := range strings.SplitSeq(languages, ",") {
		languageID, err := strconv.Atoi(language)
		if err != nil {
			return errors.New("expected languages to be a comma separated list of integers")
		}
		s.Languages = append(s.Languages, languageID)
	}
	fmt.Println(s)

	return nil
}
