// 20 september 2014

package screentext

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

var (
	lock			sync.Mutex

	// global context for image size calculations
	sizecontext	C.CGContextRef
	sizespace		C.CGColorSpaceRef
)

// screw init() not being run in tests >:(
func testInit() {
	lock.Lock()
	defer lock.Unlock()

	sizespace, sizecontext = newImage(1, 1)
}

func init() {
	testInit()
}

func newImage(width int, height int) (C.CGColorSpaceRef, C.CGContextRef) {
	colorspace := C.CGColorSpaceCreateWithName(C.kCGColorSpaceGenericRGB)
	if colorspace == nil {
		// TODO get error reason
		panic("error creating color space in NewImage()")
	}
	context := C.CGBitmapContextCreate(nil,
		C.size_t(width), C.size_t(height),
		8, 0, colorspace,
		// this matches image.RGBA
		C.kCGImageAlphaPremultipliedLast | C.kCGBitmapByteOrder32Big)
	if context == nil {
		// TODO get error reason
		panic("error creating CGContextRef in NewImage()")
	}
	return colorspace, context
}

func freeImage(colorspace C.CGColorSpaceRef, context C.CGContextRef) {
	C.CGContextRelease(context)
	C.CGColorSpaceRelease(colorspace)
}

// TODO this only supports a single line of text
func line(str string, f Font, r uint8, g uint8, b uint8) *image.RGBA {
	lock.Lock()
	defer lock.Unlock()

	line := f.toCTLine(str)
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(line)))
	width, height := realLineSize(line)
	cspace, context := newImage(width, height)
	defer freeImage(cspace, context)
	// TODO color
	// TODO drawing mode
	C.CGContextSetTextPosition(context, C.CGFloat(0), C.CGFloat(0))
	C.CTLineDraw(line, context)
	return toImage(context)
}

func lineSize(str string, f Font) (int, int) {
	lock.Lock()
	defer lock.Unlock()

	line := f.toCTLine(str)
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(line)))
	return realLineSize(line)
}

func realLineSize(line C.CTLineRef) (int, int) {
	bounds := C.CTLineGetImageBounds(line, sizecontext)
	return int(bounds.size.width), int(bounds.size.height)
}

// assumes lock is held
func toImage(context C.CGContextRef) (img *image.RGBA) {
	var srcs reflect.SliceHeader

	// no need to explicitly flush anything as far as I can see (TODO)
	// there is a CGContextFlush() but it explicitly ignores bitmap contexts
	height := C.CGBitmapContextGetHeight(context)
	stride := C.CGBitmapContextGetBytesPerRow(context)
	srcs.Data = uintptr(C.CGBitmapContextGetData(context))
	srcs.Len = int(stride * height)
	srcs.Cap = srcs.Len
	src := *((*[]uint8)(unsafe.Pointer(&srcs)))
	img = image.NewRGBA(image.Rect(0, 0, int(C.CGBitmapContextGetWidth(context)), int(height)))
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
