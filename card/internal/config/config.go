package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	// "github.com/joho/godotenv"

)

type Config struct {
	Env              string `yaml:"env" env-default:"local"`
	ConnectionString string `yaml:"connection_string" env-required:"true"`
	HTTPServer       `yaml:"http_server"`
	Secret           string        `yaml:"secret" env-required:"true"`
	Clients          ClientsConfig `yaml:"clients"`
	GRPC             GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Client struct {
	Address      string        `yaml:"address" env-required:"true"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
	STAT Client `yaml:"stat"`
}

func MustLoad() *Config {
	// if err := godotenv.Load("../../.env"); err != nil {
	// 	log.Fatal("Error loading .env file: " + err.Error())
	// }

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
