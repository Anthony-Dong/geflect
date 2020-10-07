package method

import (
	"fmt"
	"reflect"
	"testing"
)

type U int

func (U) Call() string {
	return "call"
}

/**
指针可以调用非指针函数

非指针函数无法调用指针函数
*/
func TestName(t *testing.T) {
	_type := reflect.TypeOf(new(U)).Elem()
	method, isExist := _type.MethodByName("Call")
	if !isExist {
		fmt.Println("nil")
		return
	}
	fmt.Println(method.Func.Call())
}
