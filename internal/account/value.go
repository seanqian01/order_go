package account

import (
	"fmt"
	"order_go/internal/exchange"
	"order_go/internal/utils/config"
	"strings"
)

// GetTotalValue 获取现货账户总价值（以USDT计价）
func GetTotalValue(ex exchange.Exchange) (float64, error) {
	// 获取所有币种的余额信息
	gateio, ok := ex.(*exchange.GateIO)
	if !ok {
		return 0, fmt.Errorf("不支持的交易所类型，无法获取账户总价值")
	}
	
	// 获取账户余额
	balances, err := gateio.GetClient().GetAccountBalance("")
	if err != nil {
		return 0, fmt.Errorf("获取账户余额失败: %w", err)
	}
	
	// 初始化总价值
	totalValue := 0.0
	
	// 遍历所有币种余额
	for key, value := range balances {
		// 只处理 total 余额
		if strings.HasSuffix(key, ".total") {
			currency := strings.TrimSuffix(key, ".total")
			
			// 如果是 USDT，直接加到总价值中
			if currency == "USDT" {
				totalValue += value
				config.Logger.Debugw("添加USDT余额到总价值",
					"currency", currency,
					"amount", value,
					"total_value", totalValue,
				)
			} else if value > 0 {
				// 对于非 USDT 币种且余额大于0，获取其 USDT 价格并计算价值
				price, err := ex.GetSymbolPrice(currency + "_USDT")
				if err != nil {
					// 如果获取价格失败，记录日志但继续处理其他币种
					config.Logger.Warnw("获取币种价格失败，跳过该币种",
						"currency", currency,
						"error", err.Error(),
					)
					continue
				}
				
				// 计算该币种的 USDT 价值并加到总价值中
				currencyValue := value * price
				totalValue += currencyValue
				
				config.Logger.Debugw("添加非USDT币种到总价值",
					"currency", currency,
					"amount", value,
					"price", price,
					"value", currencyValue,
					"total_value", totalValue,
				)
			}
		}
	}
	
	config.Logger.Infow("计算账户总价值完成",
		"total_value", totalValue,
	)
	
	return totalValue, nil
}
