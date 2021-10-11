package UI

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	//	"github.com/veandco/go-sdl2/ttf"

	//	"github.com/cuu/gogame/surface"
	//	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
	//	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/draw"
)

type ListScroller struct {
	Widget
	MinHeight int
	Parent    PageInterface
	Color     *color.Color

	StartX     int
	StartY     int
	EndX       int
	EndY       int
	Value      int
	CanvasHWND *sdl.Surface
}

func NewListScroller() *ListScroller {
	l := &ListScroller{}
	l.Width = 7
	l.Color = &color.Color{131, 199, 219, 255} // SkinManager().GiveColor('Front')
	return l
}

func (self *ListScroller) Init() {
	//just set the CanvasHWND
	cav_ := self.Parent.GetCanvasHWND()
	self.SetCanvasHWND(cav_)
}

func (self *ListScroller) SetCanvasHWND(canvas *sdl.Surface) {
	self.CanvasHWND = canvas
}

func (self *ListScroller) AnimateDraw(x2, y2 int) {

}

func (self *ListScroller) UpdateSize(bigheight, dirtyheight int) {
	_, h_ := self.Parent.Size()

	bodyheight := float64(h_) / float64(bigheight)
	if bodyheight > 1.0 {
		bodyheight = 1.0
	}

	margin := 4

	self.Height = int(bodyheight*float64(h_) - float64(margin))

	if self.Height < self.MinHeight {
		self.Height = self.MinHeight
	}

	self.StartX = self.Width / 2
	self.StartY = margin/2 + self.Height/2

	self.EndX = self.Width / 2
	self.EndY = h_ - margin/2 - self.Height/2

	process := float64(dirtyheight) / float64(bigheight)

	value := process * float64(self.EndY-self.StartY)

	self.Value = int(value)

}

func (self *ListScroller) Draw() {
	w_, h_ := self.Parent.Size()

	start_rect := draw.MidRect(self.PosX+self.StartX, self.StartY+self.Value, self.Width, self.Height, w_, h_)

	if self.Width < 1 {
		fmt.Println("ListScroller width error")
	} else {
		draw.AARoundRect(self.CanvasHWND, start_rect, self.Color, 3, 0, self.Color)
	}
}
