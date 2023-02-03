package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var port string
	if val, ok := os.LookupEnv("PORT"); ok {
		port = val
	} else {
		port = "9090"
	}

	srv := NewServer(ctx, port)

	log.Printf("Application started successfully! Listening on port %s...", port)

	go srv.ListenAndServe(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Application closing...")
	srv.Close(ctx)
	defer cancel()
}
