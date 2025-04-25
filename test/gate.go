package test

import (
	"flag"
	"fmt"
	"log"
	"order_go/internal/exchange/gateio"
	"order_go/internal/exchange/types"
	"order_go/internal/utils/config"
	"strconv"
	"time"
)

// 测试参数
var (
	// 测试模式
	TestMode string
	Symbol   string
	OrderID  string
	Currency string
	AccountType string
	
	// 订单参数
	OrderSide   = types.OrderSideBuy // 买入
	OrderAmount = 0.001              // 订单数量
	OrderPrice  = 0.0                // 订单价格，0表示使用市场价的90%
	
	// 精度参数
	PricePrecision  int // 价格精度（小数位数）
	AmountPrecision = 0 // 数量精度（小数位数），默认为0表示整数
)

// FormatPrice 格式化价格显示
func FormatPrice(price float64) string {
	format := fmt.Sprintf("%%.%df", PricePrecision)
	return fmt.Sprintf(format, price)
}

// FormatAmount 格式化数量显示
func FormatAmount(amount float64) string {
	// 如果精度为0，直接转为整数
	if AmountPrecision == 0 {
		return fmt.Sprintf("%d", int(amount))
	}
	
	format := fmt.Sprintf("%%.%df", AmountPrecision)
	return fmt.Sprintf(format, amount)
}

// ParseFlags 解析命令行参数
func ParseFlags() bool {
	// 定义命令行参数
	flag.StringVar(&TestMode, "test", "", "测试模式: price, order, cancel, status, balance")
	flag.StringVar(&Symbol, "symbol", "HYPE_USDT", "交易对")
	flag.StringVar(&OrderID, "orderid", "", "订单ID (用于cancel和status模式)")
	flag.StringVar(&Currency, "currency", "USDT", "币种 (用于balance模式，不指定则返回USDT)")
	flag.StringVar(&AccountType, "account", "spot", "账户类型: spot (用于balance模式)")
	flag.IntVar(&PricePrecision, "price-precision", 5, "价格精度（小数位数）")
	flag.IntVar(&AmountPrecision, "amount-precision", 0, "数量精度（小数位数），默认为0表示整数")
	
	// 订单参数
	var side string
	var amount float64
	var price float64
	flag.StringVar(&side, "side", "buy", "订单方向: buy, sell")
	flag.Float64Var(&amount, "amount", 0.0, "订单数量 (默认: 0.001)")
	flag.Float64Var(&price, "price", 0.0, "订单价格 (默认: 市场价的90%)")
	
	// 解析命令行参数
	flag.Parse()
	
	// 如果指定了测试模式，设置订单参数
	if TestMode != "" {
		SetOrderParams(side, amount, price)
		return true
	}
	
	return false
}

// getGateClient 创建并返回一个Gate.io客户端实例
func getGateClient() *gateio.Client {
	gateCfg, exists := config.GetExchangeConfig("gateio")
	if !exists {
		log.Fatal("找不到Gate.io配置")
	}
	return gateio.NewClient(gateCfg)
}

// TestGetPrice 测试获取价格
func TestGetPrice(symbol string) {
	fmt.Printf("正在获取 %s 的价格...\n", symbol)
	client := getGateClient()
	
	price, err := client.GetSymbolPrice(symbol)
	if err != nil {
		log.Fatalf("获取价格失败: %v", err)
	}
	
	fmt.Printf("%s 当前价格: %s USDT\n", symbol, FormatPrice(price))
}

// TestCreateOrder 测试创建订单
func TestCreateOrder(symbol string) {
	client := getGateClient()
	
	// 确定订单价格
	orderPrice := OrderPrice
	if orderPrice == 0.0 {
		// 只有在未指定价格时才查询市场价格
		fmt.Printf("未指定价格，正在获取 %s 的市场价格...\n", symbol)
		price, err := client.GetSymbolPrice(symbol)
		if err != nil {
			log.Fatalf("获取价格失败: %v", err)
		}
		
		fmt.Printf("%s 当前价格: %s USDT\n", symbol, FormatPrice(price))
		orderPrice = price * 0.9
		fmt.Printf("使用市场价的90%%: %s USDT\n", FormatPrice(orderPrice))
	} else {
		fmt.Printf("使用指定价格: %s USDT\n", FormatPrice(orderPrice))
	}
	
	// 创建订单
	order := &types.Order{
		Symbol:   symbol,
		Side:     OrderSide,
		Amount:   OrderAmount,
		Price:    orderPrice,
		ClientID: "t-" + strconv.FormatInt(time.Now().UnixNano(), 10), // 使用t-前缀加时间戳作为客户端订单ID
	}
	
	fmt.Printf("正在创建订单: %s %s %s @ %s USDT\n", 
		symbol, order.Side, FormatAmount(order.Amount), FormatPrice(order.Price))
	
	resp, err := client.CreateOrder(order)
	if err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	
	fmt.Printf("订单创建成功! 订单ID: %s, 状态: %s\n", resp.OrderID, resp.Status)
	fmt.Println("请记下订单ID，用于测试取消订单: go run main.go -test cancel -orderid <orderID>")
}

// TestCancelOrder 测试取消订单
func TestCancelOrder(symbol, orderID string) {
	fmt.Printf("正在取消订单 %s...\n", orderID)
	client := getGateClient()
	
	err := client.CancelOrder(symbol, orderID)
	if err != nil {
		log.Fatalf("取消订单失败: %v", err)
	}
	
	fmt.Printf("订单 %s 已成功取消\n", orderID)
}

// TestGetOrderStatus 测试获取订单状态
func TestGetOrderStatus(symbol, orderID string) {
	fmt.Printf("正在获取订单 %s 的状态...\n", orderID)
	client := getGateClient()
	
	resp, err := client.GetOrderStatus(symbol, orderID)
	if err != nil {
		log.Fatalf("获取订单状态失败: %v", err)
	}
	
	fmt.Printf("订单ID: %s\n", resp.OrderID)
	fmt.Printf("状态: %s\n", resp.Status)
	fmt.Printf("已成交数量: %s\n", FormatAmount(resp.FilledQty))
}

// TestGetBalance 测试获取账户余额
func TestGetBalance(accountType, currency string) {
	if currency == "" {
		fmt.Printf("正在获取%s账户所有币种余额...\n", getAccountTypeName(accountType))
	} else {
		fmt.Printf("正在获取%s账户 %s 的余额...\n", getAccountTypeName(accountType), currency)
	}
	
	client := getGateClient()
	
	// 使用新的GetAccountBalanceByType方法
	balances, err := client.GetAccountBalanceByType(accountType, currency)
	if err != nil {
		log.Fatalf("获取余额失败: %v", err)
	}
	
	if len(balances) == 0 {
		fmt.Println("账户没有余额")
		return
	}
	
	fmt.Printf("%s账户余额:\n", getAccountTypeName(accountType))
	fmt.Println("----------------------------------------")
	
	// 根据账户类型显示不同的表头
	if accountType == "futures" {
		fmt.Printf("%-10s %-15s %-15s %-15s\n", "币种", "可用", "未实现盈亏", "总计")
	} else {
		fmt.Printf("%-10s %-15s %-15s %-15s\n", "币种", "可用", "锁定", "总计")
	}
	
	fmt.Println("----------------------------------------")
	
	// 如果指定了币种，只显示该币种
	if currency != "" {
		if accountType == "futures" {
			available := balances[currency+".available_for_trading"]
			unrealizedPnl := balances[currency+".unrealized_pnl"]
			total := balances[currency+".total"]
			fmt.Printf("%-10s %-15.8f %-15.8f %-15.8f\n", currency, available, unrealizedPnl, total)
		} else {
			available := balances[currency+".available"]
			locked := balances[currency+".locked"]
			total := balances[currency+".total"]
			fmt.Printf("%-10s %-15.8f %-15.8f %-15.8f\n", currency, available, locked, total)
		}
		return
	}
	
	// 显示所有币种
	currencies := make(map[string]bool)
	for key := range balances {
		parts := key
		for _, suffix := range []string{".available", ".locked", ".total", ".unrealized_pnl", ".available_for_trading"} {
			if len(key) > len(suffix) && key[len(key)-len(suffix):] == suffix {
				parts = key[:len(key)-len(suffix)]
				break
			}
		}
		currencies[parts] = true
	}
	
	for currency := range currencies {
		if accountType == "futures" {
			available := balances[currency+".available_for_trading"]
			unrealizedPnl := balances[currency+".unrealized_pnl"]
			total := balances[currency+".total"]
			fmt.Printf("%-10s %-15.8f %-15.8f %-15.8f\n", currency, available, unrealizedPnl, total)
		} else {
			available := balances[currency+".available"]
			locked := balances[currency+".locked"]
			total := balances[currency+".total"]
			fmt.Printf("%-10s %-15.8f %-15.8f %-15.8f\n", currency, available, locked, total)
		}
	}
}

// getAccountTypeName 获取账户类型名称
func getAccountTypeName(accountType string) string {
	switch accountType {
	case "spot", "":
		return "现货"
	case "margin":
		return "保证金"
	case "futures":
		return "期货"
	default:
		return accountType
	}
}

// SetOrderParams 设置订单参数
func SetOrderParams(side string, amount float64, price float64) {
	// 设置订单参数
	if side != "" {
		if side == "buy" {
			OrderSide = types.OrderSideBuy
		} else if side == "sell" {
			OrderSide = types.OrderSideSell
		} else {
			log.Fatalf("无效的订单方向: %s", side)
			return
		}
	}
	
	if amount > 0 {
		OrderAmount = amount
	}
	
	if price > 0 {
		OrderPrice = price
	}
}

// RunTest 运行Gate.io测试
func RunTest() {
	// 根据测试模式执行不同的测试
	switch TestMode {
	case "price":
		// 测试获取价格
		TestGetPrice(Symbol)
	case "order":
		// 测试创建订单
		TestCreateOrder(Symbol)
	case "cancel":
		// 测试取消订单
		if OrderID == "" {
			log.Fatal("请提供订单ID: go run main.go -test cancel -orderid <orderID>")
		}
		TestCancelOrder(Symbol, OrderID)
	case "status":
		// 测试获取订单状态
		if OrderID == "" {
			log.Fatal("请提供订单ID: go run main.go -test status -orderid <orderID>")
		}
		TestGetOrderStatus(Symbol, OrderID)
	case "balance":
		// 测试获取账户余额
		TestGetBalance(AccountType, Currency)
	default:
		log.Fatalf("未知的测试模式: %s", TestMode)
	}
}