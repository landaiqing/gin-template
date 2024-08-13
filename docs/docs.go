// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/user/List": {
            "get": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "获取所有用户列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/user/delete": {
            "delete": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "删除用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户uuid",
                        "name": "uuid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/user/login": {
            "post": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "账号登录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "账号",
                        "name": "account",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/user/query_by_phone": {
            "get": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "根据手机号查询用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "手机号",
                        "name": "phone",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/user/query_by_username": {
            "get": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "根据用户名查询用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/user/query_by_uuid": {
            "get": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "根据uuid查询用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户uuid",
                        "name": "uuid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/user/register": {
            "post": {
                "tags": [
                    "鉴权模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "用户信息",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ScaAuthUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/rotate/check": {
            "post": {
                "description": "验证旋转验证码",
                "tags": [
                    "旋转验证码"
                ],
                "summary": "验证旋转验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "验证码角度",
                        "name": "angle",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "验证码key",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/rotate/get": {
            "get": {
                "description": "生成旋转验证码",
                "tags": [
                    "旋转验证码"
                ],
                "summary": "生成旋转验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/shape/check": {
            "get": {
                "description": "验证点击形状验证码",
                "tags": [
                    "点击形状验证码"
                ],
                "summary": "验证点击形状验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/shape/get": {
            "get": {
                "description": "生成点击形状验证码",
                "tags": [
                    "点击形状验证码"
                ],
                "summary": "生成点击形状验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/shape/slide/check": {
            "get": {
                "description": "验证点击形状验证码",
                "tags": [
                    "点击形状验证码"
                ],
                "summary": "验证点击形状验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "点击坐标",
                        "name": "point",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "验证码key",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/shape/slide/region/get": {
            "get": {
                "description": "验证点击形状验证码",
                "tags": [
                    "点击形状验证码"
                ],
                "summary": "验证点击形状验证码",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/text/check": {
            "get": {
                "description": "验证基础文字验证码",
                "tags": [
                    "基础文字验证码"
                ],
                "summary": "验证基础文字验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "验证码",
                        "name": "captcha",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "验证码key",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/captcha/text/get": {
            "get": {
                "description": "生成基础文字验证码",
                "tags": [
                    "基础文字验证码"
                ],
                "summary": "生成基础文字验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "验证码类型",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/sms/ali/send": {
            "get": {
                "description": "发送短信验证码",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "短信验证码"
                ],
                "summary": "发送短信验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "手机号",
                        "name": "phone",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/sms/smsbao/send": {
            "get": {
                "description": "发送短信验证码",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "短信验证码"
                ],
                "summary": "发送短信验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "手机号",
                        "name": "phone",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "dto.ScaAuthUser": {
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "头像",
                    "type": "string"
                },
                "blog": {
                    "description": "博客",
                    "type": "string"
                },
                "company": {
                    "description": "公司",
                    "type": "string"
                },
                "created_by": {
                    "description": "创建人",
                    "type": "string"
                },
                "created_time": {
                    "description": "创建时间",
                    "type": "string"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "gender": {
                    "description": "性别",
                    "type": "string"
                },
                "introduce": {
                    "description": "介绍",
                    "type": "string"
                },
                "location": {
                    "description": "地址",
                    "type": "string"
                },
                "nickname": {
                    "description": "昵称",
                    "type": "string"
                },
                "phone": {
                    "description": "电话",
                    "type": "string"
                },
                "status": {
                    "description": "状态 0 正常 1 封禁",
                    "type": "integer"
                },
                "update_by": {
                    "description": "更新人",
                    "type": "string"
                },
                "update_time": {
                    "description": "更新时间",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                },
                "uuid": {
                    "description": "唯一ID",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
