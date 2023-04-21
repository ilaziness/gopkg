package main

import (
	"fmt"

	"github.com/ilaziness/gopkg/designpattern/sg"
)

// 单例模式
// 单例结构体首字母小写，限定访问范围，再实现一个首字母大写的访问函数，相当于static方法的作用

func main() {
	msg := sg.MsgPoolInstance.GetMsg()
	fmt.Println(msg.Count)
	msg.Count = 2
	sg.MsgPoolInstance.AddMsg(msg)
	fmt.Println(sg.MsgPoolInstance.GetMsg().Count)
}
