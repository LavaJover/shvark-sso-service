package config

import (
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type SSOConfig struct {
	Env string 	`yaml:"env"`
	GRPCServer 	`yaml:"grpc_server"`
	SSODB 		`yaml:"sso_db"`
	LogConfig 	`yaml:"log_config"`
}

type GRPCServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SSODB struct {
	Dsn string `yaml:"dsn"`
}

type LogConfig struct {
	LogLevel 	string 	`yaml:"log_level"`
	LogFormat 	string 	`yaml:"log_format"`
	LogOutput 	string 	`yaml:"log_output"`
}

func MustLoad() *SSOConfig {

	// Processing env config variable and file
	configPath := os.Getenv("SSO_CONFIG_PATH")

	if configPath == ""{
		log.Fatalf("SSO_CONFIG_PATH was not found\n")
	}

	if _, err := os.Stat(configPath); err != nil{
		log.Fatalf("failed to find config file: %v\n", err)
	}

	// YAML to struct object
	var cfg SSOConfig
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil{
		log.Fatalf("failed to read config file: %v", err)
	}

	return &cfg
}