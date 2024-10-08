package configs

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

const (
	defaultConfigPath = "./configs/config.yaml"
	defaultEnvPath    = ".env"
)

type (
	Config struct {
		App    `yaml:"app"`
		GRPC   `yaml:"http"`
		PG     `yaml:"pg"`
		RMQ    `yaml:"rabbitmq"`
		Log    `yaml:"logger"`
		Notify `yaml:"Notify"`
	}

	App struct {
		Name         string        `env:"APP_NAME"            env-default:"notification-service" yaml:"name"`
		Version      string        `env:"APP_VERSION"         env-default:"1.0.0"         yaml:"version"`
		CountWorkers int           `env:"APP_WORKERS"         env-default:"24"            yaml:"workers"`
		Timeout      time.Duration `env:"APP_TIMEOUT" env-default:"5s"           yaml:"timeout"`
	}

	GRPC struct {
		Port    string        `env:"GRPC_PORT"    env-default:":8080" yaml:"port"`
		Timeout time.Duration `env:"GRPC_TIMEOUT" env-default:"5s"    yaml:"timeout"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX" env-default:"2"     yaml:"poolMax"`
		URL     string `env:"PG_URL"      env-required:"true" yaml:"url"`
	}

	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER" env-default:"rpc_server" yaml:"rpcServerExchange"`
		ClientExchange string `env:"RMQ_RPC_CLIENT" env-default:"rpc_client" yaml:"rpcClientExchange"`
		URL            string `env:"RMQ_URL"        env-required:"true"      yaml:"url"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug" yaml:"logLevel"`
	}

	SMTP struct {
		Domain   string `yaml:"domain" env:"SMTP_DOMAIN"`
		Port     int    `yaml:"port" env:"SMTP_PORT"`
		UserName string `yaml:"user_ame" env:"SMTP_USERNAME"`
		Password string `yaml:"password" env:"SMTP_PASSWORD"`
	}

	Notify struct {
		SMTP
		Mail string `yaml:"mail" env:"MAIL"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath, defaultEnvPath)
}

func MustLoadPath(configPath, envPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	_ = godotenv.Load(envPath)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
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

	if res == "" {
		res = defaultConfigPath
	}

	return res
}
