package geflect

import (
	"reflect"

	"github.com/anthony-dong/geflect/common"
	"github.com/juju/errors"
)

// get method by name
func (m *MethodsInfo) GetMethod(methodName string) (*Method, error) {
	method, isExist := m.method[methodName]
	if !isExist {
		return nil, errors.Errorf("%s method not found", methodName)
	}
	return method, nil
}

// invoke method is core !!!
func (m *Method) invoke(receiver interface{}, needCheck bool, value ...interface{}) (*ReturnValue, error) {
	var (
		receiverValue reflect.Value
	)
	if needCheck {
		if common.ViolationWithNotNil(receiver) {
			return nil, errors.Errorf("the receiver is nil")
		}
		receiverValue = reflect.ValueOf(receiver)
		if receiverValue.Type() != m.receiverType {
			return nil, errors.Errorf("the receiver type is %v, but found %v type", m.receiverType, reflect.TypeOf(receiver))
		}
	} else {
		receiverValue = m.receiver
	}
	// 校验参数是否合法
	if err := m.ViolationArgs(value...); err != nil {
		return nil, errors.Trace(err)
	}
	// 返回值
	result := &ReturnValue{}
	// 参数
	args := make([]reflect.Value, 0, len(value)+1)
	args = append(args, receiverValue)
	for _, elem := range value {
		args = append(args, reflect.ValueOf(elem))
	}
	returnValues := m.Method.Func.Call(args)
	result.value = make([]interface{}, 0, len(returnValues))
	if returnValues == nil || len(returnValues) == 0 {
		return result, nil
	}
	for _, elem := range returnValues {
		result.value = append(result.value, elem.Interface())
	}
	return result, nil
}

// call use default receiver which use easy and safe if you receiver is empty
func (m *Method) Call(value ...interface{}) (*ReturnValue, error) {
	return m.invoke(nil, false, value...)
}

// invoke like java api !!
func (m *Method) Invoke(methodName string, receiver interface{}, value ...interface{}) (*ReturnValue, error) {
	return m.invoke(receiver, true, value...)
}

// in order to not panic !!!
func (m *Method) ViolationArgs(args ...interface{}) error {
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

// that is for go is multiple return value
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
