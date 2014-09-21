// 10 june 2014

package ndraw

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

	// Text draws the given string at the given position on the current Image in the given Font.
	// The top-left corner of the drawn string will be at the given point.
	// If the given Pen is not nil, the text is outlined using that Pen.
	// If the given Brush is not nil, the text is filled using that Brush.
	// If you just want to draw text "normally", specify a non-nil Brush of the desired text color and specify a nil Pen.
	// TODO pango seems to do this vertically offset?
	Text(text string, x int, y int, f Font, p Pen, b Brush)

	// TextSize computes the size that the given text string would occupy in the given Font.
	// The reported size is in pixels.
	// It is a method of Image because some systems require a valid graphics drawing context (which Image provides) to make this calculation.
	TextSize(text string, f Font) (width int, height int)

	// Image produces a copy of i as a Go image.RGBA.
	// Note that for technical reasons, the values of the Alpha bytes of the image are undefined; you cannot reasonably blend the result of Image() with something else.
	// TODO alternative
	Image() *image.RGBA

	sysImage
}

// NewImage creates a new Image.
// It will initially be fully transparent.
func NewImage(width int, height int) Image {
	return newImage(width, height)
}
