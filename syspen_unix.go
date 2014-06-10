// 10 june 2014
package main

import (
	"fmt"
)

// #include <cairo.h>
import "C"

// cairo doesn't have solid pen objects; instead we set the various parameters individually
type sysPen struct {
	pattern		*C.cairo_pattern_t
	linewidth		uint
	// TODO join
	// TODO cap
	interval		uint
}

func tocairorgb(r uint, g uint, b uint) (C.double, C.double, C.double) {
	xr := C.double(r) / 255
	xg := C.double(g) / 255
	xb := C.double(b) / 255
	return xr, xg, xb
}

func mkSysPenRGB(r uint, g uint, b uint) (s *sysPen) {
	s = new(sysPen)
	s.pattern = C.cairo_pattern_create_rgb(tocairorgb(r, g, b))
	if status := C.cairo_pattern_status(s.pattern); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo pattern for RGB [%d %d %d]: %v", r, g, b, cairoerr(status)))
	}
	return s
}

func (s *sysPen) setLineType(linetype Line, thickness uint) {
	switch linetype {
	case Solid:
		s.interval = 0
	}
	s.linewidth = thickness
}
