package gateio

import (
	"context"
	"fmt"
	"order_go/internal/exchange/types"
	"order_go/internal/utils/config"
	"strconv"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
)

type Client struct {
	client      *gateapi.APIClient
	ctx         context.Context
	accountType string // 账户类型：spot(现货)
}

// NewClient 创建Gate.io客户端
func NewClient(cfg *config.ExchangeConfig) *Client {
	configuration := gateapi.NewConfiguration()
	
	// 创建API客户端
	client := gateapi.NewAPIClient(configuration)
	
	// 如果提供了baseURL，则修改基础URL
	if cfg.BaseURL != "" {
		client.ChangeBasePath(cfg.BaseURL)
	}
	
	// 创建带有认证信息的context
	ctx := context.WithValue(
		context.Background(),
		gateapi.ContextGateAPIV4,
		gateapi.GateAPIV4{
			Key:    cfg.ApiKey,
			Secret: cfg.ApiSecret,
		},
	)
	
	// 设置账户类型，默认为spot(现货)
	accountType := cfg.AccountType
	if accountType == "" || accountType != "spot" {
		accountType = "spot"
	}
	
	return &Client{
		client:      client,
		ctx:         ctx,
		accountType: accountType,
	}
}

// GetAccountType 获取当前账户类型
func (c *Client) GetAccountType() string {
	return c.accountType
}

// GetPositionDetail 获取资产详细信息，包括可用余额、总余额、锁定余额
func (c *Client) GetPositionDetail(currency string) (float64, float64, float64, error) {
	// 调用Gate.io API获取账户余额
	balances, _, err := c.client.SpotApi.ListSpotAccounts(c.ctx, &gateapi.ListSpotAccountsOpts{
		Currency: optional.NewString(currency),
	})
	
	if err != nil {
		return 0, 0, 0, fmt.Errorf("获取%s余额失败: %w", currency, err)
	}
	
	// 开启调试级别时才输出原始返回值
	config.Logger.Debugw("获取到原始账户余额信息",
		"currency", currency,
		"balances", balances,
	)
	
	// 如果没有找到该资产，返回0
	if len(balances) == 0 {
		return 0, 0, 0, nil
	}
	
	// 转换可用余额为float64
	available, err := strconv.ParseFloat(balances[0].Available, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("解析%s可用余额失败: %w", currency, err)
	}
	
	// 转换锁定余额为float64
	locked, err := strconv.ParseFloat(balances[0].Locked, 64)
	if err != nil {
		return available, 0, 0, fmt.Errorf("解析%s锁定余额失败: %w", currency, err)
	}
	
	// 总余额 = 可用余额 + 锁定余额
	total := available + locked
	
	config.Logger.Debugw("解析后的账户余额信息",
		"currency", currency,
		"available", available,
		"locked", locked,
		"total", total,
	)
	
	return available, total, locked, nil
}

// GetBalance 获取资产余额
func (c *Client) GetBalance(currency string) (float64, float64, error) {
	// 调用Gate.io API获取账户余额
	balances, _, err := c.client.SpotApi.ListSpotAccounts(c.ctx, &gateapi.ListSpotAccountsOpts{
		Currency: optional.NewString(currency),
	})
	
	if err != nil {
		return 0, 0, fmt.Errorf("获取%s余额失败: %w", currency, err)
	}
	
	// 开启调试级别时才输出原始返回值
	config.Logger.Debugw("获取到原始账户余额信息",
		"currency", currency,
		"balances", balances,
	)
	
	// 如果没有找到该资产，返回0
	if len(balances) == 0 {
		return 0, 0, nil
	}
	
	// 转换可用余额为float64
	available, err := strconv.ParseFloat(balances[0].Available, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("解析%s可用余额失败: %w", currency, err)
	}
	
	// 转换锁定余额为float64
	locked, err := strconv.ParseFloat(balances[0].Locked, 64)
	if err != nil {
		return available, 0, fmt.Errorf("解析%s锁定余额失败: %w", currency, err)
	}
	
	// 总余额 = 可用余额 + 锁定余额
	total := available + locked
	
	config.Logger.Debugw("解析后的账户余额信息",
		"currency", currency,
		"available", available,
		"locked", locked,
		"total", total,
	)
	
	return available, total, nil
}

// GetSymbolPrice 获取交易对价格
func (c *Client) GetSymbolPrice(symbol string) (float64, error) {
	// 使用ListTickers API获取价格
	opts := &gateapi.ListTickersOpts{
		CurrencyPair: optional.NewString(symbol),
	}
	
	tickers, _, err := c.client.SpotApi.ListTickers(c.ctx, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			return 0, fmt.Errorf("gate api error: %s - %s", e.Label, e.Message)
		}
		return 0, fmt.Errorf("获取价格失败: %w", err)
	}
	
	if len(tickers) == 0 {
		return 0, fmt.Errorf("未找到交易对 %s 的行情数据", symbol)
	}
	
	price, err := strconv.ParseFloat(tickers[0].Last, 64)
	if err != nil {
		return 0, fmt.Errorf("解析价格失败: %w", err)
	}
	
	return price, nil
}

// CreateOrder 创建订单
func (c *Client) CreateOrder(order *types.Order) (*types.OrderResponse, error) {
	// 只支持现货账户
	return c.createSpotOrder(order)
}

// createSpotOrder 创建现货订单
func (c *Client) createSpotOrder(order *types.Order) (*types.OrderResponse, error) {
	// 构建订单请求
	req := gateapi.Order{
		Text:         order.ClientID,                    // 客户端订单ID
		CurrencyPair: order.Symbol,                      // 交易对
		Side:         string(order.Side),                // buy 或 sell
		Amount:       fmt.Sprintf("%.8f", order.Amount), // 数量
		Price:        fmt.Sprintf("%.8f", order.Price),  // 价格
		Type:         "limit",                           // 限价单
	}

	// 创建订单，不需要额外的可选参数
	result, _, err := c.client.SpotApi.CreateOrder(c.ctx, req, nil)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			return nil, fmt.Errorf("gate api error: %s - %s", e.Label, e.Message)
		}
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 解析已成交数量
	filledAmount, _ := strconv.ParseFloat(result.FilledTotal, 64)

	return &types.OrderResponse{
		OrderID:   result.Id,
		Status:    result.Status,
		FilledQty: filledAmount,
	}, nil
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(symbol, orderID string) error {
	// 只支持现货账户
	return c.cancelSpotOrder(symbol, orderID)
}

// cancelSpotOrder 取消现货订单
func (c *Client) cancelSpotOrder(symbol, orderID string) error {
	// 取消订单，不需要额外的可选参数
	_, _, err := c.client.SpotApi.CancelOrder(c.ctx, orderID, symbol, nil)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			return fmt.Errorf("gate api error: %s - %s", e.Label, e.Message)
		}
		return fmt.Errorf("取消订单失败: %w", err)
	}
	return nil
}

// GetOrderStatus 获取订单状态
func (c *Client) GetOrderStatus(symbol, orderID string) (*types.OrderResponse, error) {
	// 只支持现货账户
	return c.getSpotOrderStatus(symbol, orderID)
}

// getSpotOrderStatus 获取现货订单状态
func (c *Client) getSpotOrderStatus(symbol, orderID string) (*types.OrderResponse, error) {
	// 获取订单状态，不需要额外的可选参数
	order, _, err := c.client.SpotApi.GetOrder(c.ctx, orderID, symbol, nil)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			return nil, fmt.Errorf("gate api error: %s - %s", e.Label, e.Message)
		}
		return nil, fmt.Errorf("获取订单状态失败: %w", err)
	}

	// 不在此处输出订单信息日志，避免日志过多

	// 获取成交数量（直接使用Gate.io API提供的FilledAmount字段）
	filledQty, _ := strconv.ParseFloat(order.FilledAmount, 64)
	
	// 获取成交均价
	filledPrice, _ := strconv.ParseFloat(order.AvgDealPrice, 64)
	if filledPrice == 0 && filledQty > 0 {
		// 如果没有成交均价但有成交数量，尝试使用下单价格
		filledPrice, _ = strconv.ParseFloat(order.Price, 64)
	}
	
	// 获取手续费
	fee, _ := strconv.ParseFloat(order.Fee, 64)
	
	// 正确映射订单状态
	status := order.Status
	if status == "closed" {
		// Gate.io的closed状态需要根据实际情况映射
		if filledQty > 0 {
			status = "filled"
		} else {
			status = "canceled"
		}
	}
	
	return &types.OrderResponse{
		OrderID:      order.Id,
		Status:       status,
		FilledQty:    filledQty,
		FilledPrice:  filledPrice,
		Fee:          fee,
		FeeCurrency:  order.FeeCurrency,
	}, nil
}

// GetAccountBalance 获取账户余额
func (c *Client) GetAccountBalance(currency string) (map[string]float64, error) {
	// 只支持现货账户
	return c.GetAccountBalanceByType("spot", currency)
}

// GetAccountBalanceByType 根据指定的账户类型获取账户余额
func (c *Client) GetAccountBalanceByType(accountType, currency string) (map[string]float64, error) {
	balances := make(map[string]float64)
	
	// 只支持现货账户
	if accountType != "spot" && accountType != "" {
		return nil, fmt.Errorf("不支持的账户类型: %s，目前只支持现货账户", accountType)
	}
	
	// 构建可选参数
	var opts *gateapi.ListSpotAccountsOpts
	if currency != "" {
		opts = &gateapi.ListSpotAccountsOpts{
			Currency: optional.NewString(currency),
		}
	}
	
	// 调用API获取现货账户余额
	accounts, _, err := c.client.SpotApi.ListSpotAccounts(c.ctx, opts)
	if err != nil {
		if e, ok := err.(gateapi.GateAPIError); ok {
			return nil, fmt.Errorf("gate api error: %s - %s", e.Label, e.Message)
		}
		return nil, fmt.Errorf("获取现货账户余额失败: %w", err)
	}
	
	// 如果没有余额，返回空结果
	if len(accounts) == 0 {
		if currency != "" {
			return nil, fmt.Errorf("未找到币种 %s 的现货账户余额", currency)
		}
		return balances, nil
	}
	
	// 处理返回的账户余额
	for _, account := range accounts {
		available, _ := strconv.ParseFloat(account.Available, 64)
		locked, _ := strconv.ParseFloat(account.Locked, 64)
		
		// 只返回有余额的币种
		if available > 0 || locked > 0 {
			balances[account.Currency+".available"] = available
			balances[account.Currency+".locked"] = locked
			balances[account.Currency+".total"] = available + locked
		}
	}
	
	return balances, nil
}