// +build !windows,!darwin

// 10 june 2014

package ndraw

import (
	"fmt"
	"sync"
	"image"
	"reflect"
	"unsafe"
)

// /* TODO really pangocairo? */
// #cgo pkg-config: cairo pango pangocairo
// #include <pango/pangocairo.h>
// #include <stdlib.h>
import "C"

type sysImage interface {
	// TODO
}

type imagetype struct {
	lock		sync.Mutex
	cr		*C.cairo_t
	cs		*C.cairo_surface_t
}

func cairoerr(status C.cairo_status_t) string {
	return C.GoString(C.cairo_status_to_string(status))
}

func newImage(width int, height int) Image {
	i := new(imagetype)
	i.cs = C.cairo_image_surface_create(
		C.CAIRO_FORMAT_ARGB32,
		C.int(width), C.int(height))
	if status := C.cairo_surface_status(i.cs); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo surface for image: %v", cairoerr(status)))
	}
	i.cr = C.cairo_create(i.cs)
	if status := C.cairo_status(i.cr); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("error creating cairo context for image: %v", cairoerr(status)))
	}
	return i
}

func (i *imagetype) Close() {
	i.lock.Lock()
	defer i.lock.Unlock()

	C.cairo_destroy(i.cr)
	C.cairo_surface_destroy(i.cs)
}

func (i *imagetype) Line(x0 int, y0 int, x1 int, y1 int, p Pen) {
	i.lock.Lock()
	defer i.lock.Unlock()

	p.selectInto(i.cr)
	C.cairo_new_path(i.cr)
	C.cairo_move_to(i.cr, C.double(x0), C.double(y0))
	C.cairo_line_to(i.cr, C.double(x1), C.double(y1))
	C.cairo_stroke(i.cr)
	deselectPen(i.cr)
}

func (i *imagetype) Text(str string, x int, y int, f Font, p Pen) {
	i.lock.Lock()
	defer i.lock.Unlock()

	p.selectInto(i.cr)
	C.cairo_save(i.cr)
	C.cairo_move_to(i.cr, C.double(x), C.double(y))
	pl := f.selectInto(i.cr)
	cstr := C.CString(str)
	C.pango_layout_set_text(pl, cstr, -1)
	C.free(unsafe.Pointer(cstr))
	C.pango_cairo_show_layout(i.cr, pl)
	C.cairo_restore(i.cr)
	deselectFont(pl)
	deselectPen(i.cr)
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

func (i *imagetype) Image() (img *image.RGBA) {
	i.lock.Lock()
	defer i.lock.Unlock()

	width := int(C.cairo_image_surface_get_width(i.cs))
	height := int(C.cairo_image_surface_get_height(i.cs))
	data, stride := cairoImageData(i.cs)
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
