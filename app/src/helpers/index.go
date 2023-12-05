package helpers

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
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
