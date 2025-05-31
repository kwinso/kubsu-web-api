package dto

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	config "github.com/kwinso/kubsu-web-api/internal"
	"github.com/kwinso/kubsu-web-api/internal/db"
	"github.com/kwinso/kubsu-web-api/internal/server/dto/validation"
	"github.com/kwinso/kubsu-web-api/internal/util"
)

type CreateOrUpdateSubmissionDTO struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Bio       string `json:"bio"`
	Sex       int    `json:"sex"`
	Languages []int  `json:"languages"`
}

func (s *CreateOrUpdateSubmissionDTO) Validate() validation.ValidationResult {
	languages, _ := db.Query.GetAllLanguages(context.Background())
	languageIds := make([]int32, len(languages))
	for i, language := range languages {
		languageIds[i] = language.ID
	}

	rules := []validation.Rule{
		{
			Field: "name",
			Validators: []validation.Validator{
				validation.MinString(s.Name, 2),
				validation.MaxString(s.Name, 100),
			},
			Message: "Name must be between 2 and 100 characters",
		},
		{
			Field: "phone",
			Validators: []validation.Validator{
				validation.ValidPhone(s.Phone),
			},
			Message: "Invalid phone number",
		},
		{
			Field: "email",
			Validators: []validation.Validator{
				validation.ValidEmail(s.Email),
			},
			Message: "Invalid email address",
		},
		{
			Field: "birth_date",
			Validators: []validation.Validator{
				validation.ValidDate(s.BirthDate),
			},
			Message: "Invalid birth date",
		},
		{
			Field: "sex",
			Validators: []validation.Validator{
				validation.In(s.Sex, []int{0, 1}),
			},
			Message: "Sex must be 0 or 1",
		},
		{
			Field: "bio",
			Validators: []validation.Validator{
				validation.MinString(s.Bio, 10),
				validation.MaxString(s.Bio, 100),
			},
			Message: "Bio must be between 10 and 100 characters",
		},
		{
			Field: "languages",
			Validators: []validation.Validator{
				validation.MinLength(s.Languages, 1),
				validation.EachIn(s.Languages, util.Int32ToInt(languageIds)),
			},
			Message: "There should be at least one existing language selected",
		},
	}
	return validation.RunRules(rules, s)
}

func (s *CreateOrUpdateSubmissionDTO) ParseJSON(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&s)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (s *CreateOrUpdateSubmissionDTO) ParseFormData(w http.ResponseWriter, r *http.Request) error {
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

	languages, ok := r.PostForm["languages"]
	if ok {
		for _, language := range languages {
			languageID, err := strconv.Atoi(language)
			if err != nil {
				return errors.New("expected all languages valuse to be integers")
			}
			s.Languages = append(s.Languages, languageID)
		}
	}

	return nil
}
