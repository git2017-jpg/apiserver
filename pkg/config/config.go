package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var GlobalConfig *Config

// Config is application global config
type Config struct {
	Mode         string      `mapstructure:"mode"`           // gin启动模式
	Port         string      `mapstructure:"port"`           // 启动端口
	AppName      string      `mapstructure:"app-name"`       //应用名称
	Url          string      `mapstructure:"url"`            // 应用地址,用于自检 eg. http://127.0.0.1
	MaxPingCount int         `mapstructure:"max-ping-count"` // 最大自检次数，用户健康检查
	JwtSecret    string      `mapstructure:"jwt-secret"`
	Language     string      `mapstructure:"language"` // 项目语言
	DBConfig     DBConfig    `mapstructure:"database"` // 数据库信息
	RedisConfig  RedisConfig `mapstructure:"redis"`    // redis
	LogConfig    LogConfig   `mapstructure:"log"`      // uber zap
}

// DBConfig is used to configure mysql database
type DBConfig struct {
	Dbname          string `mapstructure:"dbname"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaximumPoolSize int    `mapstructure:"maximum-pool-size"`
	MaximumIdleSize int    `mapstructure:"maximum-idle-size"`
	LogMode         bool   `mapstructure:"log-mode"`
}

// RedisConfig is used to configure redis
type RedisConfig struct {
	Addr         string `mapstructure:"address"`
	Password     string `mapstructure:"password"`
	Db           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool-size"`
	MinIdleConns int    `mapstructure:"min-idle-conns"`
	IdleTimeout  int    `mapstructure:"idle-timeout"`
}

// LogConfig is used to configure uber zap
type LogConfig struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disable-caller"`
	DisableStacktrace bool   `mapstructure:"disable-stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
	Name              string `mapstructure:"name"`
	Writers           string `mapstructure:"writers"`
	LoggerFile        string `mapstructure:"logger-file"`
	LoggerWarnFile    string `mapstructure:"logger-warn-file"`
	LoggerErrorFile   string `mapstructure:"logger-error-file"`
	LogFormatText     bool   `mapstructure:"log-format-text"`
	LogRollingPolicy  string `mapstructure:"log-rolling-policy"`
	LogBackupCount    uint   `mapstructure:"log-backup-count"`
}

// InitConfig is a loader to load config file.
func InitConfig(configFilePath string) *Config {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		// 设置默认的config
		viper.AddConfigPath("configs/dev/")
		viper.SetConfigName("config")
	}

	// 初始化配置文件
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APPLICATION")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 解析到struct
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		panic(err)
	}
	log.Println("The application configuration file is loaded successfully!")

	// 监控配置文件，并热加载
	watchConfig()

	return GlobalConfig
}

// 监控配置文件变动
// 注意：有些配置修改后，及时重新加载也要重新启动应用才行，比如端口
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("Configuration file changed: %s, reload it", in.Name)
		// 忽略错误
		InitConfig(in.Name)
	})
}
