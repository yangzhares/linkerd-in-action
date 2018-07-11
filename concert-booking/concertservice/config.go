package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	DefaultServiceEndpoint = "0.0.0.0:8182"
	DefaultDBName          = "concert"
	DefaultDBUser          = "root"
	DefaultPassword        = "pass"
	DefaultDBEndpoint      = "localhost:3306"
)

type ServiceConfig struct {
	ServiceEndpoint string `json:"service_endpoint"`
	DBName          string `json:"dbname"`
	DBUser          string `json:"dbuser"`
	Password        string `json:"password"`
	DBEndpoint      string `json:"dbendpoint"`
}

func InitConfig(f string) (*ServiceConfig, error) {
	conf := ServiceConfig{
		ServiceEndpoint: DefaultServiceEndpoint,
		DBName:          DefaultDBName,
		DBUser:          DefaultDBUser,
		Password:        DefaultPassword,
		DBEndpoint:      DefaultDBEndpoint,
	}

	content, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("Failed to read config file: %v, use default configurations!", err)
	}

	if err := json.Unmarshal(content, &conf); err != nil {
		log.Printf("Failed to init configuration: %v, use default configurations!", err)
	}

	if v := os.Getenv("SERVICE_ENDPOINT"); v != "" {
		conf.ServiceEndpoint = v
	}

	if v := os.Getenv("DBNAME"); v != "" {
		conf.DBName = v
	}

	if v := os.Getenv("DBUSER"); v != "" {
		conf.DBUser = v
	}

	if v := os.Getenv("PASSWORD"); v != "" {
		conf.Password = v
	}

	if v := os.Getenv("DBENDPOINT"); v != "" {
		conf.DBEndpoint = v
	}

	return &conf, nil
}
