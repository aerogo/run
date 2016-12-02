package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

func watch() {
	c := make(chan notify.EventInfo, 1)
	err := notify.Watch("./...", c, notify.InCloseWrite, notify.InMovedFrom, notify.InMovedTo, notify.Remove)

	if err != nil {
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
