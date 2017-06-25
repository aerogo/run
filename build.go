package main

import (
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/fatih/color"
)

var faint = color.New(color.Faint).SprintFunc()
var goBuildMessage = faint(`go build -v`)

func build() error {
	println("--------------------------------------------------------------------------------")
	println(goBuildMessage)

	cmd := exec.Command("go", "build", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Start()

	if err != nil {
		color.Red("Couldn't run 'go build'. Make sure Go is correctly installed.")
		return err
	}

	waitErr := cmd.Wait()
	duration := time.Since(start)
	ms := strconv.Itoa(int(duration.Nanoseconds() / int64(1000000)))

	println()
	println(faint(ms + " ms"))

	if waitErr != nil {
		return waitErr
	}

	return nil
}
