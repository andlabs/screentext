// 21 september 2014

package ndraw

import (
)

// Brush represents a brush.
// Brushes are used to fill shapes, text, etc.
// A Brush is created by passing a BrushSpec to NewBrush().
// As a special rule, a Brush value of nil represents a Brush that draws nothing.
type Brush interface {
	// Close frees resources allocated to the Brush.
	Close()

	sysBrush
}

// BrushSpec represents the properties of a Brush.
type BrushSpec struct {
	R			uint8	// color; alpha-premultiplied
	G			uint8
	B			uint8
	A			uint8
}

// NewBrush creates a Brush from the given BrushSpec.
func NewBrush(spec BrushSpec) Brush {
	return newBrush(spec)
}
