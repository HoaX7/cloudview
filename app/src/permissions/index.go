package permissions

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/utility"
	"fmt"
	"reflect"
	"strconv"
)

type permissionMap struct {
	// project feature permissions
	VISUALIZATION_AND_METRICS bool
	ALERTING                  bool
	SMART_DEBUGGING           bool

	// Project Member permissions
	MEMBER_MODIFY_RESOURCE_STATE bool
	MEMBER_MANAGE_METRICS_PANEL  bool
	MEMBER_REPORT_METRICS        bool

	// User permissions
	USER_MODIFY_PROJECT          bool
	USER_MODIFY_PROVIDER_ACCOUNT bool
}

type permissionConstant struct {
	VISUALIZATION_AND_METRICS    int
	ALERTING                     int
	SMART_DEBUGGING              int
	MEMBER_MODIFY_RESOURCE_STATE int
	MEMBER_MANAGE_METRICS_PANEL  int
	MEMBER_REPORT_METRICS        int
	USER_MODIFY_PROJECT          int
	USER_MODIFY_PROVIDER_ACCOUNT int
}

/*
Features - Users can choose the features they want to use
and pay for what they are using.
They can upgrade / downgrade features anytime based on their
requirements and will be billed accordingly.
*/
var (
	VISUALIZATION_AND_METRICS    = "VISUALIZATION_AND_METRICS"
	ALERTING                     = "ALERTING"
	SMART_DEBUGGING              = "SMART_DEBUGGING"
	MEMBER_MODIFY_RESOURCE_STATE = "MEMBER_MODIFY_RESOURCE_STATE"
	MEMBER_MANAGE_METRICS_PANEL  = "MEMBER_MANAGE_METRICS_PANEL"
	MEMBER_REPORT_METRICS        = "MEMBER_REPORT_METRICS"
	USER_MODIFY_PROJECT          = "USER_MODIFY_PROJECT"
	USER_MODIFY_PROVIDER_ACCOUNT = "USER_MODIFY_PROVIDER_ACCOUNT"
	/*
		We are using bitwise system assigned to each feature.
		All the features used by the user are then stored in DB,
		which can later be deserialized to figure out all the
		features he has access to.
	*/
	permissionConstants = permissionConstant{
		VISUALIZATION_AND_METRICS:    1 << 0,
		ALERTING:                     1 << 1,
		SMART_DEBUGGING:              1 << 2,
		MEMBER_MODIFY_RESOURCE_STATE: 1 << 3,
		MEMBER_MANAGE_METRICS_PANEL:  1 << 4,
		MEMBER_REPORT_METRICS:        1 << 5,
		USER_MODIFY_PROJECT:          1 << 6,
		USER_MODIFY_PROVIDER_ACCOUNT: 1 << 7,
	}

	AllProjectPermissions = []string{
		VISUALIZATION_AND_METRICS,
		ALERTING,
		SMART_DEBUGGING,
	}

	AllProjectMemberPermissions = []string{
		MEMBER_MODIFY_RESOURCE_STATE,
		MEMBER_MANAGE_METRICS_PANEL,
		MEMBER_REPORT_METRICS,
	}

	AllUserPermissions = []string{
		USER_MODIFY_PROJECT,
		USER_MODIFY_PROVIDER_ACCOUNT,
	}
)

/*
Verifies if user has all permissions in the list.
*/
func VerifyPermissions(perms []string, hex string) bool {
	if hex == "" {
		logger.Logger.Log("No permission set, hex value is 0.")
		return false
	}
	permissions, err := GetPermissions(hex)
	flag := false
	if err != nil {
		logger.Logger.Error("VerifyPermissions: Failed with unexpected error", err)
		return flag
	}
	for _, key := range perms {
		field := reflect.ValueOf(&permissions).Elem().FieldByName(key)
		if !field.IsValid() {
			continue
		}
		flag = field.Bool()
		if !flag {
			break
		}
	}
	logger.Logger.Log("VerifyPermissions: verification completed for perms:", perms,
		"returning flag:", flag)
	return flag
}

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

func GetPermissions(hex string) (permissionMap, error) {
	return deSerialize(hex)
}

// Serialize feature access to be stored in DB
func serialize(num int) string {
	return fmt.Sprintf("%x", num)
}

func deSerialize(hex string) (permissionMap, error) {
	result := permissionMap{}
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
