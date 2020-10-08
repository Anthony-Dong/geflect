package geflect

import (
	"reflect"

	"github.com/anthony-dong/geflect/common"
	"github.com/juju/errors"
)

type MethodsInfo struct {
	receiver     reflect.Value
	receiverType reflect.Type
	method       map[string]*Method
}

type Method struct {
	receiver     reflect.Value
	receiverType reflect.Type
	Method       reflect.Method
	NumIn        int
	ArgsType     []reflect.Type
}

type Reflect struct {
	value  interface{}
	_type  reflect.Type
	_value reflect.Value

	MethodsInfo
	FieldInfo
}

type FieldInfo struct {
	f     []Field
	index map[string]int // 记录索引
}

type Field struct {
	fieldType  reflect.StructField
	fieldValue reflect.Value
}

func (r *Reflect) GetValue() reflect.Value {
	return r._value
}

func (r *Reflect) GetType() reflect.Type {
	return r._type
}

func (r *Reflect) New() interface{} {
	return reflect.New(r._type).Elem().Interface()
}

// 获取class
func GetReflect(value interface{}) (*Reflect, error) {
	isNil := common.ViolationWithNotNil(value)
	if isNil {
		return nil, errors.New("the value is nil")
	}
	r := &Reflect{
		value:  value,
		_type:  reflect.TypeOf(value),
		_value: reflect.ValueOf(value),
	}
	if err := r.getMethodsInfo(); err != nil {
		return nil, err
	}
	if err := r.getFields(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Reflect) getFields() error {
	var (
		_value = reflect.Indirect(r._value)
	)
	if _value.Kind() != reflect.Struct {
		r.f = []Field{}
		return nil
		//return errors.Errorf("the %v is not struct type,the %v can not get field", _value.Type(), _value.Kind())
	}
	fieldNum := _value.NumField()
	r.f = make([]Field, 0, fieldNum)
	r.index = map[string]int{}
	for x := 0; x < fieldNum; x++ {
		filedInfo := Field{}
		filedInfo.fieldValue = _value.Field(x)
		filedInfo.fieldType = _value.Type().Field(x)
		r.f = append(r.f, filedInfo)
		r.index[filedInfo.fieldType.Name] = x
	}
	return nil
}

func (r *Reflect) getMethodsInfo() error {
	var (
		_type  = r._type
		_value = r._value
	)
	numMethod := _type.NumMethod()
	if numMethod == 0 {
		return errors.Errorf("%v not found method", r._type)
	}
	method := MethodsInfo{
		receiverType: _type,
		receiver:     _value,
		method:       map[string]*Method{},
	}
	for x := 0; x < numMethod; x++ {
		info := Method{
			receiver:     _value,
			receiverType: _type,
		}
		me := _type.Method(x)
		meType := me.Type
		in := meType.NumIn()
		info.NumIn = in - 1 // receiver 算一个参数
		info.Method = me
		info.ArgsType = []reflect.Type{}
		for x := 1; x < in; x++ {
			info.ArgsType = append(info.ArgsType, meType.In(x))
		}
		method.method[me.Name] = &info
	}
	r.MethodsInfo = method
	return nil
}
