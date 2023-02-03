package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
}

var (
	// Initialise the prometheus counter metric to count incoming requests
	reqCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_request_count_total",
			Help: "The total number of requests processed",
		},
		[]string{"code", "method"},
	)
)

// Return server object that listens on confgured port

func NewServer(ctx context.Context, port string) *Server {
	return &Server{
		server: &http.Server{Addr: fmt.Sprintf(":%s", port)},
	}
}

func (srv *Server) ListenAndServe(ctx context.Context) {
	http.HandleFunc("/hostname", srv.ServerHandler)
	http.Handle("/metrics", promhttp.Handler())

	srv.server.ListenAndServe()
}

// Handler function for GET requests on /hostname
func (srv *Server) ServerHandler(w http.ResponseWriter, r *http.Request) {
	// If non-GET request i received, return a 400 BadRequest response code
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		log.WithFields(
			log.Fields{
				"status_code": http.StatusBadRequest,
				"user_agent":  r.UserAgent(),
			},
		).Error("Request should be GET")
		reqCount.WithLabelValues(fmt.Sprint(http.StatusBadRequest), r.Method).Inc()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Implcitly writes a 200 OK response ie w.WriteHeader(http.StatusOK)
	// Returns the json object
	json.NewEncoder(w).Encode(struct {
		Timestamp string `json:"timestamp"`
		Hostname  string `json:"hostname"`
	}{
		Timestamp: time.Now().Format(time.RFC3339),
		Hostname:  os.Getenv("HOSTNAME"),
	})
	reqCount.WithLabelValues(fmt.Sprint(http.StatusOK), r.Method).Inc()
	log.WithFields(
		log.Fields{
			"status_code": http.StatusOK,
			"user_agent":  r.UserAgent(),
		},
	).Infof("Successfully returned hostname: %s", os.Getenv("HOSTNAME"))
}

// Gracefully shutsdown server by closing the main context.
func (srv *Server) Close(ctx context.Context) {
	log.Info("Graceful server shutdown!")
	if err := srv.server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful server shutdown failed:%+s", err)
	}
}
