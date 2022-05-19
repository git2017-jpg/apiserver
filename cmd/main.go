package main

import (
	"github.com/spf13/pflag"
	"monitor-apiserver/internal/router"
	"monitor-apiserver/pkg/config"
	"monitor-apiserver/pkg/log"
)

var (
	cfgDir = pflag.StringP("config dir", "c", "", "config path.")
	env    = pflag.StringP("env name", "e", "", "env var name.")
)

func main() {
	// 解析启动参数
	pflag.Parse()
	// 加载配置文件
	cf := config.InitConfig(*cfgDir + *env)
	// 初始化日志
	log.Init(&cf.LogConfig)
	// 创建应用
	application := router.NewApplication(cf)
	// 运行应用
	if err := application.Run(); err != nil {
		panic(err)
	}
}
