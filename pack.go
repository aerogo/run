package main

import (
	"os"
	"os/exec"

	"github.com/blitzprog/color"
)

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
