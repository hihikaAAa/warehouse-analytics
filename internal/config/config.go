package config

import(
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{
	Env string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path"` // для локального sqlite
	HTTPServer `yaml:"http_server"`
	DB `yaml:"db"`
	WB `yaml:"wb"`
	Ozon `yaml:"ozon"`
	Alerts `yaml:"alerts"`
}

type HTTPServer struct{
	Address string `yaml:"address" env-default:"localhost:8081"`
	ReadTimeout time.Duration `yaml:"rtimeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"wtimeout" env-default:"6s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"90s"`
	//User string `yaml:"user" env-required:"true"`
	//Password string `yaml:"password" env-required:"true"`
}

type DB struct{
	DSN string `yaml:"dsn" env:"POSTGRES_DSN"`
}

type WB struct{
	Token string `yaml:"token" env:"WB_TOKEN"`
}

type Ozon struct{
	ClientID string `yaml:"client_id" env:"OZON_CLIENT_ID"`
	APIKey string `yaml:"api_key" env:"OZON_API_KEY"`
}

type Alerts struct{
	LowStockDays int `yaml:"low_stock_days" env-default:"7"`
}

func MustLoad() *Config{
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		log.Fatal("There is no CONFIG_PATH")
	}

	if _,err:= os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath,&cfg); err!=nil{
		log.Fatalf("cannot read config file: %s",err)
	}
	
	if err := cleanenv.ReadEnv(&cfg); err != nil { // Позволяем env-переменным переопределять значения из файла
		log.Fatalf("cannot read env: %v", err)
	}
	return &cfg
}