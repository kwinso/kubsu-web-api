package validation

import (
	"encoding/json"
	"net/http"
)

type Rule struct {
	Field      string
	Validators []Validator
	Message    string
}

type ValidationResult struct {
	Errors map[string]string `json:"errors"`
}

func RunRules(rules []Rule, value any) ValidationResult {
	var errs map[string]string
	for _, rule := range rules {
		for _, function := range rule.Validators {
			if !function() {
				if errs == nil {
					errs = make(map[string]string)
				}
				errs[rule.Field] = rule.Message
			}
		}
	}
	return ValidationResult{Errors: errs}
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (vr *ValidationResult) Format(r *http.Request) string {
	if r.Header.Get("Accept") == "application/json" {
		return vr.formatJSON()
	}
	return vr.formatHTML()
}

func (vr *ValidationResult) formatJSON() string {
	s, _ := json.Marshal(vr.Errors)
	return string(s)
}

func (vr *ValidationResult) formatHTML() string {
	// comma separated list of errors
	res := ""
	for k, v := range vr.Errors {
		res += k + ": " + v + "\n"
	}
	return res
}
