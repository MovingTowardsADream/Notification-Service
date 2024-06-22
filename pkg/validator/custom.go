package custom_validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	v         *validator.Validate
	passwdErr error
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return cv.newValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return nil
}

func (cv *CustomValidator) newValidationError(field string, value interface{}, tag string, param string) error {
	switch tag {
	case "required":
		return fmt.Errorf("field %s is required", field)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}
