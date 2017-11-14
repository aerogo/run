package main

import (
	"os"
	"os/signal"
	"syscall"
)

func interceptSignals() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		exit()
	}()
}

func stopServer() {
	if server != nil && server.Process != nil {
		server.Process.Signal(os.Interrupt)
		server = nil
	}
}

func exit() {
	stopServer()
	os.Exit(1)
}
