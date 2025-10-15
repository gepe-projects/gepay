package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func GenerateMessage(err error, structT reflect.Type) map[string]any {
	var vErr validator.ValidationErrors
	if errors.As(err, &vErr) {
		messages := make(map[string]any, len(vErr))
		for _, v := range vErr {
			var fieldName string
			field, ok := structT.FieldByName(v.Field())
			if !ok {
				fieldName = v.Field()
			} else {
				fieldName = field.Tag.Get("json")
			}

			switch v.Tag() {
			case "email":
				messages[fieldName] = fmt.Sprintf("%s is not valid email", v.Value())
			case "required":
				messages[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "min":
				messages[fieldName] = fmt.Sprintf("%s must be at least %s characters", fieldName, v.Param())
			case "max":
				messages[fieldName] = fmt.Sprintf("%s must be at most %s characters", fieldName, v.Param())
			case "len":
				messages[fieldName] = fmt.Sprintf("%s must be exactly %s characters", fieldName, v.Param())
			case "uuid":
				messages[fieldName] = fmt.Sprintf("%s is not a valid UUID", fieldName)
			case "alpha":
				messages[fieldName] = fmt.Sprintf("%s must only contain alphabetic characters", fieldName)
			case "alpha_dash":
				messages[fieldName] = fmt.Sprintf("%s must only contain alphabetic characters, dashes, and underscores", fieldName)
			case "alpha_num":
				messages[fieldName] = fmt.Sprintf("%s must only contain alphabetic characters and numbers", fieldName)
			case "numeric":
				messages[fieldName] = fmt.Sprintf("%s must be a numeric value", fieldName)
			case "gt":
				messages[fieldName] = fmt.Sprintf("%s must be greater than %s", fieldName, v.Param())
			case "gte":
				messages[fieldName] = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, v.Param())
			case "lt":
				messages[fieldName] = fmt.Sprintf("%s must be less than %s", fieldName, v.Param())
			case "lte":
				messages[fieldName] = fmt.Sprintf("%s must be less than or equal to %s", fieldName, v.Param())
			case "url":
				messages[fieldName] = fmt.Sprintf("%s is not a valid URL", fieldName)
			case "hex":
				messages[fieldName] = fmt.Sprintf("%s must be a valid hexadecimal string", fieldName)
			case "date":
				messages[fieldName] = fmt.Sprintf("%s is not a valid date, e.g: 2006-01-02", fieldName)
			case "timezone":
				messages[fieldName] = fmt.Sprintf("%s is not a valid timezone e.g: UTC,+08:00,Asia,Jakarta,America,New_York", fieldName)
			case "ip":
				messages[fieldName] = fmt.Sprintf("%s is not a valid IP address", fieldName)
			}
		}
		return messages
	}
	return map[string]any{"error": err.Error()}
}
