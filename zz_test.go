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
	f := NewFont(FontSpec{
		Family:	"Helvetica",
		Size:		12,
		Bold:		true,
	})
	i := Line("hello, world", f, 0, 128, 0)
	f.Close()
	ui.Do(func() {
		w = ui.NewWindow("Test", 200, 200, ui.NewArea(i.Rect.Dx(), i.Rect.Dy(), &areaHandler{
			img:		i,
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
