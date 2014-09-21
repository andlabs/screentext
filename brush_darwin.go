// 20 september 2014

package ndraw

// #include "coregfx_darwin.h"
import "C"

// TODO there are solid color spaces and pattern color spaces :S figure out what these are and how to switch, otherwise pens can only be solid colors

// Core Graphics doesn't have solid brush objects; instead we set the various parameters individually
// any CFType objects in the brush data get retained, so we don't have to deselect anything
type sysBrush interface {
	selectInto(C.CGContextRef)
}

type brush struct {
	color	C.CGColorRef
}

func newBrush(spec BrushSpec) Brush {
	b := new(brush)
	b.color = C.CGColorCreateGenericRGB(toquartzrgba(spec.R, spec.G, spec.B, spec.A))
	// TODO check nil return?
	return b
}

func (b *brush) Close() {
	C.CGColorRelease(b.color)
}

func (b *brush) selectInto(context C.CGContextRef) {
	// TODO doesn't work?
	C.CGContextSetFillColorWithColor(context, b.color)
}
