package features

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/utility"
	"fmt"
	"reflect"
	"strconv"
)

type FeatureMap struct {
	VISUALIZATION_AND_METRICS bool
	ALERTING                  bool
	SMART_DEBUGGING           bool
}

type FeatureConstant struct {
	VISUALIZATION_AND_METRICS int
	ALERTING                  int
	SMART_DEBUGGING           int
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
	featureConstants = FeatureConstant{
		VISUALIZATION_AND_METRICS: 1 << 0,
		ALERTING:                  1 << 1,
		SMART_DEBUGGING:           1 << 2,
	}
)

func SetPermissions(permissions []string) string {
	flag := 0
	for _, key := range permissions {
		field := reflect.ValueOf(&featureConstants).Elem().FieldByName(key)
		if !field.IsValid() {
			continue
		}
		value := field.Int()
		flag = flag | int(value)
	}

	return serialize(flag)
}

func GetPermissions(hex string) (FeatureMap, error) {
	return deSerialize(hex)
}

// Serialize feature access to be stored in DB
func serialize(num int) string {
	return fmt.Sprintf("%x", num)
}

func deSerialize(hex string) (FeatureMap, error) {
	result := FeatureMap{}
	num, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		logger.Logger.Error("Failed to deserialize feature permissions:", err)
		return result, err
	}

	for _, key := range utility.GetKeys(&featureConstants) {
		field := reflect.ValueOf(&featureConstants).Elem().FieldByName(key)
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
