package config

import (
	"log"
	"todo-app/utils"

	"github.com/go-ini/ini"
)

type ConfigList struct {
	Port       string
	SQLDriver  string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	LogFile    string
	Static     string
}

var Config ConfigList

func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

func LoadConfig() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{
		Port:       cfg.Section("web").Key("port").String(),
		LogFile:    cfg.Section("web").Key("logfile").String(),
		SQLDriver:  cfg.Section("db").Key("driver").String(),
		DbHost:     cfg.Section("db").Key("host").String(),
		DbPort:     cfg.Section("db").Key("port").String(),
		DbUser:     cfg.Section("db").Key("user").String(),
		DbPassword: cfg.Section("db").Key("password").String(),
		DbName:     cfg.Section("db").Key("dbname").String(),
		Static:     cfg.Section("web").Key("static").String(),
	}
}
