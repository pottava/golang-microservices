// Package config defines application's configurations
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/pottava/golang-micro-services/app-authentication/app/misc"
)

func defaultConfig() Config {
	return Config{
		Name:            "GoMicroServices-Authentication",
		Port:            80,
		LogLevel:        4,
		AccessLog:       true,
		CorsMethods:     "",
		CorsOrigin:      "*",
		SecuredCookie:   false,
		TwitterKey:      "",
		TwitterSecret:   "",
		TwitterCallback: "",
	}
}

// NewConfig returns a config struct created by
// merging environment variables and a config file.
func NewConfig() *Config {
	temp := environmentConfig()
	config := &temp

	if !config.complete() {
		config.merge(fileConfig())
	}
	defer func() {
		config.merge(defaultConfig())
		config.trimWhitespace()
	}()
	return config
}

func environmentConfig() Config {
	return Config{
		Name:            os.Getenv("APP_NAME"),
		Host:            os.Getenv("APP_HOST"),
		Port:            misc.ParseUint16(os.Getenv("APP_PORT")),
		LogLevel:        misc.Atoi(os.Getenv("APP_LOG_LEVEL")),
		AccessLog:       misc.ParseBool(os.Getenv("APP_ACCESS_LOG")),
		CorsMethods:     os.Getenv("APP_ACCESS_CONTROL_ALLOW_METHODS"),
		CorsOrigin:      os.Getenv("APP_ACCESS_CONTROL_ALLOW_ORIGIN"),
		SecuredCookie:   misc.ParseBool(os.Getenv("APP_SECURED_COOKIE")),
		TwitterKey:      os.Getenv("APP_TWITTER_CONSUMER_KEY"),
		TwitterSecret:   os.Getenv("APP_TWITTER_CONSUMER_SECRET"),
		TwitterCallback: os.Getenv("APP_TWITTER_CONSUMER_CALLBACK"),
	}
}

func toStringArray(values string) []string {
	if misc.ZeroOrNil(values) {
		return []string{}
	}
	return strings.Split(values, ",")
}

func fileConfig() Config {
	path := misc.NVL(os.Getenv("CONFIG_FILE_PATH"), "/etc/golang-micro-services/config.json")
	file, err := os.Open(path)
	if err != nil {
		return Config{}
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read config file", "err:", err)
		return Config{}
	}
	if strings.TrimSpace(string(data)) == "" {
		return Config{}
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error reading config json data. [message] ", err)
	}
	return config
}

func (config *Config) merge(arg Config) *Config {
	mine := reflect.ValueOf(config).Elem()
	theirs := reflect.ValueOf(&arg).Elem()

	for i := 0; i < mine.NumField(); i++ {
		myField := mine.Field(i)
		if misc.ZeroOrNil(myField.Interface()) {
			myField.Set(reflect.ValueOf(theirs.Field(i).Interface()))
		}
	}
	return config
}

func (config *Config) complete() bool {
	cfgs := reflect.ValueOf(config).Elem()

	for i := 0; i < cfgs.NumField(); i++ {
		if misc.ZeroOrNil(cfgs.Field(i).Interface()) {
			return false
		}
	}
	return true
}

func (config *Config) trimWhitespace() {
	cfgs := reflect.ValueOf(config).Elem()
	cfgAttrs := reflect.Indirect(reflect.ValueOf(config)).Type()

	for i := 0; i < cfgs.NumField(); i++ {
		field := cfgs.Field(i)
		if !field.CanInterface() {
			continue
		}
		attr := cfgAttrs.Field(i).Tag.Get("trim")
		if len(attr) == 0 {
			continue
		}
		if field.Kind() != reflect.String {
			continue
		}
		str := field.Interface().(string)
		field.SetString(strings.TrimSpace(str))
	}
}

// String returns a string representation of the config.
func (config *Config) String() string {
	return fmt.Sprintf(
		"Name: %v, Port: %v, LogLevel: %v, AccessLog: %v, "+
			"TwitterKey: %v, TwitterCallback: %v, "+
			"CorsMethods: %v, CorsOrigin: %v, SecuredCookie: %v",
		config.Name, config.Port, config.LogLevel,
		config.AccessLog, config.TwitterKey, config.TwitterCallback,
		config.CorsMethods, config.CorsOrigin, config.SecuredCookie)
}
