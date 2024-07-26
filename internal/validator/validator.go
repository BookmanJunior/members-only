package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErros map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErros) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErros == nil {
		v.FieldErros = make(map[string]string)
	}

	if _, exists := v.FieldErros[key]; !exists {
		v.FieldErros[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (v *Validator) MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func (v *Validator) MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) < n
}

func (v *Validator) AreFieldsEqual(value1, value2 string) bool {
	return value1 == value2
}
