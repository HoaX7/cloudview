package helpers

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	// Invalid uuid, this uuid is set by default if no value is sent
	if u == "00000000-0000-0000-0000-000000000000" {
		return false
	}
	_, err := uuid.Parse(u)
	return err == nil
}

func CheckEmptyFields(data interface{}) error {
	value := reflect.ValueOf(data)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := value.Type().Field(i).Name
		fieldValue := field.Interface()

		// Check if the field is a string, is not optional, and is empty
		if field.Kind() == reflect.String && !isOptional(value.Type().Field(i)) && fieldValue.(string) == "" {
			return fmt.Errorf("Field %s cannot be empty", fieldName)
		}
	}
	return nil
}

func isOptional(field reflect.StructField) bool {
	optionalTag := field.Tag.Get("optional")
	return optionalTag == "true"
}
