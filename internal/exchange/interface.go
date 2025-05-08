package exchange

import (
	"fmt"
	"order_go/internal/exchange/gateio"
	"order_go/internal/exchange/types"
	"order_go/internal/models"
	"order_go/internal/utils/config"
	"strconv"
	"strings"
	"time"
)

// OrderRequest 下单请求
type OrderRequest struct {
	Symbol       string  `json:"symbol"`        // 交易对
	Price        float64 `json:"price"`         // 价格
	Amount       float64 `json:"amount"`        // 数量
	Side         string  `json:"side"`          // 买卖方向 (buy/sell)
	Type         string  `json:"type"`          // 订单类型 (limit/market)
	PositionSide string  `json:"position_side"` // 持仓方向 (open/close)
}

// OrderResponse 下单响应
type OrderResponse struct {
	OrderID     string  `json:"order_id"`      // 订单ID
	Status      string  `json:"status"`        // 订单状态
	FilledQty   float64 `json:"filled_qty"`    // 已成交数量
	FilledPrice float64 `json:"filled_price"`  // 成交均价
	Fee         float64 `json:"fee"`           // 手续费
	FeeCurrency string  `json:"fee_currency"`  // 手续费币种
}

// Exchange 交易所接口
type Exchange interface {
	// GetSymbolPrice 获取交易对价格
	GetSymbolPrice(symbol string) (float64, error)
	
	// CreateOrder 创建订单
	CreateOrder(order *OrderRequest) (*OrderResponse, error)
	
	// CancelOrder 取消订单
	CancelOrder(symbol, orderID string) error
	
	// GetOrderStatus 获取订单状态
	GetOrderStatus(symbol, orderID string) (*OrderResponse, error)
	
	// GetBalance 获取账户余额
	// 返回可用余额、总余额和错误
	GetBalance(currency string) (float64, float64, error)
	
	// GetPosition 获取持仓信息
	GetPosition(symbol string) (*models.Position, error)
}

// NewGateIO 创建GateIO交易所实例
func NewGateIO() Exchange {
	// 从配置中读取API密钥等信息
	gateCfg, exists := config.GetExchangeConfig("gateio")
	if !exists {
		// 如果配置不存在，使用空配置创建客户端
		// 在实际生产环境中应该处理这种情况
		emptyCfg := &config.ExchangeConfig{
			ApiKey:      "",
			ApiSecret:   "",
			BaseURL:     "",
			AccountType: "spot", // 默认使用现货账户
		}
		return &GateIO{client: gateio.NewClient(emptyCfg)}
	}
	return &GateIO{client: gateio.NewClient(gateCfg)}
}

// GateIO Gate.io交易所实现
type GateIO struct {
	client *gateio.Client
}

// GetClient 获取内部的Gate.io客户端
func (g *GateIO) GetClient() *gateio.Client {
	return g.client
}

// GetSymbolPrice 获取交易对价格
func (g *GateIO) GetSymbolPrice(symbol string) (float64, error) {
	return g.client.GetSymbolPrice(symbol)
}

// CreateOrder 创建订单
func (g *GateIO) CreateOrder(order *OrderRequest) (*OrderResponse, error) {
	// 转换为gateio.Client使用的Order类型
	gateOrder := &types.Order{
		Symbol:   order.Symbol,
		Side:     types.OrderSide(order.Side),
		Amount:   order.Amount,
		Price:    order.Price,
		ClientID: "t-" + strconv.FormatInt(time.Now().UnixNano(), 10), // 使用t-前缀加时间戳作为客户端ID
	}
	
	// 调用gateio.Client的CreateOrder方法
	resp, err := g.client.CreateOrder(gateOrder)
	if err != nil {
		return nil, err
	}
	
	// 转换为Exchange接口定义的OrderResponse类型
	return &OrderResponse{
		OrderID:     resp.OrderID,
		Status:      resp.Status,
		FilledQty:   resp.FilledQty,
		FilledPrice: resp.FilledPrice,
		Fee:         resp.Fee,
		FeeCurrency: resp.FeeCurrency,
	}, nil
}

// CancelOrder 取消订单
func (g *GateIO) CancelOrder(symbol, orderID string) error {
	return g.client.CancelOrder(symbol, orderID)
}

// GetOrderStatus 获取订单状态
func (g *GateIO) GetOrderStatus(symbol, orderID string) (*OrderResponse, error) {
	resp, err := g.client.GetOrderStatus(symbol, orderID)
	if err != nil {
		return nil, err
	}
	
	return &OrderResponse{
		OrderID:     resp.OrderID,
		Status:      resp.Status,
		FilledQty:   resp.FilledQty,
		FilledPrice: resp.FilledPrice,
		Fee:         resp.Fee,
		FeeCurrency: resp.FeeCurrency,
	}, nil
}

// GetBalance 获取账户余额
func (g *GateIO) GetBalance(currency string) (float64, float64, error) {
	// 获取可用余额和锁定余额
	available, _, locked, err := g.client.GetPositionDetail(currency)
	if err != nil {
		config.Logger.Warnw("获取持仓详情失败，将使用备用方法",
			"currency", currency,
			"error", err.Error(),
		)
		
		// 如果新方法失败，尝试使用旧方法
		balances, err := g.client.GetAccountBalanceByType("spot", currency)
		if err != nil {
			return 0, 0, err
		}
		
		// 获取可用余额
		available, ok := balances[currency+".available"]
		if !ok {
			return 0, 0, nil
		}
		
		// 获取总余额
		total, ok := balances[currency+".total"]
		if !ok {
			total = available // 如果没有总余额，使用可用余额
		}
		
		return available, total, nil
	}
	
	// 计算总余额
	total := available + locked
	
	config.Logger.Infow("获取到账户余额",
		"currency", currency,
		"available", available,
		"locked", locked,
		"total", total,
	)
	
	return available, total, nil
}

// GetPosition 获取持仓信息
func (g *GateIO) GetPosition(symbol string) (*models.Position, error) {
	// 从交易对中提取资产名称，例如HYPE_USDT中的HYPE
	parts := strings.Split(symbol, "_")
	if len(parts) < 1 {
		return nil, fmt.Errorf("无效的交易对格式: %s", symbol)
	}
	
	// 正确获取基础资产，对于现货交易对如HYPE_USDT，我们需要查询HYPE的余额
	asset := parts[0] // 例如HYPE_USDT中的HYPE
	
	config.Logger.Debugw("开始获取现货持仓信息",
		"symbol", symbol,
		"asset", asset,
	)
	
	// 调用client获取资产余额
	// GetBalance返回三个值：可用余额、总余额和错误
	available, _, err := g.client.GetBalance(asset)
	if err != nil {
		return nil, fmt.Errorf("获取%s余额失败: %w", asset, err)
	}
	
	// 获取locked值（真正的持仓数量）
	_, _, locked, err := g.client.GetPositionDetail(asset)
	if err != nil {
		config.Logger.Warnw("获取持仓详情失败，将使用可用余额作为持仓量",
			"asset", asset,
			"error", err.Error(),
		)
		locked = 0 // 如果获取失败，则使用零
	}
	
	// 计算总持仓量（available + locked）
	total := available + locked
	
	config.Logger.Infow("当前资产持仓状态",
		"asset", asset,
		"total", total,
	)
	
	// 如果没有持仓，返回null
	if total <= 0 {
		return nil, nil
	}
	
	// 返回模拟的持仓信息，使用total值作为持仓量
	return &models.Position{
		Symbol:     symbol,
		Size:       total,
		EntryPrice: 0, // 现货没有入场价格概念
		Leverage:   1, // 现货没有杠杆概念
		MarginType: "spot", // 标记为现货
	}, nil
}