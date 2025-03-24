# Order Go - 量化交易系统

基于Go语言和Gin框架开发的数字货币自动交易下单系统。

## 项目简介

Order Go是一个高性能的量化交易系统，支持接收交易信号，使用自定义策略进行分析，并通过多个交易所渠道（如gate.io、OKX等）执行自动下单操作。

## 技术栈

- 语言：Go
- Web框架：Gin
- 数据库：PostgreSQL
- ORM：GORM
- 配置管理：Viper
- 日志：Zap
- API文档：Swagger
- 消息队列：Redis
- 缓存：Redis
- 容器化：Docker

## 主要功能

- 接收交易信号
- 多策略管理和匹配
- 多交易所渠道管理
- 自动下单执行
- 订单状态跟踪
- 风险控制

## 项目结构

```
order_go/
├── cmd/                    # 应用程序入口
│   └── server/             # API服务器
├── configs/                # 配置文件
├── internal/               # 私有应用程序代码
│   ├── api/                # API处理程序
│   ├── exchange/           # 交易所接口
│   ├── models/             # 数据模型
│   ├── repository/         # 数据库访问层
│   ├── service/            # 业务逻辑层
│   ├── strategy/           # 交易策略实现
│   └── utils/              # 工具函数
├── pkg/                    # 可重用的库代码
├── scripts/                # 脚本和工具
├── migrations/             # 数据库迁移文件
└── docs/                   # 文档
```

## 安装和运行

### 前置条件

- Go 1.18+
- PostgreSQL 13+
- Redis (可选)

### 安装步骤

1. 克隆仓库
2. 配置环境变量
3. 初始化数据库
4. 运行应用

详细说明请参考[安装文档](./docs/installation.md)。

## 许可证

MIT