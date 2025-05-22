package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ValidatorInstance = validator.New()

func init() {
	ValidatorInstance.RegisterValidation("refreshTokenFormat", IsRefreshTokenFormat)
}

func ValidateStruct(s interface{}) error {
	err := ValidatorInstance.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				return fmt.Errorf("field '%s' validation failed on the '%s' tag", fieldError.Field(), fieldError.Tag())
			}
		}
		return err
	}
	return nil
}

func IsRefreshTokenFormat(f1 validator.FieldLevel) bool {
	token := f1.Field().String()

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	JTI := parts[1]

	_, err := uuid.Parse(JTI)
	if err != nil {
		return false
	}

	return true
}
