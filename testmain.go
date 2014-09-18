// 10 june 2014
package main

import (
	"fmt"
	"image"
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
	p := NewPen(PenSpec{
		R:			0,
		G:			128,
		B:			0,
		Line:			Solid,
		Thickness:	1,
	})
	f := NewFont(FontSpec{
		Family:	"Helvetica",
		Size:		12,
		Bold:		true,
	})
	i.Text("hello, world", 100, 20, f, p)
	p.Close()
	p = NewPen(PenSpec{
		R:			255,
		G:			0,
		B:			0,
		Line:			Solid,
		Thickness:	3,
	})
	i.Line(4, 4, 316, 236, p)
	i.Line(100, 20, 101, 21, p)
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

func main() {
	go myMain()
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
