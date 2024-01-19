package config

import (
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env"
	"github.com/spf13/viper"
)

var (
	cfg         Config
	initialized bool = false
)

type Config struct {
	PostgresConnectionString string `env:"DB_CONNECTION_STRING"`
	Port                     int    `env:"PORT" envDefault:"4000"`
	IsDebug                  bool   `env:"IS_DEBUG" envDefault:"false"`
	DefaultPage              int    `mapstructure:"default_page"`
	DefaultPerPage           int    `mapstructure:"default_per_page"`
	DefaultSortField         string `mapstructure:"default_sort_field"`
	DefaultSortOrder         string `mapstructure:"default_sort_order"`
}

func getCurrentPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func fixFilePath(path string) string {
	fixed := ""
	for _, v := range path {
		if v == '\\' {
			fixed += "/"
			continue
		}
		fixed += string(v)
	}

	return fixed
}

func BindConfig(configFilename string) {
	cfg = Config{}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	v := viper.New()
	path := fixFilePath(getCurrentPath() + "\\" + configFilename)
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	err := v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	initialized = true
}

func GetConfig() Config {
	if !initialized {
		panic("Config is not initialized. Call BindConfig() before calling GetConfig()")
	}
	return cfg
}
