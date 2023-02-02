package main

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
}

func NewServer(ctx context.Context, port string) *Server {
	// Return server object that listens on confgured port
	return &Server{
		server: &http.Server{Addr: fmt.Sprintf(":%s", port)},
	}
}

func (srv *Server) ListenAndServe(ctx context.Context) {
	http.HandleFunc("/hostname", srv.ServerHandler)
	srv.server.ListenAndServe()
}

func (srv *Server) ServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)

	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) Close(ctx context.Context) {
	log.Info("Graceful server shutdown...")
	if err := srv.server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful server shutdown failed:%+s", err)
	}
}
