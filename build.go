package main

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
)

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
