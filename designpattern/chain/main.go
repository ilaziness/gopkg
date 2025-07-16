package main

import "fmt"

// 责任链模式 - 将请求沿着处理者链传递，直到某个处理者处理它为止

// 请求结构
type Request struct {
	Level   int
	Message string
}

// 处理者接口
type Handler interface {
	SetNext(Handler)
	Handle(*Request)
}

// 基础处理者
type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(handler Handler) {
	b.next = handler
}

func (b *BaseHandler) Handle(request *Request) {
	if b.next != nil {
		b.next.Handle(request)
	}
}

// 具体处理者 - 信息级别处理者
type InfoHandler struct {
	BaseHandler
}

func (i *InfoHandler) Handle(request *Request) {
	if request.Level <= 1 {
		fmt.Printf("InfoHandler: 处理信息级别请求 - %s\n", request.Message)
	} else {
		fmt.Println("InfoHandler: 无法处理，传递给下一个处理者")
		i.BaseHandler.Handle(request)
	}
}

// 具体处理者 - 警告级别处理者
type WarningHandler struct {
	BaseHandler
}

func (w *WarningHandler) Handle(request *Request) {
	if request.Level <= 2 {
		fmt.Printf("WarningHandler: 处理警告级别请求 - %s\n", request.Message)
	} else {
		fmt.Println("WarningHandler: 无法处理，传递给下一个处理者")
		w.BaseHandler.Handle(request)
	}
}

// 具体处理者 - 错误级别处理者
type ErrorHandler struct {
	BaseHandler
}

func (e *ErrorHandler) Handle(request *Request) {
	if request.Level <= 3 {
		fmt.Printf("ErrorHandler: 处理错误级别请求 - %s\n", request.Message)
	} else {
		fmt.Println("ErrorHandler: 无法处理，传递给下一个处理者")
		e.BaseHandler.Handle(request)
	}
}

func main() {
	// 创建处理者
	infoHandler := &InfoHandler{}
	warningHandler := &WarningHandler{}
	errorHandler := &ErrorHandler{}

	// 构建责任链
	infoHandler.SetNext(warningHandler)
	warningHandler.SetNext(errorHandler)

	// 测试不同级别的请求
	requests := []*Request{
		{Level: 1, Message: "这是一个信息"},
		{Level: 2, Message: "这是一个警告"},
		{Level: 3, Message: "这是一个错误"},
		{Level: 4, Message: "这是一个严重错误"},
	}

	for _, request := range requests {
		fmt.Printf("\n处理请求: Level=%d, Message=%s\n", request.Level, request.Message)
		infoHandler.Handle(request)
	}
}
