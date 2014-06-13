// 10 june 2014
package main

import (
	"fmt"
	"image"
	"reflect"
	"unsafe"
)

// /* TODO really pangocairo? */
// #cgo pkg-config: cairo pango pangocairo
// #include <cairo.h>
import "C"

type sysImage struct {
	cr	*C.cairo_t
	cs	*C.cairo_surface_t
}

func cairoerr(status C.cairo_status_t) string {
	return C.GoString(C.cairo_status_to_string(status))
}

func mkSysImage(width int, height int) (s *sysImage) {
	s = new(sysImage)
	s.cs = C.cairo_image_surface_create(
		C.CAIRO_FORMAT_ARGB32,
		C.int(width), C.int(height))
	if status := C.cairo_surface_status(s.cs); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo surface for image: %v", cairoerr(status)))
	}
	s.cr = C.cairo_create(s.cs)
	if status := C.cairo_status(s.cr); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo context for image: %v", cairoerr(status)))
	}
	return s
}

func (s *sysImage) close() {
	C.cairo_destroy(s.cr)
	C.cairo_surface_destroy(s.cs)
}

func (s *sysImage) selectPen(p *Pen) {
	C.cairo_set_source(s.cr, p.sysPen.pattern)
	C.cairo_set_line_width(s.cr, C.double(p.sysPen.linewidth))
	// TODO join
	// TODO cap
	if p.sysPen.interval == 0 {
		C.cairo_set_dash(s.cr, nil, 0, 0)
	} else {
		interval := C.double(p.sysPen.interval)		// need to take its address
		C.cairo_set_dash(s.cr, &interval, 1, 0)		// 0 = start immediately
	}
}

func (s *sysImage) line(x0 int, y0 int, x1 int, y1 int) {
	C.cairo_new_path(s.cr)
	C.cairo_move_to(s.cr, C.double(x0), C.double(y0))
	C.cairo_line_to(s.cr, C.double(x1), C.double(y1))
	C.cairo_stroke(s.cr)
}

func cairoImageData(cs *C.cairo_surface_t) (data []uint32, stride int) {
	var sh reflect.SliceHeader

	C.cairo_surface_flush(cs)			// perform pending drawing
	height := int(C.cairo_image_surface_get_height(cs))
	stride = int(C.cairo_image_surface_get_stride(cs))
	sh.Data = uintptr(unsafe.Pointer(C.cairo_image_surface_get_data(cs)))
	sh.Len = height * stride			// should be correct for uint32
	sh.Cap = sh.Len
	data = *((*[]uint32)(unsafe.Pointer(&sh)))
	return data, stride
}

func (s *sysImage) toImage() (img *image.RGBA) {
	width := int(C.cairo_image_surface_get_width(s.cs))
	height := int(C.cairo_image_surface_get_height(s.cs))
	data, stride := cairoImageData(s.cs)
	img = image.NewRGBA(image.Rect(0, 0, width, height))
	p := 0
	q := 0
	for y := 0; y < height; y++ {
		nextp := p + img.Stride
		nextq := q + (stride / 4)
		for x := 0; x < width; x++ {
			img.Pix[p] = uint8((data[q] >> 16) & 0xFF)		// R
			img.Pix[p + 1] = uint8((data[q] >> 8) & 0xFF)		// G
			img.Pix[p + 2] = uint8(data[q] & 0xFF)			// B
			img.Pix[p + 3] = uint8((data[q] >> 24) & 0xFF)		// A
			p += 4
			q++
		}
		p = nextp
		q = nextq
	}
	return img
}
