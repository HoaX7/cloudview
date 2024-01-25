package permissions

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/utility"
	"fmt"
	"reflect"
	"strconv"
)

type PermissionMap struct {
	VISUALIZATION_AND_METRICS bool
	ALERTING                  bool
	SMART_DEBUGGING           bool
	MODIFY_RESOURCE_STATE     bool
	MANAGE_METRICS_PANEL      bool
}

type PermissionConstant struct {
	VISUALIZATION_AND_METRICS int
	ALERTING                  int
	SMART_DEBUGGING           int
	MODIFY_RESOURCE_STATE     int
	MANAGE_METRICS_PANEL      int
}

/*
Features - Users can choose the features they want to use
and pay for what they are using.
They can upgrade / downgrade features anytime based on their
requirements and will be billed accordingly.
*/
var (
	/*
		We are using bitwise system assigned to each feature.
		All the features used by the user are then stored in DB,
		which can later be deserialized to figure out all the
		features he has access to.
	*/
	permissionConstants = PermissionConstant{
		VISUALIZATION_AND_METRICS: 1 << 0,
		ALERTING:                  1 << 1,
		SMART_DEBUGGING:           1 << 2,
		MODIFY_RESOURCE_STATE:     1 << 3,
		MANAGE_METRICS_PANEL:      1 << 4,
	}
)

func SetPermissions(permissions []string) string {
	flag := 0
	for _, key := range permissions {
		field := reflect.ValueOf(&permissionConstants).Elem().FieldByName(key)
		if !field.IsValid() {
			continue
		}
		value := field.Int()
		flag = flag | int(value)
	}

	return serialize(flag)
}

func GetPermissions(hex string) (PermissionMap, error) {
	return deSerialize(hex)
}

// Serialize feature access to be stored in DB
func serialize(num int) string {
	return fmt.Sprintf("%x", num)
}

func deSerialize(hex string) (PermissionMap, error) {
	result := PermissionMap{}
	num, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		logger.Logger.Error("Failed to deserialize feature permissions:", err)
		return result, err
	}

	for _, key := range utility.GetKeys(&permissionConstants) {
		field := reflect.ValueOf(&permissionConstants).Elem().FieldByName(key)
		if !field.IsValid() {
			continue
		}
		value := field.Int()

		if num&value > 0 {
			reflect.ValueOf(&result).Elem().FieldByName(key).SetBool(true)
		}
	}

	return result, nil
}
