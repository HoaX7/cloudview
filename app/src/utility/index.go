package utility

import (
	"reflect"
)

func ConvertMapToStruct(m map[string]interface{}, s interface{}) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()

	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		if value, ok := m[field.Name]; ok {
			if field.Type.Kind() == reflect.Struct {
				// The field in the struct is a struct itself
				subStruct := reflect.New(field.Type).Interface()
				if err := ConvertMapToStruct(value.(map[string]interface{}), subStruct); err != nil {
					return err
				}
				stValue.Field(i).Set(reflect.ValueOf(subStruct).Elem())
			} else {
				// The field is not a struct, so set it directly
				stValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func IsEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	default:
		zeroValue := reflect.Zero(value.Type())
		return reflect.DeepEqual(value.Interface(), zeroValue.Interface())
	}
}

/*
Usage: GetKeys(&struct{})
*/
func GetKeys(data interface{}) []string {
	r := reflect.ValueOf(data).Elem()
	rt := r.Type()
	keys := []string{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		keys = append(keys, field.Name)
	}

	return keys
}
