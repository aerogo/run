package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

var server *exec.Cmd

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
