// 18 september 2014

package ndraw

// #include "winapi_windows.h"
import "C"

type sysBrush interface {
	get() (C.HBRUSH, C.uint8_t)
	// TODO on all platforms, add methods that would ensure a Brush is incompatible with a Pen
}

type brush struct {
	b		C.HBRUSH
	alpha	C.uint8_t
}

func newBrush(spec BrushSpec) Brush {
	var xb C.LOGBRUSH

	b := new(brush)
	xb.lbStyle = C.BS_SOLID
	b.alpha = C.uint8_t(spec.A)
	xb.lbColor = colorref(spec.R, spec.G, spec.B)
	b.b = C.newBrush(&xb)
	return b
}

func (b *brush) Close() {
	C.brushClose(b.b)
}

func (b *brush) get() (C.HBRUSH, C.uint8_t) {
	return b.b, b.alpha
}
