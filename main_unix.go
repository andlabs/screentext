// +build !windows,!darwin

// 10 june 2014

package screentext

import (
	"fmt"
	"sync"
	"image"
	"reflect"
	"unsafe"
)

// TODO clean up
// TODO leaves space above and below actual line

// /* TODO really pangocairo? */
// #cgo pkg-config: cairo pango pangocairo
// #include <pango/pangocairo.h>
// #include <stdlib.h>
import "C"

var (
	lock		sync.Mutex

	// global context for image size calculations
	// TODO rename
	cr		*C.cairo_t
	cs		*C.cairo_surface_t
)

// screw init() not being run in tests >:(
func testInit() {
	lock.Lock()
	defer lock.Unlock()

	cs, cr = newContext(1, 1)
}

func init() {
	testInit()
}

func cairoerr(status C.cairo_status_t) string {
	return C.GoString(C.cairo_status_to_string(status))
}

func newContext(width int, height int) (*C.cairo_surface_t, *C.cairo_t) {
	cs := C.cairo_image_surface_create(
		C.CAIRO_FORMAT_ARGB32,
		C.int(width), C.int(height))
	if status := C.cairo_surface_status(cs); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating new cairo surface: %v", cairoerr(status)))
	}
	cr := C.cairo_create(cs)
	if status := C.cairo_status(cr); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating new cairo context: %v", cairoerr(status)))
	}
	return cs, cr
}

func freeContext(cs *C.cairo_surface_t, cr *C.cairo_t) {
	C.cairo_destroy(cr)
	C.cairo_surface_destroy(cs)
}

// TODO this only supports a single line of text
func line(str string, f Font, r uint8, g uint8, b uint8) *image.RGBA {
	lock.Lock()
	defer lock.Unlock()

	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	width, height := realLineSize(cstr, f)
	cs, cr := newContext(width, height)
	defer freeContext(cs, cr)
	pl := f.selectInto(cr)
	C.cairo_save(cr)
	C.cairo_move_to(cr, 0, 0)		// TODO adjust for subpixel rendering?
	C.pango_layout_set_text(pl, cstr, -1)
	rr := C.double(r) / 255
	rg := C.double(g) / 255
	rb := C.double(b) / 255
	C.cairo_set_source_rgb(cr, rr, rg, rb)
	C.pango_cairo_show_layout(cr, pl)
	C.cairo_restore(cr)
	deselectFont(pl)
	return toImage(cs)
}

func lineSize(str string, f Font) (int, int) {
	lock.Lock()
	defer lock.Unlock()

	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	return realLineSize(cstr, f)
}

// assumes lock is held
func realLineSize(cstr *C.char, f Font) (int, int) {
	var width, height C.int

	pl := f.selectInto(cr)
	C.cairo_save(cr)
	C.pango_layout_set_text(pl, cstr, -1)
	C.pango_layout_get_pixel_size(pl, &width, &height)
	C.cairo_restore(cr)
	deselectFont(pl)
	return int(width), int(height)
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

// assumes lock is held
func toImage(cs *C.cairo_surface_t) (img *image.RGBA) {
	width := int(C.cairo_image_surface_get_width(cs))
	height := int(C.cairo_image_surface_get_height(cs))
	data, stride := cairoImageData(cs)
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
