package main

import "fmt"

// 代理模式: 为其他对象提供一种代理（Proxy）以控制对这个对象的访问
// 远程代理：rpc
// 虚拟代理：在需要大量资源的对象创建时，可以使用代理来延迟对象的实际创建，直到真正需要使用它。
// 安全代理：在访问敏感信息时，可以使用代理来控制对信息的访问权限
// 缓存代理: 缓存代理主要在Client与本体对象之间加上一层缓存，用于加速本体对象的访问，常见于连接数据库的场景。
// ...

// 定义接口
type Subject interface {
	Do() string
}

// 被代理类
type RealSubject struct{}

func (r *RealSubject) Do() string {
	return "RealSubject doing something"
}

// 代理类
type Proxy struct {
	realSubject *RealSubject
}

func (p *Proxy) Do() string {
	if p.realSubject == nil {
		p.realSubject = &RealSubject{}
	}
	result := "Proxy call RealSubject: "
	result += p.realSubject.Do()
	return result
}

// 通过Proxy来代理RealSubject对象执行功能
func main() {
	proxy := &Proxy{}
	fmt.Println(proxy.Do())
}
