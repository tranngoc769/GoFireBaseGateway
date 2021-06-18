package main

import (
	"context"
	"go-firebase-gateway/api"
	"go-firebase-gateway/common/auth"
	IRedis "go-firebase-gateway/internal/redis"
	redis "go-firebase-gateway/internal/redis/driver"
	"go-firebase-gateway/repository"
	"go-firebase-gateway/service"
	"io"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
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
	Firebase
}

var config Config

type Firebase struct {
	DatabaseURL string
	ConfigFile  string
	Document    string
}

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
	firebaseConfig := Firebase{
		DatabaseURL: viper.GetString(`firebase.database_url`),
		ConfigFile:  viper.GetString(`firebase.config_file`),
		Document:    viper.GetString(`firebase.database_document`),
	}
	cfg := Config{
		Dir:      config.Dir,
		Port:     viper.GetString(`main.port`),
		LogType:  viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		LogFile:  viper.GetString(`main.log_file`),
		LogAddr:  viper.GetString(`main.log_addr`),
		DB:       viper.GetString(`main.db`),
		Firebase: firebaseConfig,
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
	// FIREBASE CONFIG
	repository.FireBaseContext = context.Background()
	conf := &firebase.Config{
		DatabaseURL: config.Firebase.DatabaseURL,
	}
	opt := option.WithCredentialsFile(config.Firebase.ConfigFile)
	app, err := firebase.NewApp(repository.FireBaseContext, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing firebase app:", err)
	}
	client, err := app.Database(repository.FireBaseContext)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	repository.EventRef = client.NewRef(config.Firebase.Document)
	// END
	repository.AuthRepo = repository.NewAuthRepository()
	server := api.NewServer()
	// AUTHEN API
	userService := service.NewUserService()
	api.NewAuthHandler(server.Engine, userService)
	// FIREBASE API
	firebaseService := service.NewFireBaseService()
	api.NewFireBaseHandler(server.Engine, firebaseService)

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
