package main

import (
	"io"

	"github.com/akyoto/color"
)

// ColoredWriter writes everything in a single color
type ColoredWriter struct {
	io.Writer
	color *color.Color
}

// Write implements io.Writer
func (writer *ColoredWriter) Write(b []byte) (int, error) {
	return writer.color.Fprint(writer.Writer, string(b))
}
