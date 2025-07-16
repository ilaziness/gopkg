package main

import "fmt"

// 命令模式 - 将请求封装成对象，从而可以用不同的请求对客户进行参数化

// 命令接口
type Command interface {
	Execute()
	Undo()
}

// 接收者 - 电灯
type Light struct {
	isOn bool
}

func (l *Light) TurnOn() {
	l.isOn = true
	fmt.Println("电灯已打开")
}

func (l *Light) TurnOff() {
	l.isOn = false
	fmt.Println("电灯已关闭")
}

func (l *Light) IsOn() bool {
	return l.isOn
}

// 具体命令 - 打开电灯命令
type LightOnCommand struct {
	light *Light
}

func NewLightOnCommand(light *Light) *LightOnCommand {
	return &LightOnCommand{light: light}
}

func (c *LightOnCommand) Execute() {
	c.light.TurnOn()
}

func (c *LightOnCommand) Undo() {
	c.light.TurnOff()
}

// 具体命令 - 关闭电灯命令
type LightOffCommand struct {
	light *Light
}

func NewLightOffCommand(light *Light) *LightOffCommand {
	return &LightOffCommand{light: light}
}

func (c *LightOffCommand) Execute() {
	c.light.TurnOff()
}

func (c *LightOffCommand) Undo() {
	c.light.TurnOn()
}

// 空命令 - 用于初始化
type NoCommand struct{}

func (n *NoCommand) Execute() {}
func (n *NoCommand) Undo()    {}

// 调用者 - 遥控器
type RemoteControl struct {
	commands    []Command
	undoCommand Command
}

func NewRemoteControl() *RemoteControl {
	noCommand := &NoCommand{}
	return &RemoteControl{
		commands:    make([]Command, 7), // 7个插槽
		undoCommand: noCommand,
	}
}

func (r *RemoteControl) SetCommand(slot int, command Command) {
	if slot >= 0 && slot < len(r.commands) {
		r.commands[slot] = command
	}
}

func (r *RemoteControl) PressButton(slot int) {
	if slot >= 0 && slot < len(r.commands) && r.commands[slot] != nil {
		r.commands[slot].Execute()
		r.undoCommand = r.commands[slot]
	}
}

func (r *RemoteControl) PressUndoButton() {
	r.undoCommand.Undo()
}

// 宏命令 - 组合多个命令
type MacroCommand struct {
	commands []Command
}

func NewMacroCommand(commands []Command) *MacroCommand {
	return &MacroCommand{commands: commands}
}

func (m *MacroCommand) Execute() {
	for _, command := range m.commands {
		command.Execute()
	}
}

func (m *MacroCommand) Undo() {
	// 逆序撤销
	for i := len(m.commands) - 1; i >= 0; i-- {
		m.commands[i].Undo()
	}
}

func main() {
	// 创建接收者
	livingRoomLight := &Light{}
	kitchenLight := &Light{}

	// 创建命令
	livingRoomLightOn := NewLightOnCommand(livingRoomLight)
	livingRoomLightOff := NewLightOffCommand(livingRoomLight)
	kitchenLightOn := NewLightOnCommand(kitchenLight)
	kitchenLightOff := NewLightOffCommand(kitchenLight)

	// 创建遥控器
	remote := NewRemoteControl()

	// 设置命令
	remote.SetCommand(0, livingRoomLightOn)
	remote.SetCommand(1, livingRoomLightOff)
	remote.SetCommand(2, kitchenLightOn)
	remote.SetCommand(3, kitchenLightOff)

	// 测试单个命令
	fmt.Println("=== 测试单个命令 ===")
	remote.PressButton(0)    // 打开客厅灯
	remote.PressButton(2)    // 打开厨房灯
	remote.PressUndoButton() // 撤销上一个命令

	// 测试宏命令
	fmt.Println("\n=== 测试宏命令 ===")
	allLightsOn := NewMacroCommand([]Command{livingRoomLightOn, kitchenLightOn})
	allLightsOff := NewMacroCommand([]Command{livingRoomLightOff, kitchenLightOff})

	remote.SetCommand(4, allLightsOn)
	remote.SetCommand(5, allLightsOff)

	remote.PressButton(4) // 打开所有灯
	fmt.Println("撤销宏命令:")
	remote.PressUndoButton() // 撤销宏命令
}
