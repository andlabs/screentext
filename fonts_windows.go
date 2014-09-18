// 18 september 2014

package ndraw

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

func listFonts() []FontSpec {
	fl := make([]FontSpec, 0, 1024)		// just to get things started
	C.listFonts(unsafe.Pointer(&fl))
	return fl
}

//export listFontsAdd
func listFontsAdd(golist unsafe.Pointer, lf *C.LOGFONTW, family *C.char, size C.LONG) {
	fl := (*[]FontSpec)(golist)
	*fl = append(*fl, FontSpec{
		Family:		C.GoString(family),
		Size:			uint(size),
		// TODO this can be levelled... see what we did for cairo
		Bold:			lf.lfWeight == C.FW_BOLD,
		Italic:		lf.lfItalic != C.FALSE,
		Vertical:		false,		// TODO
		Monospace:	(lf.lfPitchAndFamily & 3) == C.FIXED_PITCH,
	})
	freestr(family)
}

type sysFont interface {
	selectInto(C.HDC) C.HFONT
	unselect(C.HDC, C.HFONT)
}

type font struct {
	f	C.HFONT
}

func newFont(fs FontSpec) Font {
	var lf C.LOGFONTW

	f := new(font)
	lf.lfWeight = C.FW_NORMAL
	if fs.Bold {
		lf.lfWeight = C.FW_BOLD
	}
	if fs.Italic {
		lf.lfItalic = C.TRUE
	}
	// TODO vertical
	if fs.Monospace {
		lf.lfPitchAndFamily = C.FIXED_PITCH
	}
	cfamily := C.CString(fs.Family)
	defer freestr(cfamily)
	f.f = C.newFont(&lf, cfamily, C.LONG(fs.Size))
	return f
}

func (f *font) Close() {
	C.fontClose(f.f)
}

func (f *font) selectInto(dc C.HDC) C.HFONT {
	return C.fontSelectInto(f.f, dc)
}

func (f *font) unselect(dc C.HDC, prev C.HFONT) {
	C.fontUnselect(f.f, dc, prev)
}
