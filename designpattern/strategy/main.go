package main

import "fmt"

// 策略模式 - 定义一系列算法，把它们一个个封装起来，并且使它们可相互替换

// 策略接口
type PaymentStrategy interface {
	Pay(amount float64) bool
}

// 具体策略 - 信用卡支付
type CreditCardPayment struct {
	cardNumber string
	cvv        string
}

func NewCreditCardPayment(cardNumber, cvv string) *CreditCardPayment {
	return &CreditCardPayment{
		cardNumber: cardNumber,
		cvv:        cvv,
	}
}

func (cc *CreditCardPayment) Pay(amount float64) bool {
	fmt.Printf("使用信用卡支付 %.2f 元 (卡号: %s)\n", amount, cc.maskCardNumber())
	// 模拟支付处理
	if amount > 0 && len(cc.cardNumber) >= 16 {
		fmt.Println("信用卡支付成功!")
		return true
	}
	fmt.Println("信用卡支付失败!")
	return false
}

func (cc *CreditCardPayment) maskCardNumber() string {
	if len(cc.cardNumber) < 4 {
		return cc.cardNumber
	}
	return "****-****-****-" + cc.cardNumber[len(cc.cardNumber)-4:]
}

// 具体策略 - 支付宝支付
type AlipayPayment struct {
	account string
}

func NewAlipayPayment(account string) *AlipayPayment {
	return &AlipayPayment{account: account}
}

func (ap *AlipayPayment) Pay(amount float64) bool {
	fmt.Printf("使用支付宝支付 %.2f 元 (账号: %s)\n", amount, ap.account)
	// 模拟支付处理
	if amount > 0 && ap.account != "" {
		fmt.Println("支付宝支付成功!")
		return true
	}
	fmt.Println("支付宝支付失败!")
	return false
}

// 具体策略 - 微信支付
type WechatPayment struct {
	phoneNumber string
}

func NewWechatPayment(phoneNumber string) *WechatPayment {
	return &WechatPayment{phoneNumber: phoneNumber}
}

func (wp *WechatPayment) Pay(amount float64) bool {
	fmt.Printf("使用微信支付 %.2f 元 (手机号: %s)\n", amount, wp.maskPhoneNumber())
	// 模拟支付处理
	if amount > 0 && len(wp.phoneNumber) == 11 {
		fmt.Println("微信支付成功!")
		return true
	}
	fmt.Println("微信支付失败!")
	return false
}

func (wp *WechatPayment) maskPhoneNumber() string {
	if len(wp.phoneNumber) != 11 {
		return wp.phoneNumber
	}
	return wp.phoneNumber[:3] + "****" + wp.phoneNumber[7:]
}

// 上下文 - 购物车
type ShoppingCart struct {
	items           []string
	totalAmount     float64
	paymentStrategy PaymentStrategy
}

func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{
		items: make([]string, 0),
	}
}

func (sc *ShoppingCart) AddItem(item string, price float64) {
	sc.items = append(sc.items, item)
	sc.totalAmount += price
	fmt.Printf("添加商品: %s, 价格: %.2f 元\n", item, price)
}

func (sc *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	sc.paymentStrategy = strategy
}

func (sc *ShoppingCart) Checkout() bool {
	if sc.paymentStrategy == nil {
		fmt.Println("请选择支付方式!")
		return false
	}

	fmt.Printf("\n购物车商品: %v\n", sc.items)
	fmt.Printf("总金额: %.2f 元\n", sc.totalAmount)
	fmt.Println("开始支付...")

	return sc.paymentStrategy.Pay(sc.totalAmount)
}

// 折扣策略接口
type DiscountStrategy interface {
	ApplyDiscount(amount float64) float64
	GetDescription() string
}

// 无折扣策略
type NoDiscountStrategy struct{}

func (nds *NoDiscountStrategy) ApplyDiscount(amount float64) float64 {
	return amount
}

func (nds *NoDiscountStrategy) GetDescription() string {
	return "无折扣"
}

// 百分比折扣策略
type PercentageDiscountStrategy struct {
	percentage float64
}

func NewPercentageDiscountStrategy(percentage float64) *PercentageDiscountStrategy {
	return &PercentageDiscountStrategy{percentage: percentage}
}

func (pds *PercentageDiscountStrategy) ApplyDiscount(amount float64) float64 {
	return amount * (1 - pds.percentage/100)
}

func (pds *PercentageDiscountStrategy) GetDescription() string {
	return fmt.Sprintf("%.0f%% 折扣", pds.percentage)
}

// 固定金额折扣策略
type FixedAmountDiscountStrategy struct {
	discountAmount float64
}

func NewFixedAmountDiscountStrategy(discountAmount float64) *FixedAmountDiscountStrategy {
	return &FixedAmountDiscountStrategy{discountAmount: discountAmount}
}

func (fads *FixedAmountDiscountStrategy) ApplyDiscount(amount float64) float64 {
	result := amount - fads.discountAmount
	if result < 0 {
		return 0
	}
	return result
}

func (fads *FixedAmountDiscountStrategy) GetDescription() string {
	return fmt.Sprintf("减 %.2f 元", fads.discountAmount)
}

// 增强版购物车，支持折扣策略
type EnhancedShoppingCart struct {
	*ShoppingCart
	discountStrategy DiscountStrategy
}

func NewEnhancedShoppingCart() *EnhancedShoppingCart {
	return &EnhancedShoppingCart{
		ShoppingCart:     NewShoppingCart(),
		discountStrategy: &NoDiscountStrategy{},
	}
}

func (esc *EnhancedShoppingCart) SetDiscountStrategy(strategy DiscountStrategy) {
	esc.discountStrategy = strategy
}

func (esc *EnhancedShoppingCart) Checkout() bool {
	if esc.paymentStrategy == nil {
		fmt.Println("请选择支付方式!")
		return false
	}

	originalAmount := esc.totalAmount
	discountedAmount := esc.discountStrategy.ApplyDiscount(originalAmount)

	fmt.Printf("\n购物车商品: %v\n", esc.items)
	fmt.Printf("原价: %.2f 元\n", originalAmount)
	fmt.Printf("折扣策略: %s\n", esc.discountStrategy.GetDescription())
	fmt.Printf("实付金额: %.2f 元\n", discountedAmount)
	fmt.Println("开始支付...")

	return esc.paymentStrategy.Pay(discountedAmount)
}

func main() {
	// 基础购物车示例
	fmt.Println("=== 基础购物车示例 ===")
	cart := NewShoppingCart()
	cart.AddItem("笔记本电脑", 5999.00)
	cart.AddItem("鼠标", 199.00)

	// 尝试不同的支付策略
	fmt.Println("\n--- 使用信用卡支付 ---")
	cart.SetPaymentStrategy(NewCreditCardPayment("1234567890123456", "123"))
	cart.Checkout()

	fmt.Println("\n--- 使用支付宝支付 ---")
	cart.SetPaymentStrategy(NewAlipayPayment("user@example.com"))
	cart.Checkout()

	fmt.Println("\n--- 使用微信支付 ---")
	cart.SetPaymentStrategy(NewWechatPayment("13800138000"))
	cart.Checkout()

	// 增强版购物车示例（支持折扣）
	fmt.Println("\n\n=== 增强版购物车示例 ===")
	enhancedCart := NewEnhancedShoppingCart()
	enhancedCart.AddItem("手机", 3999.00)
	enhancedCart.AddItem("手机壳", 99.00)

	// 应用不同的折扣策略
	fmt.Println("\n--- 应用10%折扣 ---")
	enhancedCart.SetDiscountStrategy(NewPercentageDiscountStrategy(10))
	enhancedCart.SetPaymentStrategy(NewAlipayPayment("user@example.com"))
	enhancedCart.Checkout()

	fmt.Println("\n--- 应用满减200元 ---")
	enhancedCart.SetDiscountStrategy(NewFixedAmountDiscountStrategy(200))
	enhancedCart.Checkout()
}
