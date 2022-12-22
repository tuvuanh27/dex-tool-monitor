package env

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`
}

type DbConfig struct {
	Url  string `mapstructure:"url"`
	Name string `mapstructure:"name"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type CrawlConfig struct {
	MaxBlock int `mapstructure:"max_block_per_process"`
	SafeBlock int `mapstructure:"safe_block"`
}

type Config struct {
	App    AppConfig   `mapstructure:"app"`
	Db     DbConfig    `mapstructure:"database"`
	Redis  RedisConfig `mapstructure:"redis"`
	Stable []string    `mapstructure:"stable_coin"`
	Crawl  CrawlConfig `mapstructure:"crawl"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigType("json")
	viper.AddConfigPath("../dex-tool")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	return
}

var ConfigEnv, _ = LoadConfig()
