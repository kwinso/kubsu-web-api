package validation

import (
	"encoding/json"
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

func (vr *ValidationResult) AsJsonString() string {
	s, _ := json.Marshal(vr.Errors)
	return string(s)
}
