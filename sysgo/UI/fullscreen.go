package UI

import (
	"github.com/cuu/gogame/surface"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/cuu/gogame/rect"
)

type FullScreen struct {
	Widget
	CanvasHWND *sdl.Surface
	HWND       *sdl.Surface
}

func NewFullScreen() *FullScreen {
	p := &FullScreen{}
	p.Width = Width
	p.Height = Height

	return p

}

func (self *FullScreen) SwapAndShow() {
	if self.HWND != nil {
		rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
		surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
		DisplayFlip()
	}

}

func (self *FullScreen) Draw() {

}
