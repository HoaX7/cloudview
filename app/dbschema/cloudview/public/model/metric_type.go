//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type MetricType string

const (
	MetricType_CPUUsage       MetricType = "CPU_USAGE"
	MetricType_RAMUsage       MetricType = "RAM_USAGE"
	MetricType_DiskOperations MetricType = "DISK_OPERATIONS"
)

func (e *MetricType) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "CPU_USAGE":
		*e = MetricType_CPUUsage
	case "RAM_USAGE":
		*e = MetricType_RAMUsage
	case "DISK_OPERATIONS":
		*e = MetricType_DiskOperations
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for MetricType enum")
	}

	return nil
}

func (e MetricType) String() string {
	return string(e)
}
