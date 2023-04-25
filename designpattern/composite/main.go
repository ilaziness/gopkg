package main

import "fmt"

// 组合模式

// 直接组合，就是将一个对象作为另一个对象的成员属性
type Message struct {
	Header *Header
	Body   *Body
}

type Header struct {
	Addr string
}
type Body struct {
	Content string
}

// 嵌入组合，就是利用了语言的匿名成员特性，本质上跟直接组合是一致的
type Message2 struct {
	Header
	Body
}

func main() {
	msg := Message{
		Header: &Header{Addr: "127.0.0.1"},
		Body:   &Body{Content: "content text"},
	}
	fmt.Println(msg.Header, msg.Body.Content)

	msg2 := Message2{
		Header{Addr: "192.168.0.1"},
		Body{Content: "text"},
	}
	fmt.Println(msg2, msg2.Header.Addr, msg2.Content)
}
