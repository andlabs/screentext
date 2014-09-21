// 10 june 2014

package ndraw

// Pen represents a pen.
// Pens are used to draw lines, shape outlines, text, etc.
// A Pen is created by passing a PenSpec to NewPen().
type Pen interface {
	// Close frees resources allocated to the Pen.
	Close()

	sysPen
}

// PenSpec represents the properties of a Pen.
type PenSpec struct {
	R			uint8	// color; alpha-premultiplied
	G			uint8
	B			uint8
	A			uint8
	Thickness		uint		// in pixels
	Line			Line
}

// Line represents a style of line for a Pen.
type Line uint
const (
	Solid Line = iota
)

// NewPen creates a Pen from the given PenSpec.
func NewPen(spec PenSpec) Pen {
	return newPen(spec)
}
