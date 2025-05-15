package trading

import (
	"context"
	"order_go/internal/constants"
	"order_go/internal/exchange"
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
	"sync"
	"time"
)

// OrderMonitor 订单监控器
type OrderMonitor struct {
	activeOrders sync.Map       // 当前活跃订单
	exchanges    map[string]exchange.Exchange
}

var (
	monitor     *OrderMonitor
	monitorOnce sync.Once
)

// GetOrderMonitor 获取订单监控器单例
func GetOrderMonitor() *OrderMonitor {
	monitorOnce.Do(func() {
		monitor = &OrderMonitor{
			exchanges: make(map[string]exchange.Exchange),
		}
	})
	return monitor
}

// RegisterExchange 注册交易所
func (m *OrderMonitor) RegisterExchange(name string, ex exchange.Exchange) {
	m.exchanges[name] = ex
}

// StartMonitor 开始监控订单
func (m *OrderMonitor) StartMonitor(order *models.OrderRecord, exchangeName string) {
	// 获取交易所
	ex, ok := m.exchanges[exchangeName]
	if !ok {
		config.Logger.Errorw("交易所不存在",
			"exchange", exchangeName,
			"order_id", order.OrderID,
		)
		return
	}

	// 将订单加入活跃订单列表
	m.activeOrders.Store(order.OrderID, order)

	// 启动监控协程
	go m.monitorOrder(order, ex)
}

// monitorOrder 监控订单状态
func (m *OrderMonitor) monitorOrder(order *models.OrderRecord, ex exchange.Exchange) {
	// 监控结束时从活跃订单列表中移除
	defer m.activeOrders.Delete(order.OrderID)

	// 设置监控超时时间
	ctx, cancel := context.WithTimeout(context.Background(), constants.MonitorTimeout)
	defer cancel()

	ticker := time.NewTicker(constants.MonitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 查询订单状态
			orderStatus, err := ex.GetOrderStatus(order.Symbol, order.OrderID)
			if err != nil {
				config.Logger.Errorw("查询订单状态失败",
					"error", err.Error(),
					"order_id", order.OrderID,
				)
				continue
			}

			// 更新订单状态
			updates := map[string]interface{}{
				"status": orderStatus.Status,
			}

			// 如果订单已成交，更新成交信息
			if orderStatus.Status == "filled" {
				updates["filled_price"] = orderStatus.FilledPrice
				updates["filled_amount"] = orderStatus.FilledQty
				updates["fee"] = orderStatus.Fee
				updates["fee_currency"] = orderStatus.FeeCurrency
				// 订单完成时间自动记录在updated_at字段中，不需要单独设置completed_at
			}

			// 更新数据库中的订单信息
			if err := repository.DB.Model(&models.OrderRecord{}).Where("order_id = ?", order.OrderID).Updates(updates).Error; err != nil {
				config.Logger.Errorw("更新订单信息失败",
					"error", err.Error(),
					"order_id", order.OrderID,
				)
			}

			// 如果订单已成交或已取消，结束监控
			if orderStatus.Status == "filled" || orderStatus.Status == "canceled" {
				config.Logger.Infow("订单监控结束",
					"order_id", order.OrderID,
					"status", orderStatus.Status,
					"symbol", order.Symbol,
				)
				return
			}

		case <-ctx.Done():
			// 超时前，先检查订单最新状态
			orderStatus, err := ex.GetOrderStatus(order.Symbol, order.OrderID)
			if err != nil {
				config.Logger.Errorw("超时前查询订单状态失败",
					"error", err.Error(),
					"order_id", order.OrderID,
				)
				// 查询失败，继续尝试撤单
			} else if orderStatus.Status == "filled" {
				// 订单已完全成交，不需要撤单
				config.Logger.Infow("订单已完全成交，无需撤单",
					"order_id", order.OrderID,
					"symbol", order.Symbol,
				)
				
				// 更新数据库中的订单信息
				updates := map[string]interface{}{
					"status":       "filled",
					"filled_price": orderStatus.FilledPrice,
					"filled_amount": orderStatus.FilledQty,
					"fee":          orderStatus.Fee,
					"fee_currency": orderStatus.FeeCurrency,
					// 订单完成时间自动记录在updated_at字段中
				}
				
				if err := repository.DB.Model(&models.OrderRecord{}).Where("order_id = ?", order.OrderID).Updates(updates).Error; err != nil {
					config.Logger.Errorw("更新已成交订单信息失败",
						"error", err.Error(),
						"order_id", order.OrderID,
					)
				}
				return
			} else if orderStatus.Status == "canceled" {
				// 订单已被取消，不需要再撤单
				config.Logger.Infow("订单已被取消，无需撤单",
					"order_id", order.OrderID,
					"symbol", order.Symbol,
				)
				return
			} else if orderStatus.FilledQty > 0 && orderStatus.FilledQty < order.Amount {
				// 订单部分成交，记录日志
				config.Logger.Warnw("订单部分成交，尝试撤销剩余部分",
					"order_id", order.OrderID,
					"symbol", order.Symbol,
					"filled_amount", orderStatus.FilledQty,
					"total_amount", order.Amount,
				)
			}
			
			// 尝试撤单
			config.Logger.Warnw("订单监控超时，尝试撤单",
				"order_id", order.OrderID,
				"symbol", order.Symbol,
			)

			if err := ex.CancelOrder(order.Symbol, order.OrderID); err != nil {
				config.Logger.Errorw("撤单失败",
					"error", err.Error(),
					"order_id", order.OrderID,
				)
				
				// 撤单失败，再次检查订单状态
				latestStatus, checkErr := ex.GetOrderStatus(order.Symbol, order.OrderID)
				if checkErr == nil && latestStatus.Status == "filled" {
					// 订单已成交，更新状态
					config.Logger.Infow("撤单失败但订单已成交",
						"order_id", order.OrderID,
						"symbol", order.Symbol,
					)
					
					updates := map[string]interface{}{
						"status":       "filled",
						"filled_price": latestStatus.FilledPrice,
						"filled_amount": latestStatus.FilledQty,
						"fee":          latestStatus.Fee,
						"fee_currency": latestStatus.FeeCurrency,
						// 订单完成时间自动记录在updated_at字段中
					}
					
					if updateErr := repository.DB.Model(&models.OrderRecord{}).Where("order_id = ?", order.OrderID).Updates(updates).Error; updateErr != nil {
						config.Logger.Errorw("更新已成交订单信息失败",
							"error", updateErr.Error(),
							"order_id", order.OrderID,
						)
					}
				}
			} else {
				config.Logger.Infow("撤单成功",
					"order_id", order.OrderID,
					"symbol", order.Symbol,
				)

				// 更新数据库中的订单状态
				updates := map[string]interface{}{
					"status":       "canceled",
					// 订单取消时间自动记录在updated_at字段中
				}
				
				// 如果有部分成交，记录部分成交信息
				if orderStatus != nil && orderStatus.FilledQty > 0 {
					updates["filled_price"] = orderStatus.FilledPrice
					updates["filled_amount"] = orderStatus.FilledQty
					updates["fee"] = orderStatus.Fee
					updates["fee_currency"] = orderStatus.FeeCurrency
					updates["status"] = "partially_filled"
				}
				
				if err := repository.DB.Model(&models.OrderRecord{}).Where("order_id = ?", order.OrderID).Updates(updates).Error; err != nil {
					config.Logger.Errorw("更新订单状态失败",
						"error", err.Error(),
						"order_id", order.OrderID,
					)
				}
			}

			return
		}
	}
}

// GetActiveOrders 获取当前活跃订单
func (m *OrderMonitor) GetActiveOrders() []*models.OrderRecord {
	var orders []*models.OrderRecord
	m.activeOrders.Range(func(key, value interface{}) bool {
		if order, ok := value.(*models.OrderRecord); ok {
			orders = append(orders, order)
		}
		return true
	})
	return orders
}

// CancelOrder 手动取消订单
func (m *OrderMonitor) CancelOrder(orderID string, exchangeName string) error {
	// 查找订单
	orderObj, ok := m.activeOrders.Load(orderID)
	if !ok {
		return nil // 订单不存在或已完成
	}

	order, ok := orderObj.(*models.OrderRecord)
	if !ok {
		return nil
	}

	// 获取交易所
	ex, ok := m.exchanges[exchangeName]
	if !ok {
		return nil
	}

	// 执行撤单
	if err := ex.CancelOrder(order.Symbol, order.OrderID); err != nil {
		config.Logger.Errorw("手动撤单失败",
			"error", err.Error(),
			"order_id", order.OrderID,
		)
		return err
	}

	// 更新数据库
	if err := repository.DB.Model(&models.OrderRecord{}).Where("order_id = ?", order.OrderID).Updates(map[string]interface{}{
		"status":       "canceled",
		// 订单取消时间自动记录在updated_at字段中
	}).Error; err != nil {
		config.Logger.Errorw("更新订单状态失败",
			"error", err.Error(),
			"order_id", order.OrderID,
		)
		return err
	}

	// 从活跃订单列表中移除
	m.activeOrders.Delete(orderID)

	return nil
}