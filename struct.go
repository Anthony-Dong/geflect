package geflect

import (
	"reflect"
	"sync"

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
	value      interface{}
	_type      reflect.Type
	_value     reflect.Value
	addFieldDo sync.Once
	MethodsInfo
	FieldInfo
}

type FieldInfo struct {
	f     []Field
	index map[string]int // 记录索引
}

func (f *FieldInfo) Fields() []Field {
	fields := make([]Field, 0, len(f.f))
	copy(fields, f.f)
	return fields
}

func (f *FieldInfo) FieldByName(fieldName string) (Field, error) {
	index, isExist := f.index[fieldName]
	if !isExist {
		return Field{}, errors.Errorf("not found %v field", fieldName)
	}
	return f.f[index], nil
}

func (f *Field) SetValue(fieldValue interface{}) bool {
	if f.fieldValue.CanSet() {
		if common.ViolationWithNotNil(fieldValue) {
			f.fieldValue.Set(reflect.New(f.fieldType.Type).Elem()) // set nil !!!
			return true
		}
		newValue := reflect.ValueOf(fieldValue)
		if newValue.Type() == f.fieldType.Type {
			f.fieldValue.Set(reflect.Value{})
			return true
		}
	}
	return false
}

func (f *Field) GetTag() reflect.StructTag {
	return f.fieldType.Tag
}

func (f *Field) GetName() string {
	return f.fieldType.Name
}

// it is unsafe if f.Type is ptr , public !!
func (f *Field) GetValue() interface{} {
	return f.fieldValue.Interface()
}

// true public
// false lower case file character
func (f *Field) IsPublic() bool {
	return f.fieldType.PkgPath == ""
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

func (m *MethodsInfo) GetMethod(methodName string) (*Method, error) {
	method, isExist := m.method[methodName]
	if !isExist {
		return nil, errors.Errorf("%s method not found", methodName)
	}
	return method, nil
}

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

// method
func (m *Method) Call(value ...interface{}) (*ReturnValue, error) {
	return m.invoke(nil, false, value...)
}

// invoke
func (m *Method) Invoke(methodName string, receiver interface{}, value ...interface{}) (*ReturnValue, error) {
	return m.invoke(receiver, true, value...)
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
