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

	err := v.RegisterValidation("notifyType", cv.notifyTypeValidate)
	if err != nil {
		panic(err)
	}

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
	case "notifyType":
		return fmt.Errorf("field %s is no notification type", field)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}

func (cv *CustomValidator) notifyTypeValidate(fl validator.FieldLevel) bool {
	if fl.Field().String() != "<notifyv1.NotifyType Value>" {
		return false
	}

	return true
}
