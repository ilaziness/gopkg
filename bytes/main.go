package main

// bytes包提供了操作字节slice的方法，类似strings包的功能

import (
	"bytes"
	"log"
)

func main() {
	a, b := []byte("a"), []byte("b")
	// Compare 比较 a == b返回0，a < b 返回-1，a > b 返回1
	log.Println(bytes.Compare(a, b))

	// Contains 第一个参数是否包含第二个参数
	log.Println("Contains:")
	log.Println(bytes.Contains([]byte("abc"), []byte("bc"))) //true
	log.Println(bytes.Contains([]byte("abc"), []byte("d")))  //false
	log.Println(bytes.Contains([]byte(""), []byte("")))      //true

	// ContainsAny(a, b) a是否包含
	log.Println("ContainsAny:")
	log.Println(bytes.ContainsAny([]byte("abc"), "bc")) //true
	log.Println(bytes.ContainsAny([]byte("abc"), "我"))  //false
	log.Println(bytes.ContainsAny([]byte("我是谁"), "我是")) //true
	log.Println(bytes.ContainsAny([]byte(""), ""))      //false

}
