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
	db.AutoMigrate(&User{})

	handler := NewUserHandler(db, config.BookingServiceAddr, config.ConcertServiceAddr)

	router := common.NewRouter()
	router.AddRoute("HealthCheck", "GET", "/healthcheck", util.HealthCheck)
	router.AddRoute("AddUser", "POST", "/users", handler.AddUser)
	router.AddRoute("FindUser", "GET", "/users", handler.FindUsers)
	router.AddRoute("FindUserByID", "GET", "/users/{id}", handler.FindUserByID)
	router.AddRoute("FindUserBookingsByID", "GET", "/users/{user_id}/bookings", handler.FindUserBookingsByID)

	return http.ListenAndServe(config.ServiceEndpoint, router)
}
