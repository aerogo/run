package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/rjeczalik/notify"
)

var server *exec.Cmd
var cwd string

func main() {
	cwd, _ = os.Getwd()

	interceptSignals()
	restart()
}

func restart() {
	stopServer()

	var err error

	err = pack()

	if err != nil {
		log.Fatal(err)
		return
	}

	err = build()

	if err != nil {
		log.Fatal(err)
		return
	}

	server = run()
	watch()
}

func watch() {
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch("./...", c, notify.InCloseWrite, notify.InMovedFrom, notify.InMovedTo, notify.Remove); err != nil {
		log.Fatal(err)
	}

	for fileEvent := range c {
		if server == nil {
			continue
		}

		relPath, _ := filepath.Rel(cwd, fileEvent.Path())

		// Ignore hidden files
		if strings.HasPrefix(relPath, ".") {
			continue
		}

		break
	}

	notify.Stop(c)
	restart()
}

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
		server.Process.Kill()
		server = nil
	}
}

func exit() {
	stopServer()
	os.Exit(1)
}

func pack() error {
	cmd := exec.Command("pack")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		color.Red("Couldn't run 'pack'. Make sure you ran 'go install github.com/aerogo/pack'")
		return err
	}

	waitErr := cmd.Wait()

	if waitErr != nil {
		return waitErr
	}

	return nil
}

func build() error {
	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		color.Red("Couldn't run 'go build'. Make sure Go is correctly installed.")
		return err
	}

	waitErr := cmd.Wait()

	if waitErr != nil {
		return waitErr
	}

	return nil
}

func run() *exec.Cmd {
	mainExecutable := filepath.Base(cwd)

	cmd := exec.Command("./" + mainExecutable)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		color.Red("Couldn't start the server.")
		return nil
	}

	return cmd
}
