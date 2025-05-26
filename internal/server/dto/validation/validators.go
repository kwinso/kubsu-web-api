package validation

import "regexp"

type Validator func() bool

func ValidPhone(phone string) Validator {
	return func() bool {
		re := regexp.MustCompile(`^\d{10}$`)
		return re.MatchString(phone)
	}
}

func ValidEmail(email string) Validator {
	return func() bool {
		// Simple email validation
		re := regexp.MustCompile(`@`)
		return re.MatchString(email)
	}
}
