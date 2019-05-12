package main

import (
	"io"

	"github.com/akyoto/color"
)

// coloredWriter writes everything in a single color
type coloredWriter struct {
	io.Writer
	color *color.Color
}

// Write implements io.Writer
func (writer *coloredWriter) Write(b []byte) (int, error) {
	return writer.color.Fprint(writer.Writer, string(b))
}
