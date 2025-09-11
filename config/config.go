// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config represents the application configuration structure.
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

// Init initializes the configuration with the specified environment.
func Init(env string) error {
	viper.SetConfigType("yml")
	viper.SetConfigName(env)
	viper.AddConfigPath("config/environments/")

	// 環境変数から設定値を上書きできるように設定
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // envの '_' が '.' として扱われる e.g) DATABASE_DSN -> Database.Dsn

	// 読み取り
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config file read error: %s", err)
	}

	// 構造体にマッピング
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal error: %s", err)
	}

	// 更新検知と再読み込み
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name) // TODO: logger output
	})

	return nil
}

// GetConfig returns the current configuration instance.
func GetConfig() *Config {
	return cfg
}
