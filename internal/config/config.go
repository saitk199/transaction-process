package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Adress       string        `yaml:"adress" env-default:"localhost:8090"`
	Timeout      time.Duration `yaml:"timeout"`
	IDDLETimeout time.Duration `yaml:"iddle_timeout"`
}

func MustLoad() *Config {
	configPath := "C:/Users/raizo/source/repos/go/Web/transaction-process/config/local.yaml"

	//check if file exists
	if _, error := os.Stat(configPath); os.IsNotExist(error) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if error := cleanenv.ReadConfig(configPath, &cfg); error != nil {
		log.Fatalf("cannot reading config: %s", error)
	}

	return &cfg
}
