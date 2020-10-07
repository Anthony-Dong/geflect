package core

import (
	"fmt"
	"github.com/anthony-dong/geflect/common"
	"github.com/juju/errors"
	"reflect"
	"testing"
)

type MI interface {
}

type M struct {
}

func (M) Call(str string) error {
	if str != "" {
		fmt.Println(str)
	}
	return errors.New("the str is nil")
}

func TestIsNil(t *testing.T) {
	var x = (*int)(nil)
	fmt.Println(common.ViolationWithNotNil(x))
}

func TestGetMethod(t *testing.T) {
	method, _ := GetMethod(new(M))
	fmt.Println(method)
}

func TestCall(t *testing.T) {
	method, err := GetMethod(new(M))
	if err != nil {
		t.Fatal(err)
	}
	value, err := method.Call("Call", "111")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(value.IsNil())
	fmt.Println(value.GetError())
	fmt.Println(value.GetValue())
}

func BenchmarkCall(b *testing.B) {
	method, _ := GetMethod(new(M))
	for i := 0; i < b.N; i++ {
		_, err := method.Call("Call", "111")
		if err != nil {
			panic(err)
		}
	}
}

func TestImpl(t *testing.T) {
	fmt.Println(reflect.TypeOf((*error)(nil)).Elem())
	fmt.Println(reflect.TypeOf((error)(nil)))
	var x interface{}
	x = (*error)(nil)
	fmt.Println(x == nil)
	x = (error)(nil)
	fmt.Println(x == nil)
}

func BenchmarkImpl(b *testing.B) {
	err := errors.New("err msg")
	for i := 0; i < b.N; i++ {
		if !reflect.TypeOf(err).Implements(common.ErrorType) {
			b.Fatal(err)
		}
	}
}

func BenchmarkAssert(b *testing.B) {
	err := errors.New("err msg")
	for i := 0; i < b.N; i++ {
		if err, _ := err.(error); err != nil {
			b.Fatal(err)
		}
	}
}
