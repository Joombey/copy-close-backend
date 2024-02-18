package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type DbConfig struct {
	Driver   string `yaml:"driver"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"db_name"`
}

func (cfg DbConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.UserName, cfg.Password, cfg.Address, cfg.Port, cfg.DbName)
}

func New() (cfg DbConfig, err error) {
	cfgPath, present := os.LookupEnv("COPY_CLOSE_CONFIG")
	if !present {
		cfgPath = "local.yaml"
	}

	err = cleanenv.ReadConfig(cfgPath, &cfg)
	return
}
