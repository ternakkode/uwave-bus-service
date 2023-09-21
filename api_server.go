package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartAPIServer() {
	r := gin.New()

	srv := &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	r.GET("bus-stops", GetBusStops)
	r.GET("bus-stops/:id", GetBusStopsByID)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start server : %v", err)
	}
}
