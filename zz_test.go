// 10 june 2014

package ndraw

import (
	"fmt"
	"image"
	"testing"
	"github.com/andlabs/ui"
)

type areaHandler struct {
	img		*image.RGBA
}
func (a *areaHandler) Paint(r image.Rectangle) *image.RGBA {
	return a.img.SubImage(r).(*image.RGBA)
}
func (a *areaHandler) Key(ke ui.KeyEvent) bool {
	return false
}
func (a *areaHandler) Mouse(me ui.MouseEvent) {}

var w ui.Window

func myMain() {
	fonts := ListFonts()
	for _, font := range fonts {
		fmt.Printf("%#v\n", font)
	}
	i := NewImage(320, 240)
	defer i.Close()
	b := NewBrush(BrushSpec{
		R:			0,
		G:			128,
		B:			0,
		A:			0xFF,
	})
	f := NewFont(FontSpec{
		Family:	"Helvetica",
		Size:		12,
		Bold:		true,
	})
	i.Text("hello, world", 100, 20, f, nil, b)
	b.Close()
	p := NewPen(PenSpec{
		R:			255,
		G:			0,
		B:			0,
		A:			128,
		Line:			Solid,
		Thickness:	3,
	})
	i.Line(4, 4, 316, 236, p)
	i.Line(100, 20, 101, 21, p)
	p.Close()
	p = NewPen(PenSpec{
		R:			0,
		G:			0,
		B:			255,
		A:			255,
		Line:			Solid,
		Thickness:	2,
	})
	wid, ht := i.TextSize("hello, world", f)
	i.Line(100, 20 + ht + 10, 100 + wid, 20 + ht + 10, p)
	i.Line(100 + wid + 10, 20, 100 + wid + 10, 20 + ht, p)
	p = NewPen(PenSpec{
		R:			0,
		G:			0,
		B:			255,
		A:			255,
		Line:			Solid,
		Thickness:	1,
	})
	i.Line(100, 20 + ht, 100 + wid, 20 + ht, p)
	i.Line(100 + wid, 20, 100 + wid, 20 + ht, p)
	p.Close()
	f.Close()
	ui.Do(func() {
		w = ui.NewWindow("Test", 320, 240, ui.NewArea(320, 240, &areaHandler{
			img:		i.Image(),
		}))
		w.OnClosing(func() bool {
			ui.Stop()
			return true
		})
		w.Show()
	})
}

func init() {
	go myMain()
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}

func TestDummy(t *testing.T) {
	// do nothing
}
