package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validation(u interface{}) []string {
	errs := []string{}
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		validationerrors := err.(validator.ValidationErrors)
		for _, err := range validationerrors {
			switch err.Tag() {
			case "required":
				errs = append(errs, fmt.Errorf("%s is required", err.Field()).Error())
			case "alpha":
				errs = append(errs, fmt.Errorf("%s is alphabet only", err.Field()).Error())
			case "email":
				errs = append(errs, fmt.Errorf("%s is not valid email format", err.Field()).Error())
			case "numeric":
				errs = append(errs, fmt.Errorf("%s is numeric only", err.Field()).Error())
			case "gte":
				errs = append(errs, fmt.Errorf("%s value must be greater than %s", err.Field(), err.Param()).Error())
			case "lte":
				errs = append(errs, fmt.Errorf("%s value must lower than %s", err.Field(), err.Param()).Error())
			}
		}
	}
	return errs
}
