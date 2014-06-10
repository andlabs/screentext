// 10 june 2014
package main

import (
	"sync"
)

// Pen represents a pen.
// Pens are used to draw lines, shape outlines, etc.
type Pen struct {
	lock		sync.Mutex
	sysPen	*sysPen
}

// Line represents a style of line for a Pen.
type Line uint
const (
	Solid Line = iota
)

// NewRGBPen creates a new Pen with the given opaque color.
// r, g, and b are in the range [0,255].
func NewRGBPen(r uint, g uint, b uint) *Pen {
	return &Pen{
		sysPen:	mkSysPenRGB(r, g, b),
	}
}

// Line sets the line type and thickness, in pixels, of p.
// It then returns p.
func (p *Pen) Line(linetype Line, thickness uint) *Pen {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.sysPen.setLineType(linetype, thickness)
	return p
}
