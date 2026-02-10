package config

import (
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Timeout string `mapstructure:"timeout"`
}
type MySQLConfig struct {
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	DBname          string `mapstructure:"db_name"`
	DSNParams       string `mapstructure:"dsn_params"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Mysql  MySQLConfig  `mapstructure:"mysql"`
}

var (
	cfg  *Config
	once sync.Once
)

func InitConfig() *Config {
	once.Do(func() {
		configFile := "/Users/chengyue2304/webdemo/config/config.yaml"
		v := viper.New()
		v.SetConfigFile(configFile)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("read config file fail:%s", err))
		}
		cfg = new(Config)
		if err := v.Unmarshal(cfg); err != nil {
			panic(fmt.Sprintf("解析配置文件失败：%s", err.Error()))
		}
		fmt.Printf("配置加载成功，服务端口：%s\n", cfg.Server.Port)
	})
	return cfg
}
func GetConfig() *Config {
	if cfg == nil {
		return InitConfig()
	}
	return cfg
}
