package orderid

import (
	"fmt"
	"math/rand"
	"order_go/internal/utils/config"
	"strings"
	"sync"
	"time"
)

const (
	// 订单号总长度
	OrderIDLength = 12
	
	// 订单类型标识
	TypeCrypto = 'C' // 虚拟货币订单
	TypeOther  = 'O' // 其他类型订单
	
	// 分区标识起始值
	InitialPartition = 'A'
	
	// 序列号最大值
	MaxSequence = 999
)

var (
	// 用于保护序列号生成的锁
	sequenceMutex sync.Mutex
	
	// 随机数生成器
	randSource = rand.NewSource(time.Now().UnixNano())
	randGen    = rand.New(randSource)
	
	// 当前日期的序列号映射
	// 格式：MMDD_X -> 当前序列号
	// 例如：0516_A -> 1234
	currentSequences = make(map[string]int)
	
	// 当前日期的分区标识映射
	// 格式：MMDD -> 当前分区标识
	// 例如：0516 -> 'A'
	currentPartitions = make(map[string]byte)
)

// GenerateOrderID 生成系统订单号
// 格式: C + YY(年份) + MM(月份) + DD(日期) + RR(2位随机数) + X(分区标识A-Z) + 3位自增数字
// 例如: C25051642A001
func GenerateOrderID(orderType byte) (string, error) {
	// 获取当前日期（年月日）
	now := time.Now()
	yearStr := now.Format("06")    // YY格式（年份后两位）
	dateStr := now.Format("0102") // MMDD格式
	
	// 加锁保护序列号生成
	sequenceMutex.Lock()
	defer sequenceMutex.Unlock()
	
	// 生成年月日的完整键
	fullDateKey := yearStr + dateStr
	
	// 获取当前日期的分区标识，如果不存在则初始化为'A'
	partition, exists := currentPartitions[fullDateKey]
	if !exists {
		partition = InitialPartition
		currentPartitions[fullDateKey] = partition
	}
	
	// 生成年月日+分区的键
	datePartitionKey := fmt.Sprintf("%s_%c", fullDateKey, partition)
	
	// 获取当前序列号，如果不存在则初始化为0
	sequence, exists := currentSequences[datePartitionKey]
	if !exists {
		sequence = 0
		currentSequences[datePartitionKey] = sequence
	}
	
	// 序列号自增
	sequence++
	
	// 如果序列号超过最大值，则分区标识自增
	if sequence > MaxSequence {
		// 分区标识自增
		partition++
		
		// 检查分区标识是否超过'Z'
		if partition > 'Z' {
			config.Logger.Errorw("分区标识超过最大值'Z'",
				"date", dateStr,
			)
			return "", fmt.Errorf("当日订单数量超过系统支持的最大值")
		}
		
		// 更新当前分区标识
		currentPartitions[fullDateKey] = partition
		
		// 创建新的年月日+分区键，并重置序列号
		datePartitionKey = fmt.Sprintf("%s_%c", fullDateKey, partition)
		sequence = 1
	}
	
	// 更新当前序列号
	currentSequences[datePartitionKey] = sequence
	
	// 生成两位随机数
	randomPart := fmt.Sprintf("%02d", randGen.Intn(100))
	
	// 生成订单号
	orderID := fmt.Sprintf("%c%s%s%s%c%03d", orderType, yearStr, dateStr, randomPart, partition, sequence)
	
	return orderID, nil
}

// 清理过期的序列号缓存
// 可以定期调用此函数清理旧的缓存条目
func CleanupExpiredSequences() {
	// 获取当前年月日
	currentDate := time.Now().Format("060102")
	
	sequenceMutex.Lock()
	defer sequenceMutex.Unlock()
	
	// 清理序列号缓存
	for key := range currentSequences {
		// 提取日期部分
		parts := strings.Split(key, "_")
		if len(parts) > 0 && parts[0] != currentDate {
			delete(currentSequences, key)
		}
	}
	
	// 清理分区标识缓存
	for date := range currentPartitions {
		if date != currentDate {
			delete(currentPartitions, date)
		}
	}
}
