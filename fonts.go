// 13 june 2014
package main

// Font encodes information about a font.
type Font struct {
	Family		string
	Bold			bool		// TODO can it be a factor? if not, what constitutes bold?
	Italic			bool		// italic == oblique if current backend differentiates
	Vertical		bool		// strictly gravity east/rotation 90 degrees clockwise? TODO
	Monospace	bool
}

// ListFonts computes a list of all fonts installed on the system.
// This recomputes the list on each call.
// Duplicates may be returned if information about the font is lost.
// TODO sort?
func ListFonts() []Font {
	return sysListFonts()
}
