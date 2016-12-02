package main

import "os"

var cwd string

func main() {
	cwd, _ = os.Getwd()

	interceptSignals()
	restart()
}
