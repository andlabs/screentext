// +build !windows,!darwin

// 10 june 2014

package ndraw

import (
	"fmt"
)

// #include <cairo.h>
import "C"

// cairo doesn't have solid pen objects; instead we set the various parameters individually
type sysPen interface {
	selectInto(*C.cairo_t)
	thickness() uint
}

type pen struct {
	pattern		*C.cairo_pattern_t
	linewidth		uint
	// TODO join
	// TODO cap
	interval		uint
}

// TODO split into common_unix.go
func tocairorgba(r uint8, g uint8, b uint8, a uint8) (C.double, C.double, C.double, C.double) {
	xr := C.double(r) / 255
	xg := C.double(g) / 255
	xb := C.double(b) / 255
	xa := C.double(a) / 255
	return xr, xg, xb, xa
}

func newPen(spec PenSpec) Pen {
	p := new(pen)
	p.pattern = C.cairo_pattern_create_rgba(tocairorgba(spec.R, spec.G, spec.B, spec.A))
	if status := C.cairo_pattern_status(p.pattern); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo pattern for Pen RGBA [%d %d %d %d]: %v", spec.R, spec.G, spec.B, spec.A, cairoerr(status)))
	}
	p.linewidth = spec.Thickness
	switch spec.Line {
	case Solid:
		p.interval = 0
	}
	return p
}

func (p *pen) Close() {
	C.cairo_pattern_destroy(p.pattern)
}

// assumes the image that owns cr is locked
func (p *pen) selectInto(cr *C.cairo_t) {
	C.cairo_set_source(cr, p.pattern)
	C.cairo_set_line_width(cr, C.double(p.linewidth))
	// TODO join
	// TODO cap
	if p.interval == 0 {
		C.cairo_set_dash(cr, nil, 0, 0)
	} else {
		interval := C.double(p.interval)			// need to take its address
		C.cairo_set_dash(cr, &interval, 1, 0)		// 0 = start immediately
	}
}

// assumes the image that owns cr is locked
func deselectPen(cr *C.cairo_t) {
	C.cairo_set_source_rgb(cr, 0, 0, 0)
}

func (p *pen) thickness() uint {
	return p.linewidth
}
