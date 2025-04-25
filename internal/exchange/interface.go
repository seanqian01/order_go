package exchange

import (
	"order_go/internal/exchange/types"
)

// Exchange 交易所接口
type Exchange interface {
    // 市场数据
    GetSymbolPrice(symbol string) (float64, error)
    
    // 交易功能
    CreateOrder(order *types.Order) (*types.OrderResponse, error)
    CancelOrder(symbol, orderID string) error
    GetOrderStatus(symbol, orderID string) (*types.OrderResponse, error)
    
    // 账户功能
    GetBalance(currency string) (float64, error)
}