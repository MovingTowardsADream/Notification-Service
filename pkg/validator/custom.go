package custom_validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
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

func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	err := cv.v.Struct(obj)
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
	case "email":
		return fmt.Errorf("field %s must be a valid email address", field)
	case "password":
		return cv.passwdErr
	case "min":
		return fmt.Errorf("field %s must be at least %s characters", field, param)
	case "max":
		return fmt.Errorf("field %s must be at most %s characters", field, param)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}

func (cv *CustomValidator) mailValidate(fl validator.FieldLevel) bool {

	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("field %s must be a string", fl.FieldName())
		return false
	}

	// get the value of the field
	fieldValue := fl.Field().String()

	_ = fieldValue

	return true
}

func (cv *CustomValidator) phoneValidate(fl validator.FieldLevel) bool {

	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("field %s must be a string", fl.FieldName())
		return false
	}

	// get the value of the field
	fieldValue := fl.Field().String()

	_ = fieldValue

	return true
}
