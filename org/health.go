package org

import (
	"net/http"
	"time"

	"github.com/cristianuser/go-bun-webserver/httputil"
	"github.com/uptrace/bunrouter"
)

type HealthCheckResponse struct {
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	CurrentTime time.Time `json:"current_time"`
	Version     string    `json:"version"`
}

func HealthCheckHandler(w http.ResponseWriter, req bunrouter.Request) error {

	response := HealthCheckResponse{
		Status:      "ok",
		Message:     "Service is healthy",
		CurrentTime: time.Now(),
		Version:     "1.0.0",
	}
	return httputil.JSON(w, response, http.StatusOK)
}
