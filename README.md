# geflect

## 1、method 

### Feature

- 1、支持快速的call方法，不需要反射，不需要panic，安全使用！！！
- 2、支持go的各种函数调用！！

### Quick start

```go
type M struct {
	name *string
	Age  int64
}

func (M) Call(str string) error {
	if str != "" {
		fmt.Println(str)
	}
	return errors.New("the str is nil")
}
```

下一步就可以使用了 ！！！

```go
func TestMethod(t *testing.T) {
	class, err := GetReflect(new(M))
	if err != nil {
		t.Fatal(err)
	}
	method, err := class.GetMethod("Call") // 获取method
	if err != nil {
		t.Fatal(err)
	}
	value, err := method.Call("") // 调用
	if err != nil {
		t.Fatal(errors.ErrorStack(err))
	}
	fmt.Println(value.IsNil())     // 判断是不是nil
	fmt.Println(value.GetError())  // 获取返回值是否有error，如果空||不存在就返回nil，不支持多个error
	fmt.Println(value.GetValue())  // 获取返回值
	fmt.Println(value.GetIndex(0)) // 获取返回值
}
```



## 2、Field

### Feature

- 1、支持get/set方法
- 2、支持查看是否public
- 3、支持查看tag信息

### Quick start

```go
func TestFields(t *testing.T) {
	str := "1111"
	m := &M{
		Name: &str,
		Age:  1,
	}
	// 获取reflect
	class, err := GetReflect(m)
	if err != nil {
		t.Fatal(errors.ErrorStack(err))
	}
	{
	// 字段
		field, err := class.FieldByName("Name")
		if err != nil {
			t.Fatal(errors.ErrorStack(err))
		}
		// 设置为空
		fmt.Println(field.SetValue(nil))
		fmt.Println(field.IsPublic())
	}

	{
		field, err := class.FieldByName("Age")
		if err != nil {
			t.Fatal(errors.ErrorStack(err))
		}
		fmt.Println(field.SetValue(1000))
		fmt.Println(field.IsPublic())
	}
	fmt.Println(m)
}
```

