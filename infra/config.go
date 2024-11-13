package infra

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
)

var (
	config = &AppConfig{}
	once   sync.Once
)

type AppConfig struct {
	App App
	Log Log
}

type App struct {
	Name         string
	EchoPort     string
	GinPort      string
	Version      string
	TrustProxies []string
	Upgrader     Upgrader
}

type Log struct {
	Level string
}

type Upgrader struct {
	ReadBufferSize  int
	WriteBufferSize int
}

func InitConfig() *AppConfig {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		viper.SetDefault("app.echoPort", "8081")
		viper.SetDefault("app.ginPort", "8082")

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("[INFRA] failed to read config file:", err)
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal("[INFRA] failed to unmarshal config file:", err)
		}
	})

	return config
}
