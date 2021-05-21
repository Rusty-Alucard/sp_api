package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// マッピング用の構造体
type Config struct {
	Database struct {
		Driver   string
		User     string
		Password string
		Host     string
		Port     int
		Db       string
	}
}

var cfg *Config

func Init(env string) error {
	viper.SetConfigType("yml")
	viper.SetConfigName(env)
	viper.AddConfigPath("config/environments/")

	// 環境変数から設定値を上書きできるように設定
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // envの '_' が '.' として扱われる e.g) DATABASE_DSN -> Database.Dsn

	// 読み取り
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config file read error : %s \n", err)
	}

	// 構造体にマッピング
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal error : %s \n", err)
	}

	// 更新検知と再読み込み
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name) // TODO: logger output
	})

	return nil
}

func GetConfig() *Config {
	return cfg
}
