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

func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
	err := cv.v.Struct(obj)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return fmt.Errorf(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return nil
}
