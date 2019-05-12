package main

import (
	"errors"
	"os"
	"os/exec"
)

func pack() error {
	cmd := exec.Command("pack")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		return errors.New("Couldn't run 'pack'. Make sure you ran 'go install github.com/aerogo/pack'")
	}

	return cmd.Wait()
}
