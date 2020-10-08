package geflect

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
	Name *string
	Age  int64
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

func TestMethod(t *testing.T) {
	class, err := GetReflect(new(M))
	if err != nil {
		t.Fatal(err)
	}
	method, err := class.GetMethod("Call")
	if err != nil {
		t.Fatal(err)
	}
	value, err := method.Call("")
	if err != nil {
		t.Fatal(errors.ErrorStack(err))
	}
	fmt.Println(value.IsNil())     // 判断是不是nil
	fmt.Println(value.GetError())  // 获取返回值是否有error，如果空||不存在就返回nil，不支持多个error
	fmt.Println(value.GetValue())  // 获取返回值
	fmt.Println(value.GetIndex(0)) // 获取返回值
}

func TestFields(t *testing.T) {
	str := "1111"
	m := &M{
		Name: &str,
		Age:  1,
	}
	class, err := GetReflect(m)
	if err != nil {
		t.Fatal(errors.ErrorStack(err))
	}
	{
		field, err := class.FieldByName("Name")
		if err != nil {
			t.Fatal(errors.ErrorStack(err))
		}
		new, _ := field.GetValue().(*string)
		*new = "111111"
		fmt.Println(field.GetValue().(*string))
	}

	{
		field, err := class.FieldByName("Age")
		if err != nil {
			t.Fatal(errors.ErrorStack(err))
		}
		fmt.Println(field.SetValue(1000))
		fmt.Println(field.IsPublic())
	}
	fmt.Println(*m.Name)
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

func TestClone(t *testing.T) {
}
