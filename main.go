package main

import (
	"go-firebase-gateway/api"
	"go-firebase-gateway/common/auth"
	util "go-firebase-gateway/common/redis"
	"go-firebase-gateway/internal/middleware"
	IRedis "go-firebase-gateway/internal/redis"
	redis "go-firebase-gateway/internal/redis/driver"
	"go-firebase-gateway/repository"
	"go-firebase-gateway/service"
	"io"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type Config struct {
	Dir      string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port     string
	LogType  string
	LogLevel string
	LogFile  string
	LogAddr  string
	DB       string
}
type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	SSLMode         string
	Timeout         int
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime int
}

var config Config

func init() {
	if err := env.Parse(&config); err != nil {
		log.Error("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Dir:      config.Dir,
		Port:     viper.GetString(`main.port`),
		LogType:  viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		LogFile:  viper.GetString(`main.log_file`),
		LogAddr:  viper.GetString(`main.log_addr`),
		DB:       viper.GetString(`main.db`),
	}

	var err error
	IRedis.Redis, err = redis.NewRedis(redis.Config{
		Addr:         viper.GetString(`redis.address`),
		Password:     viper.GetString(`redis.password`),
		DB:           viper.GetInt(`redis.database`),
		PoolSize:     30,
		PoolTimeout:  20,
		IdleTimeout:  10,
		ReadTimeout:  20,
		WriteTimeout: 15,
	})
	util.CallBackListID = viper.GetString(`callback_list_requestid`)
	util.CallBackRequestHash = viper.GetString(`callback_request_hash`)
	util.HookCallStatusListID = viper.GetString(`hook_list_lead_id`)
	util.HookCallStatusHash = viper.GetString(`hook_hash_lead_id`)
	if err != nil {
		panic(err)
	}
	auth.NewAuthUtil(auth.Config{
		ExpiredTime: viper.GetInt(`oauth.expired_in`),
		TokenType:   viper.GetString(`oauth.tokenType`),
	})
	config = cfg
}

func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)
	repository.AuthRepo = repository.NewAuthRepository()
	server := api.NewServer()
	// AUTHEN API
	userService := service.NewUserService()
	api.NewAuthHandler(server.Engine, userService)
	// FIREBASE API
	firebaseService := service.NewFireBaseService()
	api.NewFireBaseHandler(server.Engine, firebaseService)
	middleware.ApiPath = "api/uaa/oauth/token?grant_type=client_credentials"
	middleware.BasicAuth = "dGVsX2RzYV9jbGllbnQ6QXdvaXVyYVNpb2ZoYW9mMDc0cnQ="
	server.Start(config.Port)
}
func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "GELF":
		gelfWriter, err := gelf.NewUDPWriter(cfg.LogAddr)
		if err != nil {
			log.Error("main", "setAppLogger", err.Error())
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.SetOutput(io.MultiWriter(os.Stdout, file, gelfWriter))
		}
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.Error("main ", "Log File "+cfg.LogFile+" error")
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}
