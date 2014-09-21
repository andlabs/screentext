// 18 september 2014

package ndraw

// #include "winapi_windows.h"
import "C"

// TODO un-premultiply the alpha here and in Brush (or do it in imageInternalBlend())

type sysPen interface {
	get() (C.HPEN, C.uint8_t)
}

type pen struct {
	p		C.HPEN
	alpha	C.uint8_t
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
	p.alpha = C.uint8_t(spec.A)
	xp.brush.lbColor = colorref(spec.R, spec.G, spec.B)
	xp.nSegments = 0
	p.p = C.newPen(&xp)
	return p
}

func (p *pen) Close() {
	C.penClose(p.p)
}

func (p *pen) get() (C.HPEN, C.uint8_t) {
	return p.p, p.alpha
}
