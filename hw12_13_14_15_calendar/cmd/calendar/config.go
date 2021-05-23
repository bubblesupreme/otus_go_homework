package main

import (
	"fmt"

	log "github.com/bubblesupreme/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

const (
	defaultEnvString = "-1" // default value for fields which can be initialized by environment variables
	defaultEnvInt    = -1   // default value for fields which can be initialized by environment variables
)

type Config struct {
	Logger   LoggerConf
	DataBase DBMSConf
	Server   ServerConf
	Storage  StorageConf
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

type DBMSConf struct {
	Login    string `mapstructure:"login"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
}

type ServerConf struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type StorageConf struct {
	Type string `mapstructure:"type"`
}

func NewConfig() (Config, error) {
	c := Config{}
	c.DataBase.Login = defaultEnvString
	c.DataBase.DBName = defaultEnvString
	c.DataBase.Password = defaultEnvString
	c.DataBase.Port = defaultEnvInt
	c.DataBase.Host = defaultEnvString
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Error("unable to decode into struct", nil)
		return c, err
	}

	if c.DataBase.Login == defaultEnvString {
		ok := false
		c.DataBase.Login, ok = viper.Get("dblogin").(string)
		if !ok {
			return c, fmt.Errorf(
				"database login is not set. Define environment variable 'POSTGRES_USER' "+
					"or \"database\":\"login\" in the config file or check it is not equal '%s'",
				defaultEnvString)
		}
	}
	if c.DataBase.DBName == defaultEnvString {
		ok := false
		c.DataBase.DBName, ok = viper.Get("dbname").(string)
		if !ok {
			return c, fmt.Errorf(
				"database name is not set. Define environment variable 'POSTGRES_DB' "+
					"or \"database\":\"dbname\" in the config file or check it is not equal '%s'",
				defaultEnvString)
		}
	}
	if c.DataBase.Password == defaultEnvString {
		ok := false
		c.DataBase.Password, ok = viper.Get("dbpassword").(string)
		if !ok {
			return c, fmt.Errorf(
				"database login is not set. Define environment variable 'POSTGRES_PASSWORD' "+
					"or \"database\":\"password\" in the config file or check it is not equal '%s'",
				defaultEnvString)
		}
	}
	if c.DataBase.Port == defaultEnvInt {
		ok := false
		c.DataBase.Password, ok = viper.Get("dbport").(string)
		if !ok {
			return c, fmt.Errorf(
				"database login is not set. Define environment variable 'POSTGRES_PORT' "+
					"or \"database\":\"port\" in the config file or check it is not equal '%d'",
				defaultEnvInt)
		}
	}
	if c.DataBase.Host == defaultEnvString {
		ok := false
		c.DataBase.Password, ok = viper.Get("dbhost").(string)
		if !ok {
			return c, fmt.Errorf(
				"database login is not set. Define environment variable 'POSTGRES_HOST' "+
					"or \"database\":\"host\" in the config file or check it is not equal '%s'",
				defaultEnvString)
		}
	}

	return c, nil
}
