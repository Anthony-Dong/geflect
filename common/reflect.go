package common

import (
	"reflect"
)

var (
	ErrorType     = reflect.TypeOf((*error)(nil)).Elem()
	isNilCheckMap = map[reflect.Kind]interface{}{
		reflect.Chan:          nil,
		reflect.Func:          nil,
		reflect.Map:           nil,
		reflect.Ptr:           nil,
		reflect.UnsafePointer: nil,
		reflect.Interface:     nil,
		reflect.Slice:         nil,
	}
)

// Chan, Func, Map, Ptr, UnsafePointer:
// Interface, Slice

// true is nil
func ViolationWithNotNil(value interface{}) bool {
	if value == nil {
		return true
	}
	_value := reflect.ValueOf(value)
	_, isExist := isNilCheckMap[_value.Kind()]
	if !isExist {
		return false
	}
	if _value.IsNil() {
		return true
	}
	return false
}
