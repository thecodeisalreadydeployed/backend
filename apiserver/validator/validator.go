package validator

import (
	govalidator "github.com/go-playground/validator/v10"
)

var validate *govalidator.Validate

type ValidationError struct {
	Field string
	Tag   string
	Value string
}

func Init() {
	validate = govalidator.New()
}

func CheckStruct(s interface{}) []ValidationError {
	var errors []ValidationError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(govalidator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field: err.StructNamespace(),
				Tag:   err.Tag(),
				Value: err.Param(),
			})
		}
	}
	return errors
}
