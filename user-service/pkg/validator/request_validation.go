package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(data interface{}) error {
	var errorMessages []string

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", err.Field()))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid email address", err.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is not valid", err.Field()))
			}
		}
		return errors.New("Error validate :" + joinMessage(errorMessages))
	}

	return nil
}

func joinMessage(errorMessages []string) string {
	result := ""
	for i, msg := range errorMessages {
		if i > 0 {
			result += ", "
		} 
		result += msg
	}
	return result
}