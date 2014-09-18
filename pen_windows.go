// 18 september 2014

package ndraw

import (
	// ...
)

// #include "winapi_windows.h"
import "C"

type sysPen interface {
	// TODO
}

type pen struct {
	p	C.HPEN
}

var lineTypes = map[Line]C.DWORD{
	Solid:	C.PS_SOLID,
}

func newPen(spec PenSpec) Pen {
	var xp C.struct_xpen

	p := new(pen)
	xp.style = C.PS_GEOMETRIC | lineTypes[spec.Line]
	xp.width = C.DWORD(spec.Thickness)
	xp.brush.lbStyle = C.BS_SOLID
	xp.brush.lbColor = colorref(spec.R, spec.G, spec.B)
	xp.nSegments = 0
	p.p = C.newPen(&xp)
	return p
}
