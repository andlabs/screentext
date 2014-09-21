// 20 september 2014

package ndraw

// #include "coregfx_darwin.h"
import "C"

// TODO there are solid color spaces and pattern color spaces :S figure out what these are and how to switch, otherwise pens can only be solid colors

// Core Graphics doesn't have solid pen objects; instead we set the various parameters individually
// any CFType objects in the pen data get retained, so we don't have to deselect anything
type sysPen interface {
	selectInto(C.CGContextRef)
}

type pen struct {
	color	C.CGColorRef
	width	C.CGFloat
}

func toquartzrgba(r uint8, g uint8, b uint8, a uint8) (C.CGFloat, C.CGFloat, C.CGFloat, C.CGFloat) {
	xr := C.CGFloat(r) / 255
	xg := C.CGFloat(g) / 255
	xb := C.CGFloat(b) / 255
	xa := C.CGFloat(a) / 255
	return xr, xg, xb, xa
}

func newPen(spec PenSpec) Pen {
	p := new(pen)
	p.color = C.CGColorCreateGenericRGB(toquartzrgba(spec.R, spec.G, spec.B, spec.A))
	// TODO check nil return?
	p.width = C.CGFloat(spec.Thickness)
	// TODO line dashing
	return p
}

func (p *pen) Close() {
	C.CGColorRelease(p.color)
}

func (p *pen) selectInto(context C.CGContextRef) {
	C.CGContextSetStrokeColorWithColor(context, p.color)
	C.CGContextSetLineWidth(context, p.width)
}
