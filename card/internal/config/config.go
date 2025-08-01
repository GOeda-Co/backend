package config

import (
	"fmt"
	// "log"
	"os"
	"time"

	// "github.com/ilyakaznacheev/cleanenv"
	// "github.com/tomatoCoderq/card/internal/config"
	"gopkg.in/yaml.v3"
	// "github.com/joho/godotenv"
)

type Config struct {
	Env              string `yaml:"env" env-default:"local"`
	ConnectionString string `yaml:"connection_string" env-required:"true"`
	// HTTPServer       `yaml:"http_server"`
	Secret  string        `yaml:"secret" env-required:"true"`
	Clients ClientsConfig `yaml:"clients"`
	GRPC    GRPCConfig    `yaml:"grpc"`
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
	// SSO Client `yaml:"sso"`
	STAT Client `yaml:"stat"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		possiblePaths := []string{
			"config/config.yaml",
			"../../config/local.yaml",
			"../../config/config.yaml",
			"/app/config/config.yaml",
			"card/config/config.yaml",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
		if configPath == "" {
			panic(fmt.Errorf("could not find config file in any of the expected locations"))
		}
	}

	fmt.Println("Using config file:", configPath)

	// Read and expand config file
	raw, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("could not read config file %s: %w", configPath, err))
	}
	expanded := os.ExpandEnv(string(raw))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		panic(fmt.Errorf("could not unmarshal config file %s: %w", configPath, err))
	}

	fmt.Println(cfg.ConnectionString)
	// configPath := os.Getenv("CONFIG_PATH")
	// if configPath == "" {
	// 	log.Fatal("CONFIG_PATH is not set")
	// }
	// // check if file exists
	// if _, err := os.Stat(configPath); os.IsNotExist(err) {
	// 	log.Fatalf("config file does not exist: %s", configPath)
	// }

	// var cfg Config

	return &cfg
}
