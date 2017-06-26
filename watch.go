package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

var watchExtensions = map[string]bool{
	".go":      true,
	".js":      true,
	".pixy":    true,
	".scarlet": true,
	".json":    true,
}

func watch() {
	c := make(chan notify.EventInfo, 1)
	err := notify.Watch("./...", c, notify.InCloseWrite, notify.InMovedFrom, notify.InMovedTo, notify.Remove)

	if err != nil {
		log.Fatal(err)
	}

	for fileEvent := range c {
		// if server == nil {
		// 	continue
		// }

		relPath, _ := filepath.Rel(cwd, fileEvent.Path())

		// Ignore hidden files
		if strings.HasPrefix(relPath, ".") {
			continue
		}

		extension := filepath.Ext(relPath)

		// Only care about certain extensions
		_, ok := watchExtensions[extension]

		if !ok {
			continue
		}

		break
	}

	notify.Stop(c)
	restart()
}
