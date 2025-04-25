package main

import (
	"fmt"
	"io/ioutil"
	"order_go/internal/api/routes"
	"order_go/internal/database"
	"order_go/internal/queue"
	"order_go/internal/repository"
	"order_go/internal/utils/config"
	"order_go/test"

	"github.com/gin-gonic/gin"
)

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

	// 设置运行模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 禁用Gin的路由日志输出
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	
	// 生产环境下配置可信代理
	if config.AppConfig.Server.Mode == "release" {
		router.SetTrustedProxies([]string{"127.0.0.1"})
	}

	// 注册路由
	routes.RegisterSignalRoutes(router)
	
	// 添加启动日志
	config.Logger.Info("服务器启动，监听端口: " + config.AppConfig.Server.Port)
	
	// 启动服务
	router.Run(fmt.Sprintf(":%s", config.AppConfig.Server.Port))
}