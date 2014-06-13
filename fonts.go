// 13 june 2014
package main

// Font encodes information about a font.
type Font struct {
	Family		string
	Size			uint		// in points
	Bold			bool		// TODO can it be a factor? if not, what constitutes bold?
	Italic			bool		// italic == oblique if current backend differentiates
	Vertical		bool		// strictly gravity east/rotation 90 degrees clockwise? TODO
	Monospace	bool
}

// ListFonts computes a list of all fonts installed on the system.
// This recomputes the list on each call.
// The Size field of each returned Font shall be 0.
// Duplicates may be returned if information about the font is lost.
// TODO sort?
func ListFonts() []Font {
	return sysListFonts()
}

// Font selects the current Font for drawing.
// TODO behavior if Size == 0
func (i *Image) Font(f Font) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.sysImage.setFont(f)
}

// Text draws the given string at the given position on the current Image in the given Pen and Font.
// The top-left corner of the drawn string will be at the given point.
// TODO pango seems to do this vertically offset?
// TODO what if no Font was selected?
func (i *Image) Text(text string, x int, y int) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.sysImage.text(text, x, y)
}