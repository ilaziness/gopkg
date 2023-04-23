package main

import "fmt"

// 建造者模式
// 用来构建复杂对象，比如ORM构建查询对象过程 db.from(...).where(....).order(...)....query()

type Message struct {
	Header *Header
	Body   *Body
}

type Header struct {
	Schema string
	Addr   string
	Port   uint
}

type Body struct {
	Items map[string]string
}

type MessageBuilder struct {
	Msg *Message
}

func Build() *MessageBuilder {
	return &MessageBuilder{
		Msg: &Message{
			Header: &Header{},
			Body:   &Body{},
		},
	}
}

func (b *MessageBuilder) WithSchma(s string) *MessageBuilder {
	b.Msg.Header.Schema = s
	return b
}

func (b *MessageBuilder) WithAddr(addr string) *MessageBuilder {
	b.Msg.Header.Addr = addr
	return b
}

func (b *MessageBuilder) WithPort(Port uint) *MessageBuilder {
	b.Msg.Header.Port = Port
	return b
}

func (b *MessageBuilder) WithBody(d map[string]string) *MessageBuilder {
	b.Msg.Body.Items = d
	return b
}

func main() {
	msg := Build().WithSchma("https").WithAddr("127.0.0.1").WithPort(9000).WithBody(map[string]string{"a": "123"})
	fmt.Printf("%+v\n", msg.Msg.Header)
	fmt.Printf("%+v\n", msg.Msg.Body)
}
