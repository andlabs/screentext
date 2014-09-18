// 10 june 2014
package main

import (
	"image"
)

// Image represents an in-memory image.
// It satisfies Go's draw.Image.
type Image interface {
	// Close cleans up all resources and renders the image invalid.
	Close()

	// Line draws a line from (x0,y0) to (x1,y1) with the given Pen.
	Line(x0 int, y0 int, x1 int, y1 int, p Pen)

	// Text draws the given string at the given position on the current Image in the given Pen and Font.
	// The top-left corner of the drawn string will be at the given point.
	// TODO pango seems to do this vertically offset?
	Text(text string, x int, y int, f Font, p Pen)

	// Image produces a copy of i as a Go image.RGBA.
	Image() *image.RGBA

	sysImage
}

// NewImage creates a new Image.
// It will initially be fully transparent.
func NewImage(width int, height int) Image {
	return newImage(width, height)
}
