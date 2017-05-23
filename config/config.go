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
	Postgres struct {
		Driver       string `toml:"driver"`
		User         string `toml:"user"`
		Password     string `toml:"password"`
		Port         string `toml:"port"`
		Dbname       string `toml:"dbname"`
		Host         string `toml:"host"`
		MaxOpenConns int    `toml:"maxOpenConns"`
		MaxIdleConns int    `toml:"maxIdleConns"`
		SchemaOutput string `toml:"schemaOutput"`
		TableOutput  string `toml:"tableOutput"`
	} `toml:"postgres"`
	Mongo struct {
	} `toml:"mongo"`
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
