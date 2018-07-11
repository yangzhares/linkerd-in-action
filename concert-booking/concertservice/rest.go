package main

import (
	"log"
	"net/http"

	"github.com/yangzhares/linkerd-in-action/concert-booking/common"
	"github.com/yangzhares/linkerd-in-action/concert-booking/db"
	"github.com/yangzhares/linkerd-in-action/concert-booking/util"
)

func ServeAPI(config *ServiceConfig) error {
	db, err := db.InitDB(config.DBName, config.DBUser, config.Password, config.DBEndpoint)
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()
	db.AutoMigrate(&Concert{})

	handler := NewConcertHandler(db)

	router := common.NewRouter()
	router.AddRoute("HealthCheck", "GET", "/healthcheck", util.HealthCheck)
	router.AddRoute("AddConcert", "POST", "/concerts", handler.AddConcert)
	router.AddRoute("FindConcertByID", "GET", "/concerts/{id}", handler.FindConcertByID)
	router.AddRoute("FindConcerts", "GET", "/concerts", handler.FindConcerts)

	return http.ListenAndServe(config.ServiceEndpoint, router)
}
