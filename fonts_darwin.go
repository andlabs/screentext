// 20 september 2014

package ndraw

import (
	"unsafe"
)

// TODO
import "C"

func listFonts() (fonts []FontSpec) {
	// TODO
	return nil
}

type sysFont interface {
	toCTLine(text string) C.CTLineRef
}

type font struct {
	f	C.CTFontRef
}

func newFont(spec FontSpec) Font {
	f := new(font)

	cfamily := C.CString(spec.Family)
	defer C.free(unsafe.Pointer(cfamily))
	familyref := C.CFStringCreateWithCString(nil, cfamily, C.kCFStringEncodingUTF8)
	if familyref == nil {
		// TODO get error reason
		panic("error creating CFString for NewFont() family name")
	}
	defer C.CFRelease(familyref)

	// TODO other fields

	desc := C.CTFontDescriptorCreateWithAttributes(attrs)
	// TODO check for nil?
	defer C.CFRelease(desc)		// TODO correct?

	basefont := C.CTFontCreateWithName(desc, C.CGFloat(spec.Size), nil)
	// TODO check for nil?
	defer C.CFRelease(basefont)	// TODO correct?

	traits := C.CTFontSymbolicTraits(0)
	if spec.Bold {
		traits |= C.kCTFontBoldTrait
	}
	if spec.Italic {
		traits |= C.kCTFontItalicTrait
	}
	// TODO specify these two? they ARE used by OS X since a font can have both variants in one set
	// TODO if so, make sure these are set in the other ports too
	if spec.Monospace {
		traits |= C.kCTFontMonoSpaceTrait
	}
	if spec.Vertical {
		traits |= C.kCTFontVerticalTrait
	}

	// 0.0 preserves original size; the second traits is the bit mask
	f.f = C.CTFontCreateCopyWithSymbolicTraits(basefont, 0.0, nil, traits, traits)
	if f.f == nil {
		// TODO get reason
		panic("error creating attributed CTFont in NewFont()")
	}
	// TODO fast-track if traits == 0?

	return f
}

func (f *font) Close() {
	// TODO is this correct?
	C.CFRelease(f.f)
}

func (f *font) toCTLine(text string) C.CTLineRef {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	strref := C.CFStringCreateWithCString(nil, cstr, C.kCFStringEncodingUTF8)
	if strref == nil {
		// TODO get error reason
		panic("error creating CFString for drawing text")
	}
	defer C.CFRelease(strref)

	attrs := C.CFDictionaryCreateMutable(nil, 1, &C.kCFTypeDictionaryKeyCallBacks, &C.kCFTypeDictionaryValueCallBacks)
	if attrs == nil {
		// TODO get error reason
		panic("error creating text attribute list for drawing text")
	}
	defer C.CFRelease(attrs)
	C.CFDictionaryAddValue(attrs, unsafe.Pointer(C.kCTFontAttributeName), unsafe.Pointer(f.f))

	attrstr := C.CFAttributedStringCreate(nil, strref, attrs)
	if attrstr == nil {
		// TODO get error reason
		panic("error creating attributed string for drawing text")
	}
	defer C.CFRelease(attstr)

	line := C.CTLineCreateWithAttributedString(attrstr)
	if line == nil {
		// TODO get error reason
		panic("error creating CTLine for drawing text")
	}

	return line
}
