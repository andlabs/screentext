// 18 september 2014

package ndraw

import (
	"sync"
	"image"
	"reflect"
	"unsafe"
)

// #cgo CFLAGS: --std=c99
// #cgo LDFLAGS: -luser32 -lkernel32 -lgdi32
// #include "winapi_windows.h"
import "C"

type sysImage interface {
	// TODO
}

type imagetype struct {
	lock		sync.Mutex
	bitmap	C.HBITMAP
	dc		C.HDC
	prev		C.HBITMAP
	ppvBits	unsafe.Pointer
	width	int		// save these here
	height	int
}

func newImage(width int, height int) Image {
	i := new(imagetype)
	i.bitmap = C.newBitmap(C.int(width), C.int(height), &i.ppvBits)
	i.dc = C.newDCForBitmap(i.bitmap, &i.prev)
	i.width = width
	i.height = height
	return i
}

func (i *imagetype) Close() {
	i.lock.Lock()
	defer i.lock.Unlock()

	C.imageClose(i.bitmap, i.dc, i.prev)
}

// TODO this is [x0y0, x0y1) - the pixel at (x1,y1) is not drawn; check everything
func (i *imagetype) Line(x0 int, y0 int, x1 int, y1 int, p Pen) {
	i.lock.Lock()
	defer i.lock.Unlock()

	prev := p.selectInto(i.dc)
	C.moveTo(i.dc, C.int(x0), C.int(y0))
	C.lineTo(i.dc, C.int(x1), C.int(y1))
	p.unselect(i.dc, prev)
}

// TODO this only supports a single line of text
func (i *imagetype) Text(str string, x int, y int, f Font, p Pen) {
	i.lock.Lock()
	defer i.lock.Unlock()

	prevfont := f.selectInto(i.dc)
	prevpen := p.selectInto(i.dc)
	cstr := C.CString(str)
	defer freestr(cstr)
	C.drawText(i.dc, cstr, C.int(x), C.int(y))
	p.unselect(i.dc, prevpen)
	f.unselect(i.dc, prevfont)
}

// TODO merge with the cairo implementation
func (i *imagetype) Image() (img *image.RGBA) {
	i.lock.Lock()
	defer i.lock.Unlock()

	var ppvBits reflect.SliceHeader

	width := i.width
	height := i.height
	ppvBits.Data = uintptr(i.ppvBits)
	ppvBits.Len = width * height
	ppvBits.Cap = ppvBits.Len
	data := *((*[]uint32)(unsafe.Pointer(&ppvBits)))
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
			// the alpha value is inverted because right now it acts as a written flag: 0xFF means not touched by GDI, 0x00 means touched by GDI
			img.Pix[p + 3] ^= 0xFF
			p += 4
			q++
		}
		p = nextp
		q = nextq
	}
	return img
}