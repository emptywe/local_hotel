package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	defAddr         = "127.0.0.1:8080" // Address for local development
	dockerAddr      = "0.0.0.0:8080"   // Address for docker containers
	sessionLifetime = time.Hour * 24
)

var pageStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_page_status_count", // metric name
		Help: "Count of status returned by page.",
	},
	[]string{"page", "status"}, // labels
)

func init() {
	// we need to register the counter so prometheus can collect this metric
	prometheus.MustRegister(pageStatus)
}

func main() {
	fmt.Println("entry point")
}
