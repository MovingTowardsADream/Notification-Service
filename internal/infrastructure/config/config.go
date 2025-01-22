package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	_defaultEnvPath    = ".env"
	_defaultConfigPath = "./configs/prod.yaml"
)

type (
	Config struct {
		App           `yaml:"app"`
		GRPC          `yaml:"grpc"`
		Storage       `yaml:"storage"`
		Messaging     `yaml:"messaging"`
		SMTP          `yaml:"smtp"`
		PhoneSender   `yaml:"phone_sender"`
		Observability `yaml:"observability"`
		Log           `yaml:"logger"`
	}

	App struct {
		Name         string        `env:"APP_NAME"            env-default:"notification-service" yaml:"name"`
		Version      string        `env:"APP_VERSION"         env-default:"1.0.0"         yaml:"version"`
		CountWorkers int           `env:"APP_WORKERS"         env-default:"24"            yaml:"countWorkers"`
		Timeout      time.Duration `env:"APP_TIMEOUT" env-default:"5s"           yaml:"timeout"`
	}

	GRPC struct {
		Port    string        `env:"GRPC_PORT"    env-default:":8080" yaml:"port"`
		Timeout time.Duration `env:"GRPC_TIMEOUT" env-default:"5s"    yaml:"timeout"`
	}

	Storage struct {
		PoolMax      int32         `env:"PG_POOL_MAX" env-default:"2"     yaml:"poolMax"`
		URL          string        `env:"PG_URL"      env-required:"true" yaml:"url"`
		ConnAttempts int           `env:"PG_CONN_ATTEMPTS" yaml:"connAttempts"`
		ConnTimeout  time.Duration `env:"PG_CONN_TIMEOUT" yaml:"connTimeout"`
	}

	MessagingServer struct {
		RPCExchange     string        `yaml:"rpcExchange"`
		GoroutinesCount int           `yaml:"goroutinesCount"`
		WaitTime        time.Duration `yaml:"waitTime"`
		Attempts        int           `yaml:"attempts"`
		Timeout         time.Duration `yaml:"timeout"`
	}

	MessagingClient struct {
		RPCExchange string        `yaml:"rpcExchange"`
		WaitTime    time.Duration `yaml:"waitTime"`
		Attempts    int           `yaml:"attempts"`
		Timeout     time.Duration `yaml:"timeout"`
	}

	Messaging struct {
		Server MessagingServer `yaml:"server"`
		Client MessagingClient `yaml:"client"`
		URL    string          `env:"RMQ_URL"        env-required:"true"      yaml:"url"`
		Topics []string        `env:"RMQ_TOPICS"     env-required:"true"      yaml:"topics"`
	}

	SMTP struct {
		Domain   string `yaml:"domain" env:"SMTP_DOMAIN"`
		Port     int    `yaml:"port" env:"SMTP_PORT"`
		UserName string `yaml:"userName" env:"SMTP_USERNAME"`
		Password string `yaml:"password" env:"SMTP_PASSWORD"`
	}

	PhoneSender struct {
		AccountSID          string `env:"TWILIO_ACCOUNT_SID" env-required:"true" yaml:"accountSID"`
		AuthToken           string `env:"TWILIO_AUTH_TOKEN" env-required:"true" yaml:"authToken"`
		MessagingServiceSID string `env:"TWILIO_MESSAGING_SERVICE_SID" env-required:"true" yaml:"messagingServiceSID"`
	}

	RetryConfig struct {
		Enabled         bool          `yaml:"enabled"`
		InitialInterval time.Duration `yaml:"initialInterval"`
		MaxInterval     time.Duration `yaml:"maxInterval"`
		MaxElapsedTime  time.Duration `yaml:"maxElapsedTime"`
	}

	Trace struct {
		Host        string `env:"TRACE_PROVIDER_HOST" env-required:"true"`
		Port        string `yaml:"port"`
		RetryConfig `yaml:"retryConfig"`
	}

	Observability struct {
		Trace `yaml:"trace"`
	}

	Log struct {
		Level string  `env:"LOG_LEVEL" env-default:"testing" yaml:"logLevel"`
		Path  *string `env:"LOG_PATH" yaml:"logPath"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath, _defaultEnvPath)
}

func MustLoadPath(configPath, envPath string) *Config {
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

	flag.StringVar(&res, "config", _defaultConfigPath, "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
