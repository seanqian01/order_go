package models

import "time"

// OrderParams 下单参数
type OrderParams struct {
	Symbol       string  `json:"symbol"`        // 交易对
	Price        float64 `json:"price"`         // 价格
	Action       string  `json:"action"`        // 交易动作 (buy/sell)
	OrderType    string  `json:"order_type"`    // 订单类型 (limit/market)
	PositionSide string  `json:"position_side"` // 持仓方向 (open/close)
	Amount       float64 `json:"amount"`        // 下单数量
}

// OrderRecord 订单记录模型
type OrderRecord struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SystemOrderID  string    `json:"system_order_id" gorm:"uniqueIndex"` // 系统订单号（12位，包含类型和日期信息）
	OrderID        string    `json:"order_id" gorm:"uniqueIndex"`       // 交易所订单ID
	StrategyID     uint      `json:"strategy_id"`                       // 关联的策略ID
	ExchangeID     uint      `json:"exchange_id"`                       // 交易所ID
	Symbol         string    `json:"symbol"`                            // 交易对
	ContractCode   string    `json:"contract_code"`                     // 合约代码
	ContractType   string    `json:"contract_type"`                     // 合约类型 (spot/futures)
	OrderType      string    `json:"order_type"`                        // 订单类型 (limit/market)
	Price          float64   `json:"price"`                             // 价格
	Amount         float64   `json:"amount"`                            // 数量
	Action         string    `json:"action"`                            // 交易动作 (buy/sell)
	PositionSide   string    `json:"position_side"`                     // 持仓方向 (open/close)
	Status         string    `json:"status"`                            // 订单状态 (created/pending/filled/canceled/failed)
	FilledPrice    float64   `json:"filled_price"`                      // 成交价格
	FilledAmount   float64   `json:"filled_amount"`                     // 成交数量
	Fee            float64   `json:"fee"`                               // 手续费
	FeeCurrency    string    `json:"fee_currency"`                      // 手续费币种
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"` // 订单更新时间，完成时也会更新
}

// TableName 指定表名
func (OrderRecord) TableName() string {
	return "order_records"
}

// Position 持仓模型
type Position struct {
	Symbol     string  `json:"symbol"`      // 交易对
	Size       float64 `json:"size"`        // 持仓量 (正数为多仓，负数为空仓)
	EntryPrice float64 `json:"entry_price"` // 入场价格
	Leverage   int     `json:"leverage"`    // 杠杆倍数
	MarginType string  `json:"margin_type"` // 保证金类型 (isolated/cross)
}
