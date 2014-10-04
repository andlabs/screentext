// 18 september 2014

package screentext

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

//export ndraw_xpanic
func ndraw_xpanic(msg *C.char, lasterr C.DWORD) {
	panic(fmt.Errorf("%s: %s", C.GoString(msg), syscall.Errno(lasterr)))
}

func freestr(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func colorref(r uint8, g uint8, b uint8) C.COLORREF {
	return C.colorref(C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
}
