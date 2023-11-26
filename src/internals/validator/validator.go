package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

var UsernameRX = regexp.MustCompile("^[a-zA-Z]")

//var PasswordRX = regexp.MustCompile(`^\\(?=.*[A-Z])(\\?=.*[a-z])(?=.*\d)(\\?=.*[\W_])`)

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NoBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func NoEmptyTextSlice(value []string) bool {
	if value == nil {
		return false
	} else {
		return len(value) <= 0
	}
}

func IsPositiveNumber(value int) bool {
	rx := regexp.MustCompile(`^\d+$`)
	return rx.MatchString(string(value))
}

func IsStrongPassword(value string) bool {
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
