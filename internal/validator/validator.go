package validator

import (
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already
	// initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue() returns true if a value is in a list of specific permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// ValidateInt checks if a value is empty and if it's a valid integer.
// Returns the parsed int64, a boolean for validity, and an error message.
func ValidateInt(value string, bitSize int) (bool, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return false, "This field cannot be blank"
	}
	_, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return false, "Invalid number"
	}
	return true, ""
}

func ValidateFloat(value string, bitSize int) (bool, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return false, "This field cannot be blank"
	}
	_, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return false, "Invalid float"
	}
	return true, ""
}
