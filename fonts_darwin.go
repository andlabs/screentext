// 20 september 2014

package ndraw

import (
	"unsafe"
)

// #include "coregfx_darwin.h"
import "C"

func listFonts() (fonts []FontSpec) {
	collection := C.CTFontCollectionCreateFromAvailableFonts(nil)
	// TODO check for nil return?
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(collection)))		// TODO correct?

	descs := C.CTFontCollectionCreateMatchingFontDescriptors(collection)
	// TODO check for nil return?
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(descs)))		// TODO correct?

	n := int(C.CFArrayGetCount(descs))
	fonts = make([]FontSpec, n)
	for i := 0; i < n; i++ {
		desc := C.CTFontDescriptorRef(C.CFArrayGetValueAtIndex(descs, C.CFIndex(i)))
		name := C.CTFontDescriptorCopyAttribute(desc, C.kCTFontFamilyNameAttribute)
		if name != nil {
			namestr := C.CFStringRef(unsafe.Pointer(name))
			family := C.CFStringGetCStringPtr(namestr, C.kCFStringEncodingUTF8)
			if family == nil {
				panic("CFStringGetCStringPtr() failed; TODO implement the long way")
			}
			fonts[i].Family = C.GoString(family)
			C.CFRelease(name)
		}
		traits := C.CTFontDescriptorCopyAttribute(desc, C.kCTFontTraitsAttribute)
		if traits != nil {
			traitsd := C.CFDictionaryRef(unsafe.Pointer(traits))
			traitsnum := C.CFNumberRef(C.CFDictionaryGetValue(traitsd, unsafe.Pointer(C.kCTFontSymbolicTrait)))
			if traitsnum != nil {
				// CTFontSymbolicTraits is uint32_t
				var traitsv uint32

				if C.CFNumberGetValue(traitsnum, C.kCFNumberSInt32Type, unsafe.Pointer(&traitsv)) != C.true {
					// TODO get error reason
					panic("error extracting traits value from CFNumber in ListFonts()")
				}
				fonts[i].Bold = (traitsv & C.kCTFontBoldTrait) != 0
				fonts[i].Italic = (traitsv & C.kCTFontItalicTrait) != 0
				fonts[i].Monospace = (traitsv & C.kCTFontMonoSpaceTrait) != 0
				fonts[i].Vertical = (traitsv & C.kCTFontVerticalTrait) != 0
				// do not release; Get rule
			}
			C.CFRelease(traits)
		}
	}

	return fonts
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
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(familyref)))

	basefont := C.CTFontCreateWithName(familyref, C.CGFloat(spec.Size), nil)
	// TODO check for nil?
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(basefont)))	// TODO correct?

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
	C.CFRelease(C.CFTypeRef(unsafe.Pointer(f.f)))
}

func (f *font) toCTLine(text string) C.CTLineRef {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	strref := C.CFStringCreateWithCString(nil, cstr, C.kCFStringEncodingUTF8)
	if strref == nil {
		// TODO get error reason
		panic("error creating CFString for drawing text")
	}
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(strref)))

	attrs := C.CFDictionaryCreateMutable(nil, 1, &C.kCFTypeDictionaryKeyCallBacks, &C.kCFTypeDictionaryValueCallBacks)
	if attrs == nil {
		// TODO get error reason
		panic("error creating text attribute list for drawing text")
	}
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(attrs)))
	C.CFDictionaryAddValue(attrs, unsafe.Pointer(C.kCTFontAttributeName), unsafe.Pointer(f.f))

	attrstr := C.CFAttributedStringCreate(nil, strref, attrs)
	if attrstr == nil {
		// TODO get error reason
		panic("error creating attributed string for drawing text")
	}
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(attrstr)))

	line := C.CTLineCreateWithAttributedString(attrstr)
	if line == nil {
		// TODO get error reason
		panic("error creating CTLine for drawing text")
	}

	return line
}
