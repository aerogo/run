package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		err := server.Process.Signal(os.Interrupt)

		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(50 * time.Millisecond)

		err = server.Process.Kill()

		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(50 * time.Millisecond)
		server = nil
	}
}

func exit() {
	stopServer()
	os.Exit(1)
}
