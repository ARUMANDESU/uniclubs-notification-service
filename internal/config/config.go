package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env      string `yaml:"env" env:"ENV" env-default:"local"`
	Rabbitmq `yaml:"rabbitmq"`
	Mailer   `yaml:"mailer"`
}

type Rabbitmq struct {
	User     string `yaml:"user" env:"RABBITMQ_USER"`
	Password string `yaml:"password" env:"RABBITMQ_PASSWORD"`
	Host     string `yaml:"host" env:"RABBITMQ_HOST"`
	Port     string `yaml:"port" env:"RABBITMQ_PORT"`
}

type Mailer struct {
	Host     string `yaml:"host" env:"MAILER_HOST"`
	Port     int    `yaml:"port" env:"MAILER_port"`
	Username string `yaml:"username" env:"MAILER_USERNAME"`
	Password string `yaml:"password" env:"MAILER_PASSWORD"`
	Sender   string `yaml:"sender" env:"MAILER_SENDER"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return MustLoadFromEnv()
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	cfg, err := LoadByPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func LoadByPath(configPath string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("there is no config file: %w", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func MustLoadFromEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Env empty")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
