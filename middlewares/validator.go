package middlewares

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func translateValidationError(err *ValidationError) string {
	switch err.Tag {
	case "required":
		return fmt.Sprintf("%s is required", err.Field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field)
	case "uuid":
		return fmt.Sprintf("%s must be in uuid-v4 format", err.Field)
	case "unique":
		return fmt.Sprintf("%s must be unique", err.Field)
	case "gte":
		return fmt.Sprintf("%s must be more or equal than %s charachters", err.Field, err.Param)
	case "gt":
		return fmt.Sprintf("%s must be more than %s charachters", err.Field, err.Param)
	case "lte":
		return fmt.Sprintf("%s must be less or equal than %s charachters", err.Field, err.Param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s charachters", err.Field, err.Param)
	default:
		return "undefined validation error"
	}
}

func parseValidationErrors(errors []ValidationError) map[string]string {
	result := map[string]string{}
	for _, err := range errors {
		result[err.Field] = translateValidationError(&err)
	}
	return result
}

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	field, found := structType.FieldByName(fieldName)
	if !found {
		return fieldName // fallback to the struct field name
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return fieldName // fallback to the struct field name
	}
	return jsonTag
}

func ValidatePayload(payload interface{}) (map[string]string, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}
		validationErrors := err.(validator.ValidationErrors)
		result := map[string]string{}
		payloadType := reflect.TypeOf(payload)

		for _, err := range validationErrors {
			field := err.StructField()
			name := getJSONFieldName(payloadType, field)
			tag := err.Tag()
			param := err.Param()
			valError := ValidationError{
				Field: name,
				Tag:   tag,
				Param: param,
			}
			result[name] = translateValidationError(&valError)
		}
		return result, nil
	}
	return nil, nil

}
