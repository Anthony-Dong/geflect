package geflect

import (
	"reflect"

	"github.com/anthony-dong/geflect/common"
	"github.com/juju/errors"
)

// get all field and that is safe
func (f *FieldInfo) Fields() []Field {
	fields := make([]Field, 0, len(f.f))
	copy(fields, f.f)
	return fields
}

// like java api
func (f *FieldInfo) FieldByName(fieldName string) (Field, error) {
	index, isExist := f.index[fieldName]
	if !isExist {
		return Field{}, errors.Errorf("not found %v field", fieldName)
	}
	return f.f[index], nil
}

// like java api
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

// like java api
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
