package types

// OrderSide 订单方向
type OrderSide string

const (
    OrderSideBuy  OrderSide = "buy"
    OrderSideSell OrderSide = "sell"
)

// String 实现Stringer接口
func (s OrderSide) String() string {
    return string(s)
}

// OrderType 订单类型
type OrderType string

const (
    OrderTypeLimit  OrderType = "limit"
    OrderTypeMarket OrderType = "market"
)

// Order 统一的订单结构
type Order struct {
    Symbol    string    `json:"symbol"`     // 交易对
    Side      OrderSide `json:"side"`       // 买卖方向
    Type      OrderType `json:"type"`       // 订单类型
    Price     float64   `json:"price"`      // 价格
    Amount    float64   `json:"amount"`     // 数量
    ClientID  string    `json:"client_id"`  // 客户端订单ID
}

// OrderResponse 下单响应
type OrderResponse struct {
    OrderID     string  `json:"order_id"`      // 订单ID
    Status      string  `json:"status"`        // 订单状态
    FilledQty   float64 `json:"filled_qty"`    // 已成交数量
    FilledPrice float64 `json:"filled_price"`  // 成交均价
    Fee         float64 `json:"fee"`           // 手续费
    FeeCurrency string  `json:"fee_currency"`  // 手续费币种
    Error       error   `json:"-"`            // 错误信息
}