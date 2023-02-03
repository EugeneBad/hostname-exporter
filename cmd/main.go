package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

const (
	port string = "9090"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	srv := NewServer(ctx, port)

	log.Printf("Application started successfully! Listening on port %s...", port)

	// With goroutine, launch http server on set port
	go srv.ListenAndServe(ctx)

	// Use interrupt signal channel to block program from exiting immediately
	// Only proceeds when an interrupt signal (Ctrl+C) writes a msg to the channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Application closing...")
	srv.Close(ctx)
	defer cancel()
}
