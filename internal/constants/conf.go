package constants

import "time"

// 合约类型常量
const (
    ContractTypeAStock   = 1 // 大A股票
    ContractTypeFutures  = 2 // 商品期货
    ContractTypeETF      = 3 // ETF金融指数
    ContractTypeCrypto   = 4 // 虚拟货币
)

// 交易所类型常量
const (
    ExchangeTypeSpot     = "spot"    // 现货
    ExchangeTypeFutures  = "futures" // 期货
)

// 订单监控常量
const (
    // MonitorTimeout 订单监控超时时间
    MonitorTimeout = 2 * time.Minute

    // MonitorInterval 订单监控间隔时间
    MonitorInterval = 5 * time.Second
)

// 获取合约类型名称
func GetContractTypeName(contractType int) string {
    switch contractType {
    case ContractTypeAStock:
        return "大A股票"
    case ContractTypeFutures:
        return "商品期货"
    case ContractTypeETF:
        return "ETF金融指数"
    case ContractTypeCrypto:
        return "虚拟货币"
    default:
        return "未知合约类型"
    }
}

// 获取交易所类型名称
func GetExchangeTypeName(exchangeType string) string {
    switch exchangeType {
    case ExchangeTypeSpot:
        return "现货"
    case ExchangeTypeFutures:
        return "期货"
    default:
        return "未知交易所类型"
    }
}