// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/webhook": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "接收并处理交易信号，同时将信号发送到处理队列和存储队列",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "signals"
                ],
                "summary": "接收交易信号",
                "parameters": [
                    {
                        "description": "交易信号",
                        "name": "signal",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TradingSignal"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "信号处理成功",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "密钥无效",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "503": {
                        "description": "服务不可用",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.TradingSignal": {
            "type": "object",
            "required": [
                "action",
                "alert_title",
                "contractType",
                "price",
                "scode",
                "strategy_id",
                "symbol",
                "time_circle"
            ],
            "properties": {
                "action": {
                    "description": "交易动作",
                    "type": "string",
                    "example": "buy"
                },
                "alert_title": {
                    "description": "提醒标题",
                    "type": "string",
                    "example": "BTC买入信号"
                },
                "contractType": {
                    "description": "合约类型",
                    "type": "string",
                    "example": "spot"
                },
                "created_at": {
                    "description": "创建时间",
                    "type": "string",
                    "example": "2025-04-28T09:00:00+08:00"
                },
                "id": {
                    "description": "信号ID",
                    "type": "integer",
                    "example": 1
                },
                "price": {
                    "description": "价格",
                    "type": "number",
                    "example": 50000
                },
                "scode": {
                    "description": "交易对简码",
                    "type": "string",
                    "example": "BTC"
                },
                "secretkey": {
                    "description": "API密钥，不存储到数据库",
                    "type": "string",
                    "example": "your-secret-key"
                },
                "strategy_id": {
                    "description": "策略ID",
                    "type": "string",
                    "example": "1"
                },
                "symbol": {
                    "description": "交易对",
                    "type": "string",
                    "example": "BTC_USDT"
                },
                "time_circle": {
                    "description": "时间周期",
                    "type": "string",
                    "example": "5m"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "X-API-Key",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Order Go API",
	Description:      "交易信号处理系统 API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
