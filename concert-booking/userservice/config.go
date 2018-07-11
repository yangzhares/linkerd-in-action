package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	DefaultServiceEndpoint    = "0.0.0.0:8180"
	DefaultDBName             = "user"
	DefaultDBUser             = "root"
	DefaultPassword           = "pass"
	DefaultDBEndpoint         = "localhost:3306"
	DefaultBookingServiceAddr = "localhost:8181"
	DefaultConcertServiceAddr = "localhost:8182"
)

type ServiceConfig struct {
	ServiceEndpoint    string `json:"service_endpoint"`
	DBName             string `json:"dbname"`
	DBUser             string `json:"dbuser"`
	Password           string `json:"password"`
	DBEndpoint         string `json:"dbendpoint"`
	BookingServiceAddr string `json:"booking_service_addr"`
	ConcertServiceAddr string `json:"concert_service_addr"`
}

func InitConfig(f string) (*ServiceConfig, error) {
	conf := ServiceConfig{
		ServiceEndpoint:    DefaultServiceEndpoint,
		DBName:             DefaultDBName,
		DBUser:             DefaultDBUser,
		Password:           DefaultPassword,
		DBEndpoint:         DefaultDBEndpoint,
		BookingServiceAddr: DefaultBookingServiceAddr,
		ConcertServiceAddr: DefaultConcertServiceAddr,
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

	if v := os.Getenv("BOOKING_SERVICE_ADDR"); v != "" {
		conf.BookingServiceAddr = v
	}

	if v := os.Getenv("CONCERT_SERVICE_ADDR"); v != "" {
		conf.ConcertServiceAddr = v
	}

	return &conf, nil
}
