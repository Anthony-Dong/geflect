# geflect

## 1、method 

```go
type M struct {
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
func TestCall(t *testing.T) {
	method, err := GetMethod(new(M))
	if err != nil {
		t.Fatal(err)
	}
	value, err := method.Call("Call", "111")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(value.IsNil()) // 判断是不是nil
	fmt.Println(value.GetError()) // 获取返回值是否有error，如果空||不存在就返回nil，不支持多个error
	fmt.Println(value.GetValue()) // 获取返回值
}
```



