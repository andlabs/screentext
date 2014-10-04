// 18 september 2014

package screentext

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

// screw init() not being run in tests >:(
func testInit() {
	lock.Lock()
	defer lock.Unlock()

	C.init()
}

func init() {
	testInit()
}

// TODO this only supports a single line of text
func line(str string, f Font, r uint8, g uint8, b uint8) *image.RGBA {
	lock.Lock()
	defer lock.Unlock()

	font := f.get()
	cstr := C.CString(str)
	defer freestr(cstr)
	i := C.drawText(cstr, font)
	defer C.imageClose(i)
	return toImage(i, r, g, b)
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

// same premultiplication that GDI AlphaBlend() says to do; see http://msdn.microsoft.com/en-us/library/dd183393%28v=vs.85%29.aspx
func premultiply(c uint8, alpha uint8) uint8 {
	part := int(c) * int(alpha)
	return uint8(part / 0xFF)
}

// assumes lock is held
func toImage(i *C.struct_image, r uint8, g uint8, b uint8) (img *image.RGBA) {
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
			// GDI doesn't natively support antialiasing text to transparent
			// but here's a clever trick I wish I had thought of:
			// we set the image background to white in newImage() in image_windows.c
			// we set the text color to black in drawText() in image_windows.c
			// this means the pixel color can be used as an alpha, with white being fully transparent and black being fully opaque
			// (all three color components sould be equal by definition)
			// we then manually alpha-premultiply the color
			// full credit for this goes to arx at http://stackoverflow.com/a/26025936/3408572
			// TODO http://stackoverflow.com/questions/26023798/is-it-possible-to-render-antialiased-text-onto-a-transparent-background-with-pur#comment40943112_26025936
			alpha := uint8((data[q] >> 16) & 0xFF)			// use red component
			// white is 0xFF and black is 0x00; we need the opposite
			alpha = 255 - alpha
			img.Pix[p] = premultiply(r, alpha)			// R
			img.Pix[p + 1] = premultiply(g, alpha)		// G
			img.Pix[p + 2] = premultiply(b, alpha)		// B
			img.Pix[p + 3] = alpha					// A
			p += 4
			q++
		}
		p = nextp
		q = nextq
	}
	return img
}
