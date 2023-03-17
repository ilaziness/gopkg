package main

import (
	"fmt"
	"unsafe"
)

// unsafe.Pointer的操作规则
// 任何类型的指针都可以转化成 unsafe.Pointer；
// unsafe.Pointer 可以转化成任何类型的指针；
// uintptr 可以转换为 unsafe.Pointer；
// unsafeP.ointer 可以转换为 uintptr；

type Admin struct {
	Name     string
	Age      int
	Language string
}

func main() {
	i := 30
	iPtr1 := &i

	fmt.Println(i, iPtr1)
	iPtr2 := (*int64)(unsafe.Pointer(iPtr1))
	*iPtr2 = 8
	fmt.Println(i, iPtr2)

	admin := Admin{
		Name: "test1",
		Age:  20,
	}
	ptr := &admin
	name := (*string)(unsafe.Pointer(ptr))
	*name = "123"
	fmt.Println(*ptr, admin)

	// Offsetof() 获取结构体成员的偏移量，进而获取到成员地址
	age := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Offsetof(ptr.Age)))
	*age = 25
	fmt.Println(*ptr, admin)

	// Sizeof() 函数可以获取成员大小，进而计算出成员的地址
	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(int(0)) + unsafe.Sizeof(string(""))))
	//lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(0) + unsafe.Sizeof("")))
	*lang = "en"
	fmt.Println(*ptr, admin)

	lang = (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Offsetof(ptr.Language)))
	*lang = "cn"
	fmt.Println(*ptr, admin)

	// string byte零拷贝互相转换，go版本>=1.20
	s := "abc"
	//b := []byte("efgh")
	b := []byte{'e', 'f', 'g', 'h'}
	fmt.Println(s, StringToBytes(s))
	fmt.Println(b, BytesToString(b))
	fmt.Println(b, BytesToString2(b))
	// output:
	// abc [97 98 99]
	// [101 102 103 104] efgh
	// [101 102 103 104] efgh

	s = ""
	b = make([]byte, 0)
	fmt.Println(s, StringToBytes(s))
	fmt.Println(b, BytesToString(b))
	fmt.Println(b, BytesToString2(b))
	// output:
	//   []
	// []
	// []
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

func BytesToString2(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
