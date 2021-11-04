package configs

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
)

var (
	cfg *ini.File
	confPath string
)

func init()  {
	flag.StringVar(&confPath, "conf", "configs/", "Define conf file path")
	flag.Parse()

	if confPath[len(confPath) - 1] != '/' {
		confPath = confPath + "/"
	}

	fmt.Println("Usage confPath: " + confPath)
	confPath = confPath +  "setting.ini"

	var err error
	cfg, err = ini.Load(confPath)
	if err != nil {
		log.Fatal("Fail to read file: ", err)
		panic(err)
	}
	fmt.Println("Load configs file: " + confPath)
}

func Load() {
	loadApp()
	loadDatabase()
}

func loadApp()  {
	sec, _ := cfg.GetSection("app")
	logPath := sec.Key("LogPath").MustString("/tmp/logs/")
	if logPath[len(logPath) - 1] != '/' {
		logPath = logPath + "/"
	}
	AppConfig.LogPath = logPath
	AppConfig.LogFile = sec.Key("LogFile").String()

	log.Println("Init App Config:", AppConfig)
}

func loadDatabase() {
	sec, _ := cfg.GetSection("database")

	DBConfig.Database = sec.Key("Database").String()
	DBConfig.Type = sec.Key("Type").String()
	DBConfig.Host = sec.Key("Host").String()
	DBConfig.User = sec.Key("User").String()
	DBConfig.Password = sec.Key("Password").String()
	DBConfig.MaxIdleConns = sec.Key("MaxIdleConns").MustInt(8)
	DBConfig.MaxOpenConns = sec.Key("MaxOpenConns").MustInt(16)

	log.Println("Init Database Config:", DBConfig)
}