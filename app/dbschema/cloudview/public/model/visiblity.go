//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type Visiblity string

const (
	Visiblity_Public  Visiblity = "PUBLIC"
	Visiblity_Private Visiblity = "PRIVATE"
)

func (e *Visiblity) Scan(value interface{}) error {
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
	case "PUBLIC":
		*e = Visiblity_Public
	case "PRIVATE":
		*e = Visiblity_Private
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for Visiblity enum")
	}

	return nil
}

func (e Visiblity) String() string {
	return string(e)
}
