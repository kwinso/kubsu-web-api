package validation

import (
	"regexp"
)

// TODO: Probably this thing should return message instead of a bool to allow different errors
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

func ValidDate(date string) Validator {
	return func() bool {
		// Simple date validation
		re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		return re.MatchString(date)
	}
}

func EachIn[T comparable](vals []T, allowed []T) Validator {
	return func() bool {
		for _, v := range vals {
			if !In(v, allowed)() {
				return false
			}
		}
		return true
	}
}

// Checks if the value is in the allowed list
func In[T comparable](val T, allowed []T) Validator {
	return func() bool {
		for _, a := range allowed {
			if a == val {
				return true
			}
		}
		return false
	}
}

func MinLength(vals []int, min int) Validator {
	return func() bool {
		return len(vals) >= min
	}
}

func MinString(s string, min int) Validator {
	return func() bool {
		return len(s) >= min
	}
}

func MaxString(s string, max int) Validator {
	return func() bool {
		return len(s) <= max
	}
}
