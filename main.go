// 10 june 2014

package ndraw

import (
	"image"
)

// Line draws the given line of text at the given position on the current Image in the given Font.
// The color is specified as an RGB value with no alpha information.
// TODO pango seems to do this vertically offset?
func Line(text string, f Font, r uint8, g uint8, b uint8) *image.RGBA {
	return line(text, f, r, g, b)
}

// LineSize computes the size that the given line of text would occupy in the given Font.
// The reported size is in pixels.
// It is a method of Image because some systems require a valid graphics drawing context (which Image provides) to make this calculation.
func LineSize(text string, f Font) (width int, height int) {
	return lineSize(text, f)
}
