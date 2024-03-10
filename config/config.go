package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type fileConfig struct {
	Db *DbConfig `yaml:"db"`
}

type DbConfig struct {
	Driver   string `yaml:"driver"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"db_name"`
}

var cfg = fileConfig{}

func init() {
	cfgPath, present := os.LookupEnv("COPY_CLOSE_CONFIG")
	if !present {
		path, _ := os.Getwd()
		cfgPath = fmt.Sprintf("%s/config/default.yaml", path)
	}

	cleanenv.ReadConfig(cfgPath, &cfg)
}

func GetDSN() string {
	db := cfg.Db
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.UserName, db.Password, db.Address, db.Port, db.DbName)
}
