package main

import "fmt"

// 抽象工厂模式 - 提供一个创建一系列相关或相互依赖对象的接口，而无需指定它们具体的类

// 抽象产品 - 按钮
type Button interface {
	Render()
}

// 抽象产品 - 文本框
type TextBox interface {
	Render()
}

// 具体产品 - Windows按钮
type WindowsButton struct{}

func (w *WindowsButton) Render() {
	fmt.Println("Rendering Windows Button")
}

// 具体产品 - Windows文本框
type WindowsTextBox struct{}

func (w *WindowsTextBox) Render() {
	fmt.Println("Rendering Windows TextBox")
}

// 具体产品 - Mac按钮
type MacButton struct{}

func (m *MacButton) Render() {
	fmt.Println("Rendering Mac Button")
}

// 具体产品 - Mac文本框
type MacTextBox struct{}

func (m *MacTextBox) Render() {
	fmt.Println("Rendering Mac TextBox")
}

// 抽象工厂接口
type GUIFactory interface {
	CreateButton() Button
	CreateTextBox() TextBox
}

// 具体工厂 - Windows工厂
type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (w *WindowsFactory) CreateTextBox() TextBox {
	return &WindowsTextBox{}
}

// 具体工厂 - Mac工厂
type MacFactory struct{}

func (m *MacFactory) CreateButton() Button {
	return &MacButton{}
}

func (m *MacFactory) CreateTextBox() TextBox {
	return &MacTextBox{}
}

// 客户端代码
type Application struct {
	factory GUIFactory
	button  Button
	textBox TextBox
}

func NewApplication(factory GUIFactory) *Application {
	return &Application{
		factory: factory,
		button:  factory.CreateButton(),
		textBox: factory.CreateTextBox(),
	}
}

func (a *Application) Render() {
	a.button.Render()
	a.textBox.Render()
}

func main() {
	// 根据操作系统选择不同的工厂
	var factory GUIFactory

	osType := "windows" // 可以从环境变量或配置中获取

	switch osType {
	case "windows":
		factory = &WindowsFactory{}
	case "mac":
		factory = &MacFactory{}
	default:
		factory = &WindowsFactory{}
	}

	app := NewApplication(factory)
	app.Render()
}
