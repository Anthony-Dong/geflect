package core

import (
	"reflect"

	"github.com/anthony-dong/geflect/common"
	"github.com/juju/errors"
)

type Method struct {
	receiver     reflect.Value
	receiverType reflect.Type
	method       map[string]*MethodInfo
}

func GetMethod(value interface{}) (*Method, error) {
	isNil := common.ViolationWithNotNil(value)
	if isNil {
		return nil, errors.New("the value is nil")
	}
	_type := reflect.TypeOf(value)
	numMethod := _type.NumMethod()
	if numMethod == 0 {
		return nil, errors.Errorf("%v not found method", _type)
	}
	method := Method{
		receiver:     reflect.ValueOf(value),
		receiverType: _type,
		method:       map[string]*MethodInfo{},
	}

	for x := 0; x < numMethod; x++ {
		info := MethodInfo{}
		info.ArgsType = []reflect.Type{}
		me := _type.Method(x)
		meType := me.Type
		in := meType.NumIn()
		info.NumIn = in - 1 // receiver 算一个参数
		info.Method = me
		for x := 1; x < in; x++ {
			info.ArgsType = append(info.ArgsType, meType.In(x))
		}
		method.method[me.Name] = &info
	}
	return &method, nil
}

type MethodInfo struct {
	Method   reflect.Method
	NumIn    int
	ArgsType []reflect.Type
}

func (m MethodInfo) ViolationArgs(args ...interface{}) error {
	if (args == nil || len(args) == 0) && m.NumIn == 0 {
		return nil
	}
	if len(args) != len(m.ArgsType) {
		return errors.Errorf("method need %d args but find %d args", len(m.ArgsType), len(args))
	}
	for index, elem := range args {
		if m.ArgsType[index] != reflect.TypeOf(elem) {
			return errors.Errorf("method need %v type arg in %d index arg but find %v type arg", m.ArgsType[index], index, reflect.TypeOf(elem))
		}
	}
	return nil
}

func (m *Method) invoke(methodName string, receiver interface{}, needCheck bool, value ...interface{}) (*ReturnValue, error) {
	if needCheck {
		if common.ViolationWithNotNil(receiver) {
			return nil, errors.Errorf("the receiver is nil")
		}
		if reflect.TypeOf(receiver) != m.receiverType {
			return nil, errors.Errorf("the receiver type is %v, but found %v type", m.receiverType, reflect.TypeOf(receiver))
		}
	}
	// 判断方法是否存在
	method, isExist := m.method[methodName]
	if !isExist {
		return nil, errors.Errorf("%s method not found", methodName)
	}

	// 校验参数是否合法
	if err := method.ViolationArgs(value...); err != nil {
		return nil, errors.Trace(err)
	}

	// 返回值
	result := &ReturnValue{}
	// 参数
	args := make([]reflect.Value, 0, len(value)+1)
	args = append(args, reflect.ValueOf(receiver))
	for _, elem := range value {
		args = append(args, reflect.ValueOf(elem))
	}
	returnValues := method.Method.Func.Call(args)
	result.value = make([]interface{}, 0, len(returnValues))
	if returnValues == nil || len(returnValues) == 0 {
		return result, nil
	}
	for _, elem := range returnValues {
		result.value = append(result.value, elem.Interface())
	}
	return result, nil
}

func (m *Method) Call(methodName string, value ...interface{}) (*ReturnValue, error) {
	return m.invoke(methodName, nil, false, value...)
}

func (m *Method) Invoke(methodName string, receiver interface{}, value ...interface{}) (*ReturnValue, error) {
	return m.invoke(methodName, receiver, true, value...)
}

type ReturnValue struct {
	value []interface{}
}

func (r *ReturnValue) GetValue() []interface{} {
	return r.value
}

func (r *ReturnValue) IsNil() bool {
	return r.value == nil || len(r.value) == 0
}

func (r *ReturnValue) GetError() error {
	if r.IsNil() {
		return nil
	}
	for _, elem := range r.value {
		// 断言速度大约是反射判断imp的30倍起步，所以直接断言即可
		if err, _ := elem.(error); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReturnValue) GetIndex(index int) interface{} {
	if len(r.value) >= index {
		return nil
	}
	return r.value[index]
}
