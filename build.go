package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/akyoto/color"
)

var (
	faint          = color.New(color.Faint).SprintFunc()
	goBuildMessage = faint(`go build`)
)

func build() error {
	fmt.Println(goBuildMessage)

	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &coloredWriter{
		Writer: os.Stderr,
		color:  color.New(color.FgRed),
	}
	start := time.Now()
	err := cmd.Start()

	if err != nil {
		return errors.New("Couldn't run 'go build'. Make sure Go is correctly installed.")
	}

	err = cmd.Wait()
	duration := time.Since(start)
	ms := strconv.Itoa(int(duration.Nanoseconds() / int64(1000000)))

	fmt.Println()
	fmt.Println(faint(ms + " ms"))

	if err != nil {
		return err
	}

	return nil
}
