package main

import (
	"fmt"
	"io/ioutil"
	"order_go/internal/account"
	"order_go/internal/api/routes"
	"order_go/internal/cache"
	"order_go/internal/database"
	"order_go/internal/exchange"
	"order_go/internal/queue"
	"order_go/internal/repository"
	"order_go/internal/strategy"
	"order_go/internal/utils/config"
	"order_go/internal/validator"
	"order_go/test"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "order_go/docs" // 导入生成的 docs 包
)

// @title           Order Go API
// @version         1.0
// @description     交易信号处理系统 API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 解析命令行参数，如果是测试模式则执行测试
	if test.ParseFlags() {
		test.RunTest()
		return
	}

	// 正常启动服务
	startServer()
}

// printAccountTotalValue 计算并输出账户总价值
func printAccountTotalValue() {
	// 创建Gate.io交易所实例
	ex := exchange.NewGateIO()
	
	// 调用账户总价值计算函数
	totalValue, err := account.GetTotalValue(ex)
	if err != nil {
		config.Logger.Errorw("计算账户总价值失败",
			"error", err.Error(),
		)
		return
	}
	
	// 输出账户总价值
	config.Logger.Infow("账户总价值",
		"total_value_usdt", fmt.Sprintf("%.2f USDT", totalValue),
	)
}

// startServer 启动正常的服务
func startServer() {
	// 初始化日志
	if err := config.InitLogger(); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}

	// 初始化数据库
	db := database.ConnectDB()
	repository.DB = db

	// 执行数据库迁移
	database.MigrateDB(db)

	// 初始化消息队列
	queue.InitSignalQueue()

	// 初始化策略管理器
	strategy.GetManager().InitStrategies()

	// 计算并输出账户总价值
	printAccountTotalValue()
	
	// 启动账户总价值缓存更新器
	cache.StartAccountValueCacheUpdater()
	
	// 校验交易对交易额度设置
	if err := validator.ValidateContractPositionRatios(); err != nil {
		config.Logger.Warnw("交易对交易额度校验失败，请检查配置",
			"error", err.Error(),
		)
		// 注意：这里不会停止系统启动，只是输出警告日志
	}

	// 设置运行模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 禁用Gin的路由日志输出
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	
	// 配置CORS中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // 允许前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	// 生产环境下配置可信代理
	if config.AppConfig.Server.Mode == "release" {
		router.SetTrustedProxies([]string{"127.0.0.1"})
	}

	// 注册路由
	routes.RegisterSignalRoutes(router)
	// 注册后台管理路由
routes.RegisterAdminRoutes(router)
	
	// 添加 Swagger 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// 添加启动日志
	config.Logger.Info("服务器启动，监听端口: " + config.AppConfig.Server.Port)
	
	// 启动服务
	router.Run(fmt.Sprintf(":%s", config.AppConfig.Server.Port))
}