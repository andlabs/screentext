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
func (a *areaHandler) Mouse(me ui.MouseEvent) bool {
	return false
}

func myMain() {
	fonts := ListFonts()
	for _, font := range fonts {
		fmt.Printf("%#v\n", font)
	}
	i := NewImage(320, 240)
	defer i.Close()
	i.Pen(NewRGBPen(255, 0, 0).Line(Solid, 3))
	i.Line(4, 4, 316, 236)
	w := ui.NewWindow("Test", 320, 240)
	w.Open(ui.NewArea(320, 240, &areaHandler{
		img:		i.Image(),
	}))
	<-w.Closing
}

func main() {
	err := ui.Go(myMain)
	if err != nil {
		panic(err)
	}
}
