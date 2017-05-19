package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Database   databaseConfig
	File       fileConfig
	Processing processingConfig
}

type databaseConfig struct {
	Driver       string
	User         string
	Password     string
	Port         string
	Dbname       string
	Host         string
	MaxOpenConns int
	MaxIdleConns int
	Name         string
}

type fileConfig struct {
	PathName  string
	Separator string
}

type processingConfig struct {
	NumberOfThreads int
	BatchSizeCommit int
}

func GetConfig() TomlConfig {
	var config TomlConfig

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
	}

	return config
}
