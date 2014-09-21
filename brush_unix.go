// +build !windows,!darwin

// 10 june 2014

package ndraw

import (
	"fmt"
)

// #include <cairo.h>
import "C"

// cairo doesn't have solid brush objects; instead we set the various parameters individually
type sysBrush interface {
	selectInto(*C.cairo_t)
}

type brush struct {
	pattern		*C.cairo_pattern_t
}

func newBrush(spec BrushSpec) Brush {
	b := new(brush)
	b.pattern = C.cairo_pattern_create_rgba(tocairorgba(spec.R, spec.G, spec.B, spec.A))
	if status := C.cairo_pattern_status(b.pattern); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo pattern for Brush RGBA [%d %d %d %d]: %v", spec.R, spec.G, spec.B, spec.A, cairoerr(status)))
	}
	return b
}

func (b *brush) Close() {
	C.cairo_pattern_destroy(b.pattern)
}

// assumes the image that owns cr is locked
func (b *brush) selectInto(cr *C.cairo_t) {
	C.cairo_set_source(cr, b.pattern)
}

// assumes the image that owns cr is locked
func deselectBrush(cr *C.cairo_t) {
	C.cairo_set_source_rgb(cr, 0, 0, 0)
}
