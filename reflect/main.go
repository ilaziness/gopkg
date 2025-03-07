package main

import (
	"fmt"
	"reflect"
)

// reflect.Type：表示 Go 语言的类型信息。
// reflect.Value：表示 Go 语言的值信息。
// reflect.Kind：表示类型的种类（如 int、string、struct 等）。

func main() {
	var x float64 = 3.14
	//---------------获取变量的类型和值
	fmt.Println("Type:", reflect.TypeOf(x))   // 输出: Type: float64
	fmt.Println("Value:", reflect.ValueOf(x)) // 输出: Value: 3.14

	//-------------获取值的种类
	v := reflect.ValueOf(x)
	fmt.Println("Kind:", v.Kind()) // 输出: Kind: float64

	//------------------ 修改变量的值
	v = reflect.ValueOf(&x).Elem() // 获取 x 的地址并解引用
	v.SetFloat(2.71)               // 修改值
	fmt.Println("New value:", x)   // 输出: New value: 2.71

	//-----------------调用函数
	// 定义一个函数
	fn := func(a, b int) int {
		return a + b
	}
	// 获取函数的值
	v = reflect.ValueOf(fn)
	// 准备参数
	args := []reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf(2),
	}
	// 调用函数
	result := v.Call(args)
	fmt.Println("Result:", result[0].Int()) // 输出: Result: 3

	//----------------获取结构体字段
	p := Person{Name: "Alice", Age: 30}
	v = reflect.ValueOf(p)
	// 获取字段
	fmt.Println("Field 0:", v.Field(0)) // 输出: Field 0: Alice
	fmt.Println("Field 1:", v.Field(1)) // 输出: Field 1: 30

	//---------------修改结构体字段
	p2 := &Person{Name: "Alice", Age: 30}
	v = reflect.ValueOf(p2).Elem()
	// 修改字段
	v.Field(0).SetString("Bob")
	v.Field(1).SetInt(25)
	fmt.Println("Updated Person:", *p2) // 输出: Updated Person: {Bob 25}

	//------------------ 遍历结构体字段
	p = Person{Name: "Alice", Age: 30}
	t := reflect.TypeOf(p)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("Field %d: %s (%s)\n", i, field.Name, field.Type)
	}

	//----------------- 操作切片
	s := []int{1, 2, 3}
	v = reflect.ValueOf(s)

	// 获取切片长度
	fmt.Println("Length:", v.Len()) // 输出: Length: 3

	// 获取切片元素
	for i := 0; i < v.Len(); i++ {
		fmt.Printf("Element %d: %d\n", i, v.Index(i).Int())
	}

	//--------------------- 操作map
	m := map[string]int{"a": 1, "b": 2}
	v = reflect.ValueOf(m)
	// 获取映射的键
	keys := v.MapKeys()
	for _, key := range keys {
		value := v.MapIndex(key)
		fmt.Printf("Key: %s, Value: %d\n", key.String(), value.Int())
	}

	// -------------------- 创建新变量
	var x2 int
	t = reflect.TypeOf(x2)
	v = reflect.New(t)                  // 创建一个新的 int 变量
	fmt.Println("New value:", v.Elem()) // 输出: New value: 0

	//--------------------- 动态调用函数
	dyfun() // 输出: Result: 3

	// ------------------- 调用结构体方法
	// 调用 Add 方法的结果: 8
	//调用 Multiply 方法的结果: 15
	callStructMethod()

}

func dyfun() {
	var fn func(int, int) int64
	v := reflect.ValueOf(&fn).Elem()

	// 动态创建函数
	newFn := reflect.MakeFunc(v.Type(), func(args []reflect.Value) []reflect.Value {
		a := args[0].Int()
		b := args[1].Int()
		return []reflect.Value{reflect.ValueOf(a + b)}
	})

	v.Set(newFn)

	// 调用动态创建的函数
	result := fn(1, 2)
	fmt.Println("Result:", result) // 输出: Result: 3
}

type Person struct {
	Name string
	Age  int
}

// 定义一个结构体
type Calculator struct{}

// 定义一个方法
func (c Calculator) Add(a, b int) int {
	return a + b
}

// 定义另一个方法
func (c Calculator) Multiply(a, b int) int {
	return a * b
}

func callStructMethod() {
	// 创建结构体实例
	calc := Calculator{}

	// 获取结构体的 reflect.Value
	v := reflect.ValueOf(calc)

	// 动态调用 Add 方法
	methodName := "Add"
	method := v.MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("方法 %s 不存在\n", methodName)
		return
	}

	// 准备参数
	args := []reflect.Value{
		reflect.ValueOf(3),
		reflect.ValueOf(5),
	}

	// 调用方法
	result := method.Call(args)
	fmt.Printf("调用 %s 方法的结果: %d\n", methodName, result[0].Int())

	// 动态调用 Multiply 方法
	methodName = "Multiply"
	method = v.MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("方法 %s 不存在\n", methodName)
		return
	}

	// 调用方法
	result = method.Call(args)
	fmt.Printf("调用 %s 方法的结果: %d\n", methodName, result[0].Int())
}
