//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var MetricType = &struct {
	CPUUsage       postgres.StringExpression
	RAMUsage       postgres.StringExpression
	DiskOperations postgres.StringExpression
}{
	CPUUsage:       postgres.NewEnumValue("CPU_USAGE"),
	RAMUsage:       postgres.NewEnumValue("RAM_USAGE"),
	DiskOperations: postgres.NewEnumValue("DISK_OPERATIONS"),
}
