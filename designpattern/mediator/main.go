package main

import "fmt"

// 中介者模式 - 定义一个中介对象来封装一系列对象之间的交互

// 中介者接口
type Mediator interface {
	SendMessage(message string, colleague Colleague)
	AddColleague(colleague Colleague)
}

// 同事接口
type Colleague interface {
	Send(message string)
	Receive(message string)
	SetMediator(mediator Mediator)
	GetName() string
}

// 具体中介者 - 聊天室
type ChatRoom struct {
	colleagues []Colleague
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		colleagues: make([]Colleague, 0),
	}
}

func (cr *ChatRoom) AddColleague(colleague Colleague) {
	cr.colleagues = append(cr.colleagues, colleague)
	colleague.SetMediator(cr)
	fmt.Printf("%s 加入了聊天室\n", colleague.GetName())
}

func (cr *ChatRoom) SendMessage(message string, sender Colleague) {
	fmt.Printf("[聊天室] %s 发送消息: %s\n", sender.GetName(), message)
	for _, colleague := range cr.colleagues {
		if colleague != sender {
			colleague.Receive(message)
		}
	}
}

// 具体同事 - 用户
type User struct {
	name     string
	mediator Mediator
}

func NewUser(name string) *User {
	return &User{name: name}
}

func (u *User) SetMediator(mediator Mediator) {
	u.mediator = mediator
}

func (u *User) Send(message string) {
	if u.mediator != nil {
		u.mediator.SendMessage(message, u)
	}
}

func (u *User) Receive(message string) {
	fmt.Printf("[%s] 收到消息: %s\n", u.name, message)
}

func (u *User) GetName() string {
	return u.name
}

// 机器人用户
type Bot struct {
	name     string
	mediator Mediator
}

func NewBot(name string) *Bot {
	return &Bot{name: name}
}

func (b *Bot) SetMediator(mediator Mediator) {
	b.mediator = mediator
}

func (b *Bot) Send(message string) {
	if b.mediator != nil {
		b.mediator.SendMessage(message, b)
	}
}

func (b *Bot) Receive(message string) {
	fmt.Printf("[机器人 %s] 收到消息: %s\n", b.name, message)
	// 机器人自动回复
	if b.mediator != nil {
		b.mediator.SendMessage("我是机器人，收到了您的消息", b)
	}
}

func (b *Bot) GetName() string {
	return b.name
}

// 智能家居中介者示例
type SmartHomeMediator interface {
	DeviceChanged(device SmartDevice, event string)
	RegisterDevice(device SmartDevice)
}

type SmartDevice interface {
	SetMediator(mediator SmartHomeMediator)
	GetName() string
	TriggerEvent(event string)
}

// 智能家居控制中心
type SmartHomeController struct {
	devices []SmartDevice
}

func NewSmartHomeController() *SmartHomeController {
	return &SmartHomeController{
		devices: make([]SmartDevice, 0),
	}
}

func (shc *SmartHomeController) RegisterDevice(device SmartDevice) {
	shc.devices = append(shc.devices, device)
	device.SetMediator(shc)
	fmt.Printf("设备 %s 已注册到智能家居系统\n", device.GetName())
}

func (shc *SmartHomeController) DeviceChanged(device SmartDevice, event string) {
	fmt.Printf("[智能家居] %s 触发事件: %s\n", device.GetName(), event)

	// 根据不同设备的事件，触发其他设备的联动
	switch device.GetName() {
	case "门锁":
		if event == "开锁" {
			shc.triggerDeviceAction("灯光", "开启")
			shc.triggerDeviceAction("空调", "启动")
		}
	case "温度传感器":
		if event == "温度过高" {
			shc.triggerDeviceAction("空调", "降温")
			shc.triggerDeviceAction("窗帘", "关闭")
		}
	case "光线传感器":
		if event == "光线变暗" {
			shc.triggerDeviceAction("灯光", "开启")
		}
	}
}

func (shc *SmartHomeController) triggerDeviceAction(deviceName, action string) {
	for _, device := range shc.devices {
		if device.GetName() == deviceName {
			fmt.Printf("  -> 联动触发: %s %s\n", deviceName, action)
			break
		}
	}
}

// 智能门锁
type SmartLock struct {
	name     string
	mediator SmartHomeMediator
}

func NewSmartLock() *SmartLock {
	return &SmartLock{name: "门锁"}
}

func (sl *SmartLock) SetMediator(mediator SmartHomeMediator) {
	sl.mediator = mediator
}

func (sl *SmartLock) GetName() string {
	return sl.name
}

func (sl *SmartLock) TriggerEvent(event string) {
	if sl.mediator != nil {
		sl.mediator.DeviceChanged(sl, event)
	}
}

func (sl *SmartLock) Unlock() {
	fmt.Printf("%s: 门锁已开启\n", sl.name)
	sl.TriggerEvent("开锁")
}

// 温度传感器
type TemperatureSensor struct {
	name     string
	mediator SmartHomeMediator
}

func NewTemperatureSensor() *TemperatureSensor {
	return &TemperatureSensor{name: "温度传感器"}
}

func (ts *TemperatureSensor) SetMediator(mediator SmartHomeMediator) {
	ts.mediator = mediator
}

func (ts *TemperatureSensor) GetName() string {
	return ts.name
}

func (ts *TemperatureSensor) TriggerEvent(event string) {
	if ts.mediator != nil {
		ts.mediator.DeviceChanged(ts, event)
	}
}

func (ts *TemperatureSensor) DetectHighTemperature() {
	fmt.Printf("%s: 检测到高温\n", ts.name)
	ts.TriggerEvent("温度过高")
}

// 智能灯光
type SmartLight struct {
	name     string
	mediator SmartHomeMediator
}

func NewSmartLight() *SmartLight {
	return &SmartLight{name: "灯光"}
}

func (sl *SmartLight) SetMediator(mediator SmartHomeMediator) {
	sl.mediator = mediator
}

func (sl *SmartLight) GetName() string {
	return sl.name
}

func (sl *SmartLight) TriggerEvent(event string) {
	if sl.mediator != nil {
		sl.mediator.DeviceChanged(sl, event)
	}
}

// 空调
type AirConditioner struct {
	name     string
	mediator SmartHomeMediator
}

func NewAirConditioner() *AirConditioner {
	return &AirConditioner{name: "空调"}
}

func (ac *AirConditioner) SetMediator(mediator SmartHomeMediator) {
	ac.mediator = mediator
}

func (ac *AirConditioner) GetName() string {
	return ac.name
}

func (ac *AirConditioner) TriggerEvent(event string) {
	if ac.mediator != nil {
		ac.mediator.DeviceChanged(ac, event)
	}
}

// 窗帘
type SmartCurtain struct {
	name     string
	mediator SmartHomeMediator
}

func NewSmartCurtain() *SmartCurtain {
	return &SmartCurtain{name: "窗帘"}
}

func (sc *SmartCurtain) SetMediator(mediator SmartHomeMediator) {
	sc.mediator = mediator
}

func (sc *SmartCurtain) GetName() string {
	return sc.name
}

func (sc *SmartCurtain) TriggerEvent(event string) {
	if sc.mediator != nil {
		sc.mediator.DeviceChanged(sc, event)
	}
}

func main() {
	fmt.Println("=== 聊天室中介者模式示例 ===")

	// 创建聊天室
	chatRoom := NewChatRoom()

	// 创建用户
	alice := NewUser("Alice")
	bob := NewUser("Bob")
	charlie := NewUser("Charlie")
	bot := NewBot("助手")

	// 用户加入聊天室
	chatRoom.AddColleague(alice)
	chatRoom.AddColleague(bob)
	chatRoom.AddColleague(charlie)
	chatRoom.AddColleague(bot)

	fmt.Println()

	// 用户发送消息
	alice.Send("大家好!")
	bob.Send("你好 Alice!")
	charlie.Send("机器人在吗?")

	fmt.Println("\n=== 智能家居中介者模式示例 ===")

	// 创建智能家居控制中心
	homeController := NewSmartHomeController()

	// 创建智能设备
	smartLock := NewSmartLock()
	tempSensor := NewTemperatureSensor()
	smartLight := NewSmartLight()
	airConditioner := NewAirConditioner()
	smartCurtain := NewSmartCurtain()

	// 注册设备到控制中心
	homeController.RegisterDevice(smartLock)
	homeController.RegisterDevice(tempSensor)
	homeController.RegisterDevice(smartLight)
	homeController.RegisterDevice(airConditioner)
	homeController.RegisterDevice(smartCurtain)

	fmt.Println()

	// 触发设备事件，观察联动效果
	fmt.Println("--- 场景1: 用户回家开门 ---")
	smartLock.Unlock()

	fmt.Println("\n--- 场景2: 温度传感器检测到高温 ---")
	tempSensor.DetectHighTemperature()
}
