package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
)

var (
	C     *Config
	Redis RedisConfig
	Web   WebConfig
	Bsz   BszConfig
)

var (
	configPath string
	DistPath   string
)

var defaultConfig = Config{
	Redis: RedisConfig{
		Address:   "redis:6379",
		Password:  "",
		Database:  0,
		Prefix:    "bsz_",
		MaxIdle:   10,
		MaxActive: 100,
	},
	Web: WebConfig{
		Address: "0.0.0.0:8080",
		Cors:    "*",
		Log:     true,
		Debug:   false,
	},
	Bsz: BszConfig{
		Expire:    0,
		JwtSecret: "bsz",
	},
}

func init() {
	// get config file path
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.StringVar(&DistPath, "d", "dist", "dist path")
	flag.Parse()

	var data []byte
	var err error
	C = &Config{}

	if data, err = os.ReadFile(configPath); err == nil {
		if err = yaml.Unmarshal(data, C); err != nil {
			log.Fatal("Error parsing config file:\r\n" + err.Error())
		}
	} else {
		log.Println("Error reading config file:\r\n" + err.Error())
		log.Println("Using default config", defaultConfig)
		C = &defaultConfig
	}

	// READ FROM ENV
	if _redisAddr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		C.Redis.Address = _redisAddr
	}
	if _redisPwd, ok := os.LookupEnv("REDIS_PWD"); ok {
		C.Redis.Password = _redisPwd
	}
	if _redisDb, ok := os.LookupEnv("REDIS_DATABASE"); ok {
		if _redisDbInt, err := strconv.Atoi(_redisDb); err == nil {
			C.Redis.Database = _redisDbInt
		}
	}
	if _log, ok := os.LookupEnv("LOG_ENABLE"); ok {
		if _logBool, err := strconv.ParseBool(_log); err == nil {
			C.Web.Log = _logBool
		}
	}
	if _debug, ok := os.LookupEnv("DEBUG_ENABLE"); ok {
		if _debugBool, err := strconv.ParseBool(_debug); err == nil {
			C.Web.Debug = _debugBool
		}
	}
	if _jwt, ok := os.LookupEnv("JWT_SECRET"); ok {
		C.Bsz.JwtSecret = _jwt
	}

	Redis = C.Redis
	Web = C.Web
	Bsz = C.Bsz
}
