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
	db.AutoMigrate(&Booking{})

	handler := NewBookingHandler(db, config.ConcertServiceAddr)

	router := common.NewRouter()
	router.AddRoute("HealthCheck", "GET", "/healthcheck", util.HealthCheck)
	router.AddRoute("AddBooking", "POST", "/bookings", handler.AddBooking)
	router.AddRoute("FindBookings", "GET", "/bookings/{user_id}", handler.FindBookings)

	return http.ListenAndServe(config.ServiceEndpoint, router)
}
