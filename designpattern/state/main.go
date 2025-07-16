package main

import "fmt"

// 状态模式 - 允许对象在内部状态改变时改变它的行为，对象看起来好像修改了它的类

// 状态接口
type State interface {
	InsertCoin(machine *VendingMachine)
	SelectProduct(machine *VendingMachine)
	DispenseProduct(machine *VendingMachine)
	RefundCoin(machine *VendingMachine)
	GetStateName() string
}

// 上下文 - 自动售货机
type VendingMachine struct {
	state        State
	coinInserted bool
	productCount int
}

func NewVendingMachine(productCount int) *VendingMachine {
	vm := &VendingMachine{
		coinInserted: false,
		productCount: productCount,
	}

	if productCount > 0 {
		vm.state = &NoCoinState{}
	} else {
		vm.state = &SoldOutState{}
	}

	return vm
}

func (vm *VendingMachine) SetState(state State) {
	fmt.Printf("状态切换: %s -> %s\n", vm.state.GetStateName(), state.GetStateName())
	vm.state = state
}

func (vm *VendingMachine) InsertCoin() {
	vm.state.InsertCoin(vm)
}

func (vm *VendingMachine) SelectProduct() {
	vm.state.SelectProduct(vm)
}

func (vm *VendingMachine) DispenseProduct() {
	vm.state.DispenseProduct(vm)
}

func (vm *VendingMachine) RefundCoin() {
	vm.state.RefundCoin(vm)
}

func (vm *VendingMachine) GetCurrentState() string {
	return vm.state.GetStateName()
}

func (vm *VendingMachine) ReleaseCoin() {
	if vm.coinInserted {
		fmt.Println("硬币已退回")
		vm.coinInserted = false
	}
}

func (vm *VendingMachine) ReleaseProduct() {
	if vm.productCount > 0 {
		fmt.Println("商品已出货")
		vm.productCount--
	}
}

func (vm *VendingMachine) GetProductCount() int {
	return vm.productCount
}

func (vm *VendingMachine) HasCoin() bool {
	return vm.coinInserted
}

func (vm *VendingMachine) SetCoinInserted(inserted bool) {
	vm.coinInserted = inserted
}

// 具体状态 - 无硬币状态
type NoCoinState struct{}

func (ncs *NoCoinState) InsertCoin(machine *VendingMachine) {
	fmt.Println("硬币已投入")
	machine.SetCoinInserted(true)
	machine.SetState(&HasCoinState{})
}

func (ncs *NoCoinState) SelectProduct(machine *VendingMachine) {
	fmt.Println("请先投入硬币")
}

func (ncs *NoCoinState) DispenseProduct(machine *VendingMachine) {
	fmt.Println("请先投入硬币")
}

func (ncs *NoCoinState) RefundCoin(machine *VendingMachine) {
	fmt.Println("没有硬币可退回")
}

func (ncs *NoCoinState) GetStateName() string {
	return "无硬币状态"
}

// 具体状态 - 有硬币状态
type HasCoinState struct{}

func (hcs *HasCoinState) InsertCoin(machine *VendingMachine) {
	fmt.Println("硬币已经投入，请选择商品")
}

func (hcs *HasCoinState) SelectProduct(machine *VendingMachine) {
	fmt.Println("商品已选择，正在出货...")
	machine.SetState(&DispensingState{})
}

func (hcs *HasCoinState) DispenseProduct(machine *VendingMachine) {
	fmt.Println("请先选择商品")
}

func (hcs *HasCoinState) RefundCoin(machine *VendingMachine) {
	machine.ReleaseCoin()
	machine.SetState(&NoCoinState{})
}

func (hcs *HasCoinState) GetStateName() string {
	return "有硬币状态"
}

// 具体状态 - 出货状态
type DispensingState struct{}

func (ds *DispensingState) InsertCoin(machine *VendingMachine) {
	fmt.Println("正在出货，请稍等")
}

func (ds *DispensingState) SelectProduct(machine *VendingMachine) {
	fmt.Println("正在出货，请稍等")
}

func (ds *DispensingState) DispenseProduct(machine *VendingMachine) {
	machine.ReleaseProduct()
	machine.SetCoinInserted(false)

	if machine.GetProductCount() > 0 {
		machine.SetState(&NoCoinState{})
	} else {
		fmt.Println("商品已售完")
		machine.SetState(&SoldOutState{})
	}
}

func (ds *DispensingState) RefundCoin(machine *VendingMachine) {
	fmt.Println("正在出货，无法退币")
}

func (ds *DispensingState) GetStateName() string {
	return "出货状态"
}

// 具体状态 - 售完状态
type SoldOutState struct{}

func (sos *SoldOutState) InsertCoin(machine *VendingMachine) {
	fmt.Println("商品已售完，硬币已退回")
}

func (sos *SoldOutState) SelectProduct(machine *VendingMachine) {
	fmt.Println("商品已售完")
}

func (sos *SoldOutState) DispenseProduct(machine *VendingMachine) {
	fmt.Println("商品已售完")
}

func (sos *SoldOutState) RefundCoin(machine *VendingMachine) {
	if machine.HasCoin() {
		machine.ReleaseCoin()
	} else {
		fmt.Println("没有硬币可退回")
	}
}

func (sos *SoldOutState) GetStateName() string {
	return "售完状态"
}

// 注意：交通灯示例在这里只是概念演示，实际应用中应该为交通灯单独设计状态接口
// 这里为了简化，复用了售货机的状态接口，但在实际项目中不建议这样做

func main() {
	fmt.Println("=== 自动售货机状态模式示例 ===")

	// 创建有3个商品的售货机
	machine := NewVendingMachine(3)

	fmt.Printf("初始状态: %s, 商品数量: %d\n\n", machine.GetCurrentState(), machine.GetProductCount())

	// 测试正常购买流程
	fmt.Println("--- 正常购买流程 ---")
	machine.InsertCoin()
	machine.SelectProduct()
	machine.DispenseProduct()
	fmt.Printf("当前状态: %s, 剩余商品: %d\n\n", machine.GetCurrentState(), machine.GetProductCount())

	// 测试退币功能
	fmt.Println("--- 测试退币功能 ---")
	machine.InsertCoin()
	machine.RefundCoin()
	fmt.Printf("当前状态: %s\n\n", machine.GetCurrentState())

	// 测试无硬币时的操作
	fmt.Println("--- 无硬币时的操作 ---")
	machine.SelectProduct()
	machine.DispenseProduct()
	machine.RefundCoin()
	fmt.Println()

	// 购买剩余商品直到售完
	fmt.Println("--- 购买剩余商品 ---")
	for machine.GetProductCount() > 0 {
		fmt.Printf("剩余商品: %d\n", machine.GetProductCount())
		machine.InsertCoin()
		machine.SelectProduct()
		machine.DispenseProduct()
		fmt.Printf("购买后状态: %s\n", machine.GetCurrentState())
	}

	// 测试售完状态
	fmt.Println("\n--- 测试售完状态 ---")
	machine.InsertCoin()
	machine.SelectProduct()
	machine.DispenseProduct()
}
