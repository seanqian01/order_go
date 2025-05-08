package validator

import (
	"fmt"
	"order_go/internal/models"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
)

// ValidateContractPositionRatios 验证所有交易对的交易额度设置
// 确保所有激活状态的交易对的交易额度总和不超过100%
func ValidateContractPositionRatios() error {
	// 查询所有激活状态的交易对
	var contractCodes []models.ContractCode
	result := repository.DB.Where("status = ?", true).Find(&contractCodes)
	if result.Error != nil {
		return fmt.Errorf("查询交易对失败: %w", result.Error)
	}

	// 计算所有交易对的交易额度总和
	totalRatio := 0.0
	for _, contract := range contractCodes {
		totalRatio += contract.MaxPositionRatio
	}

	// 检查总和是否超过100%
	if totalRatio > 100.0 {
		err := fmt.Errorf("交易对交易额度总和(%.2f)超过了100%%，请调整交易对的MaxPositionRatio设置", totalRatio)
		config.Logger.Errorw(err.Error(),
			"total_ratio", totalRatio,
			"contract_count", len(contractCodes),
		)
		return err
	}

	// 输出当前交易对交易额度总和
	config.Logger.Infow("交易对交易额度校验通过",
		"total_ratio", fmt.Sprintf("%.2f%%", totalRatio),
		"contract_count", len(contractCodes),
	)

	// 计算剩余未分配的交易额度比例
	remainRatio := 100.0 - totalRatio
	config.Logger.Infow("剩余未分配交易额度比例",
		"remain_ratio", fmt.Sprintf("%.2f%%", remainRatio),
	)

	return nil
}
