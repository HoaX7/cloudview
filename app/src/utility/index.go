package utility

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

func ContainsString(array []string, target string) (bool, string) {
	// Join the array into a single string with a delimiter
	joined := strings.Join(array, ",")
	if target == "" {
		return false, joined
	}

	// Use strings.Contains to check if the target string is present
	return strings.Contains(joined, target), joined
}

/*
Ugly workaround to parse `double precision[][]` data from postgres column.
We need to manually parse it since `pq` lib does not support 2D arrays.
*/
func ParseDoublePrecision2DToFloat2D(str []uint8) *[][]float64 {
	data := string(str)
	dimensions := strings.Split(strings.Trim(data, "{}"), "},{")

	var twoDResult [][]float64
	for _, dime := range dimensions {
		values := strings.Split(dime, ",")
		n1, _ := strconv.ParseFloat(values[0], 64)
		n2, _ := strconv.ParseFloat(values[1], 64)
		nums := []float64{n1, n2}
		twoDResult = append(twoDResult, nums)
	}

	return &twoDResult
}

/*
ExtractQueryParams updates the provided struct fields with query parameters from the request
NOTE: the `query` tag in the `struct` is required to extract query fields.

TODO - Improve this function to be robust to support pointer types
and other types such as int32, uint32, float etc...
*/
func ExtractQueryParams(req *http.Request, d interface{}) error {
	dType := reflect.TypeOf(d)

	if err := shouldBeStruct(dType); err != nil {
		return err
	}

	// Data Holder Value
	dhVal := reflect.ValueOf(d)

	// Loop over all the fields present in struct (Title, Body, JSON)
	for i := 0; i < dType.Elem().NumField(); i++ {

		// Give me ith field. Elem() is used to dereference the pointer
		field := dType.Elem().Field(i)

		// Get the value from field tag i.e in case of Title it is "title"
		key := field.Tag.Get("query")

		// Get the type of field
		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := req.URL.Query().Get(key)

		//  Get reference of field value provided to input `d`
		result := dhVal.Elem().Field(i)

		// set the value to string field
		// for other kinds we need to use different functions like; SetInt, Set etc
		if kind == reflect.String {
			result.SetString(val)
		} else if kind == reflect.Int || kind == reflect.Int64 {
			if val == "" {
				val = "0"
			}
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return errors.New("Unable to parse int: " + err.Error())
			}
			result.SetInt(i)
		} else {
			return errors.New("Unsupported type:" + kind.String())
		}

	}
	return nil
}
func shouldBeStruct(d reflect.Type) error {
	td := d.Elem()
	if td.Kind() != reflect.Struct {
		errStr := fmt.Sprintf("destination should be %v, found %v", reflect.Struct, td.Kind())
		return errors.New(errStr)
	}
	return nil
}
