// 18 september 2014

package ndraw

import (
	"sync"
	"image"
	"reflect"
	"unsafe"
)

// #cgo CFLAGS: --std=c99
// #cgo LDFLAGS: -luser32 -lkernel32 -lgdi32 -lmsimg32
// #include "winapi_windows.h"
import "C"

var lock sync.Mutex

func init() {
	lock.Lock()
	defer lock.Unlock()

	C.init()
}

// TODO this only supports a single line of text
func line(str string, f Font, r uint8, g uint8, b uint8) *image.RGBA {
	lock.Lock()
	defer lock.Unlock()

	font := f.get()
	cstr := C.CString(str)
	defer freestr(cstr)
	i := C.drawText(cstr, font, C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
	defer C.imageClose(i)
	return toImage(i)
}

func lineSize(str string, f Font) (int, int) {
	lock.Lock()
	defer lock.Unlock()

	font := f.get()
	cstr := C.CString(str)
	defer freestr(cstr)
	size := C.textSize(cstr, font)
	return int(size.cx), int(size.cy)
}

// assumes lock is held
// TODO merge with the cairo implementation
func toImage(i *C.struct_image) (img *image.RGBA) {
	var s reflect.SliceHeader

	width := int(i.width)
	height := int(i.height)
	s.Data = uintptr(unsafe.Pointer(i.ppvBits))
	s.Len = width * height
	s.Cap = s.Len
	data := *((*[]uint32)(unsafe.Pointer(&s)))
	stride := width * 4
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
			// img.Pix[p + 3] is either 0x00 (written by GDI) or 0xFF (not written  by GDI)
			// (this is why we fill ppvBits with 0xFF in image_windows.c newImage())
			// but alpha is the opposite, so we invert
			img.Pix[p + 3] ^= 0xFF
			p += 4
			q++
		}
		p = nextp
		q = nextq
	}
	return img
}
