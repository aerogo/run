package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/akyoto/color"
	"github.com/rjeczalik/notify"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	// Notify us about kill signals
	terminator := make(chan os.Signal, 1)
	signal.Notify(terminator, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-terminator
		// stopServer()
		os.Exit(1)
	}()

	// Get working directory
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	// Initial pack
	packError := pack()

	if packError != nil {
		color.Red(packError.Error())
	}

	fmt.Println()

	// Initial build
	var buildError error

	if packError == nil {
		buildError = build()

		if buildError != nil {
			color.Red(buildError.Error())
		}
	}

	// Start server
	var server *exec.Cmd

	if packError == nil && buildError == nil {
		server, err = startServer()

		if err != nil {
			color.Red(err.Error())
		}
	}

	// Create a channel for file system events
	eventChannel := make(chan notify.EventInfo, 1)
	err = notify.Watch("./...", eventChannel, notify.Write, notify.Rename, notify.Remove)
	batchingTime := 100 * time.Millisecond
	acceptNewEvents := time.Now()

	if err != nil {
		return err
	}

	defer close(eventChannel)
	defer notify.Stop(eventChannel)

	for fileEvent := range eventChannel {
		if time.Now().Before(acceptNewEvents) {
			continue
		}

		relPath, err := filepath.Rel(cwd, fileEvent.Path())

		if err != nil {
			return err
		}

		// Ignore hidden files
		if strings.HasPrefix(relPath, ".") {
			continue
		}

		// Get file extension
		extension := filepath.Ext(relPath)

		// Specify rebuild function
		var rebuild func() error

		switch extension {
		case ".pixy", ".scarlet", ".js":
			// When template files change, run the asset packer.
			rebuild = func() error {
				err := pack()

				if err != nil {
					return err
				}

				fmt.Println()
				return build()
			}

		case ".go", ".mod":
			// When Go files change, rebuild the executable.
			rebuild = build

		case ".json":
			// Just restart the server.
			rebuild = func() error {
				return nil
			}

		default:
			continue
		}

		// Log
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("%s %s\n", faint("Detected change in"), color.YellowString(relPath))
		fmt.Println("--------------------------------------------------------------------------------")

		// Don't accept new events for some time.
		stopEvents := func(duration time.Duration) {
			acceptNewEvents = time.Now().Add(duration)
		}

		// Execute rebuild function
		err = rebuild()

		if err != nil {
			color.Red(err.Error())
			stopEvents(batchingTime)
			continue
		}

		// Stop old server
		err = stopServer(server)

		if err != nil {
			color.Red(err.Error())
			stopEvents(batchingTime)
			continue
		}

		fmt.Println("--------------------------------------------------------------------------------")
		server, err = startServer()

		if err != nil {
			color.Red(err.Error())
		}

		stopEvents(batchingTime)
	}

	return nil
}
