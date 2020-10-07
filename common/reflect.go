package common

import "reflect"

var (
	ErrorType = reflect.TypeOf((*error)(nil)).Elem()
)

// true is nil
func ViolationWithNotNil(value interface{}) bool {
	if value == nil {
		return true
	}
	if reflect.ValueOf(value).IsNil() {
		return true
	}
	return false
}
