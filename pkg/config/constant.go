package config

import "time"

const (
	MinGoVersion            = 1.16                                    // MinGoVersion 最小 Go 版本
	ProjectVersion          = "v1.2.8"                                // ProjectVersion 项目版本
	ProjectName             = "go-gin-api"                            // ProjectName 项目名称
	ProjectDomain           = "http://127.0.0.1"                      // ProjectDomain 项目域名
	ProjectPort             = ":9999"                                 // ProjectPort 项目端口
	ProjectAccessLogFile    = "./logs/" + ProjectName + "-access.log" // ProjectAccessLogFile 项目访问日志存放文件
	ProjectCronLogFile      = "./logs/" + ProjectName + "-cron.log"   // ProjectCronLogFile 项目后台任务日志存放文件
	ProjectInstallMark      = "INSTALL.lock"                          // ProjectInstallMark 项目安装完成标识
	HeaderLoginToken        = "Token"                                 // HeaderLoginToken 登录验证 Token，Header 中传递的参数
	HeaderSignToken         = "Authorization"                         // HeaderSignToken 签名验证 Authorization，Header 中传递的参数
	HeaderSignTokenDate     = "Authorization-Date"                    // HeaderSignTokenDate 签名验证 Date，Header 中传递的参数
	HeaderSignTokenTimeout  = time.Minute * 2                         // HeaderSignTokenTimeout 签名有效期为 2 分钟
	RedisKeyPrefixLoginUser = ProjectName + ":login-user:"            // RedisKeyPrefixLoginUser Redis Key 前缀 - 登录用户信息
	RedisKeyPrefixSignature = ProjectName + ":signature:"             // RedisKeyPrefixSignature Redis Key 前缀 - 签名验证信息
	ZhCN                    = "zh-cn"                                 // ZhCN 简体中文 - 中国
	EnUS                    = "en-us"                                 // EnUS 英文 - 美国
	MaxRequestsPerSecond    = 10000                                   // MaxRequestsPerSecond 每秒最大请求量
	LoginSessionTTL         = time.Hour * 24                          // LoginSessionTTL 登录有效期为 24 小时
	RequestId               = "request_id"                            // RequestId 请求id名称
	TimeLayout              = "2006-01-02 15:04:05"                   // TimeLayout 时间格式
	TimeLayoutMs            = "2006-01-02 15:04:05.000"
	UserID                  = "user_id" // UserID 用户id key
)
