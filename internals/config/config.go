package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr        string        `yaml:"address" env-default:"localhost:8081"`
	Timeout     time.Duration `yaml:"timeout"  env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"20s"`
	Username    string        `yaml:"username" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Config struct {
	Env         string               `yaml:"env" env-required:"true"`
	StoragePath string               `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"` //nameless struct
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Unable to get config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) { // no compare with nil
		log.Fatalf("Invalid config path: %s", configPath)
	}

	fmt.Println(configPath)

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Unable to parse the config: %s", err)
	}

	return &cfg
}
