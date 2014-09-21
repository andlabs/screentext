// 20 september 2014

package ndraw

import (
	"sync"
	"image"
	"reflect"
	"unsafe"
)

// #cgo CFLAGS: -mmacosx-version-min=10.6 -DMACOSX_DEPLOYMENT_TARGET=10.6
// #cgo LDFLAGS: -mmacosx-version-min=10.6 -framework ApplicationServices
// #include "coregfx_darwin.h"
import "C"

// TODO use layers?

type sysImage interface {
	// TODO
}

type imagetype struct {
	lock			sync.Mutex
	context		C.CGContextRef
	colorspace	C.CGColorSpaceRef
}

func newImage(width int, height int) Image {
	i := new(imagetype)
	i.colorspace = C.CGColorSpaceCreateWithName(C.kCGColorSpaceGenericRGB)
	if i.colorspace == nil {
		// TODO get error reason
		panic("error creating color space in NewImage()")
	}
	i.context = C.CGBitmapContextCreate(nil,
		C.size_t(width), C.size_t(height),
		8, 0, i.colorspace,
		// this matches image.RGBA
		C.kCGImageAlphaPremultipliedLast | C.kCGBitmapByteOrder32Big)
	if i.context == nil {
		// TODO get error reason
		panic("error creating CGContextRef in NewImage()")
	}
	// now we want the context's origin to be the upper-left
	// see https://wiki.mozilla.org/NPAPI:CoreGraphicsDrawing
	C.CGContextTranslateCTM(i.context, 0.0, C.CGFloat(height))
	C.CGContextScaleCTM(i.context, 1.0, -1.0)
	return i
}

func (i *imagetype) Close() {
	C.CGContextRelease(i.context)
	C.CGColorSpaceRelease(i.colorspace)
}

func (i *imagetype) Line(x0 int, y0 int, x1 int, y1 int, p Pen) {
	i.lock.Lock()
	defer i.lock.Unlock()

	p.selectInto(i.context)
	C.CGContextBeginPath(i.context)
	C.CGContextMoveToPoint(i.context, C.CGFloat(x0), C.CGFloat(y0))
	C.CGContextAddLineToPoint(i.context, C.CGFloat(x1), C.CGFloat(y1))
	C.CGContextStrokePath(i.context)
}

// TODO fills blah blah blah
func (i *imagetype) Text(str string, x int, y int, f Font, p Pen) {
	i.lock.Lock()
	defer i.lock.Unlock()

	p.selectInto(i.context)
	line := f.toCTLine(str)
	C.CGContextSetTextDrawingMode(i.context, C.kCGTextStroke)
	C.CGContextSetTextPosition(i.context, C.CGFloat(x), C.CGFloat(y))
	C.CTLineDraw(line, i.context)
	C.CFRelease(C.CFTypeRef(unsafe.Pointer(line)))
}

func (i *imagetype) Image() (img *image.RGBA) {
	i.lock.Lock()
	defer i.lock.Unlock()

	var srcs reflect.SliceHeader

	// no need to explicitly flush anything as far as I can see (TODO)
	// there is a CGContextFlush() but it explicitly ignores bitmap contexts
	height := C.CGBitmapContextGetHeight(i.context)
	stride := C.CGBitmapContextGetBytesPerRow(i.context)
	srcs.Data = uintptr(C.CGBitmapContextGetData(i.context))
	srcs.Len = int(stride * height)
	srcs.Cap = srcs.Len
	src := *((*[]uint8)(unsafe.Pointer(&srcs)))
	img = image.NewRGBA(image.Rect(0, 0, int(C.CGBitmapContextGetWidth(i.context)), int(height)))
	p := 0
	q := 0
	for y := 0; y < int(height); y++ {
		nextp := p + img.Stride
		nextq := q + int(stride)
		copy(img.Pix[p:nextp], src[q:nextq])
		p = nextp
		q = nextq
	}
	return img
}
