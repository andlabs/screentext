// 18 september 2014

package ndraw

// #include "winapi_windows.h"
import "C"

type sysPen interface {
	selectInto(C.HDC) C.HPEN
	unselect(C.HDC, C.HPEN)
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

func (p *pen) Close() {
	C.penClose(p.p)
}

func (p *pen) selectInto(dc C.HDC) C.HPEN {
	return C.penSelectInto(p.p, dc)
}

func (p *pen) unselect(dc C.HDC, prev C.HPEN) {
	C.penUnselect(p.p, dc, prev)
}