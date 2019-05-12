package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func startServer() (*exec.Cmd, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	mainExecutable := filepath.Base(cwd)
	cmd := exec.Command("./" + mainExecutable)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()

	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func stopServer(server *exec.Cmd) error {
	// Server already stopped?
	if server == nil || server.Process == nil {
		return nil
	}

	// Send interrupt signal
	err := server.Process.Signal(os.Interrupt)

	if err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	// Send kill signal
	return server.Process.Kill()
}
