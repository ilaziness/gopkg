package main

import "fmt"

// 原型模式
// 解决对象复制问题，对象提供clone方法复制，使用者不需要知道实现细节就可以复制对象

type Prototype interface {
	Clone() Prototype
}

// 一般用来复制复杂对象，比如用创建者模式创建的复杂对象
type Message struct {
	Title  string
	Body   string
	Lenght int
}

func (m *Message) Clone() Prototype {
	msg := *m
	return &msg
}

func main() {
	msg := &Message{
		"title",
		"body",
		4,
	}

	repeatMsg := msg.Clone().(*Message)

	fmt.Printf("%+v\t%p\n", msg, msg)
	fmt.Printf("%+v\t%p\n", *repeatMsg, repeatMsg)
}
