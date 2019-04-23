package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/akyoto/color"
)

var server *exec.Cmd

func restart() {
	pack()
	err := build()
	stopServer()

	if err == nil {
		run()
	}

	watch()
}

func run() {
	mainExecutable := filepath.Base(cwd)

	cmd := exec.Command("./" + mainExecutable)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		color.Red("Couldn't start the server.")
		server = nil
		return
	}

	server = cmd
}
