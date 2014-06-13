// 13 june 2014
package main

import (
	"unsafe"
	"reflect"
)

// #include <pango/pangocairo.h>
// #include <stdlib.h>
import "C"

func sysListFonts() (fonts []Font) {
	fonts = make([]Font, 0, 1024)		// initial cap to avoid lots of small allocations

	var xfamilies **C.PangoFontFamily
	var nfamilies C.int
	var sh reflect.SliceHeader

	fontmap := C.pango_cairo_font_map_get_default()
	C.pango_font_map_list_families(fontmap, &xfamilies, &nfamilies)
	sh.Data = uintptr(unsafe.Pointer(xfamilies))
	sh.Len = int(nfamilies)
	sh.Cap = sh.Cap
	families := *(*[]*C.PangoFontFamily)(unsafe.Pointer(&sh))

	for _, family := range families {
		var xfaces **C.PangoFontFace
		var nfaces C.int
		var sh reflect.SliceHeader

		name := C.GoString(C.pango_font_family_get_name(family))
		monospace := C.pango_font_family_is_monospace(family) != C.FALSE
		C.pango_font_family_list_faces(family, &xfaces, &nfaces)
		sh.Data = uintptr(unsafe.Pointer(xfaces))
		sh.Len = int(nfaces)
		sh.Cap = sh.Len
		faces := *(*[]*C.PangoFontFace)(unsafe.Pointer(&sh))

		for _, face := range faces {
			f := Font{
				Family:		name,
				Monospace:	monospace,
			}
			desc := C.pango_font_face_describe(face)
			set := C.pango_font_description_get_set_fields(desc)
			if set & C.PANGO_FONT_MASK_STYLE != 0 {
				style := C.pango_font_description_get_style(desc)
				f.Italic = (style == C.PANGO_STYLE_ITALIC) ||
					(style == C.PANGO_STYLE_OBLIQUE)
			}
			if set & C.PANGO_FONT_MASK_WEIGHT != 0 {
				weight := C.pango_font_description_get_weight(desc)
				f.Bold = (weight >= C.PANGO_WEIGHT_BOLD)		// TODO
			}
			if set & C.PANGO_FONT_MASK_GRAVITY != 0 {
				gravity := C.pango_font_description_get_gravity(desc)
				f.Vertical = (gravity == C.PANGO_GRAVITY_EAST)		// TODO
			}
			fonts = append(fonts, f)
			C.pango_font_description_free(desc)
		}

		C.g_free(C.gpointer(unsafe.Pointer(xfaces)))
	}

	C.g_free(C.gpointer(unsafe.Pointer(xfamilies)))
	return fonts
}

func (s *sysImage) setFont(f Font) {
	if s.pl != nil {		// free old PangoLayout
		C.g_object_unref(C.gpointer(unsafe.Pointer(s.pl)))
	}
	s.pl = C.pango_cairo_create_layout(s.cr)
	desc := C.pango_font_description_new()
	cfamily := C.CString(f.Family)
	C.pango_font_description_set_family(desc, cfamily)
	C.free(unsafe.Pointer(cfamily))
	C.pango_font_description_set_size(desc, C.gint(f.Size * C.PANGO_SCALE))
	if f.Bold {
		C.pango_font_description_set_weight(desc, C.PANGO_WEIGHT_BOLD)
	}
	if f.Italic {
		C.pango_font_description_set_style(desc, C.PANGO_STYLE_ITALIC)
	}
	if f.Vertical {
		C.pango_font_description_set_gravity(desc, C.PANGO_GRAVITY_EAST)
	}
	C.pango_layout_set_font_description(s.pl, desc)
	C.pango_font_description_free(desc)		// copy owned by s.pl according to the pangocairo docs
}

func (s *sysImage) text(str string, x int, y int) {
	C.cairo_save(s.cr)
	C.cairo_move_to(s.cr, C.double(x), C.double(y))
	cstr := C.CString(str)
	C.pango_layout_set_text(s.pl, cstr, -1)
	C.free(unsafe.Pointer(cstr))
	C.pango_cairo_show_layout(s.cr, s.pl)
	C.cairo_restore(s.cr)
}
