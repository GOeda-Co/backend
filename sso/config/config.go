package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {  
    Env            string     `yaml:"env" env-default:"local"`  
    ConnectionString    string     `yaml:"connection_string" env-required:"true"`  
    GRPC           GRPCConfig `yaml:"grpc"`  
    MigrationsPath string  
    TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`  
}  

type GRPCConfig struct {  
    Address    string           `yaml:"address" env-default:":44044"`  
    Timeout time.Duration `yaml:"timeout"`  
}

func MustLoad() *Config {  
    configPath := fetchConfigPath()  
    if configPath == "" {  
        panic("config path is empty") 
    }  

    // check if file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        panic("config file does not exist: " + configPath)
    }

    var cfg Config

    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        panic("config path is empty: " + err.Error())
    }

    return &cfg
}

func fetchConfigPath() string {
    var res string

    flag.StringVar(&res, "config", "", "path to config file")
    flag.Parse()
    
    if res == "" {
        fmt.Println("sdsdsd")
        res = os.Getenv("CONFIG_PATH")
    }
    fmt.Println(res)

    return res
}