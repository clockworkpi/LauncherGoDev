package UI

import (
	"fmt"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/image"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
)

type MultiIconItem struct {
	IconItem

	IconWidth  int
	IconHeight int
}

func NewMultiIconItem() *MultiIconItem {
	m := &MultiIconItem{}
	m.IconIndex = 0
	m.IconWidth = 18
	m.IconHeight = 18
	return m
}

func (self *MultiIconItem) CreateImgSurf() {
	if self.ImgSurf == nil && self.ImageName != "" {
		self.ImgSurf = image.Load(self.ImageName)
	}
}

func (self *MultiIconItem) DrawTopLeft() {
	if self.Parent == nil {
		fmt.Println("Error: MultiIconItem Draw Parent nil")
		return
	}

	parent_x, parent_y := self.Parent.Coord()

	if self.Label != nil {
		lab_w, lab_h := self.Label.Size()
		if self.Align == ALIGN["VCenter"] {

			self.Label.NewCoord(self.PosX-lab_w/2+parent_x, self.PosY+self.Height/2+6+parent_y)
		} else if self.Align == ALIGN["HLeft"] {
			self.Label.NewCoord(self.PosX+self.Width/2+3+parent_x, self.PosY-lab_h/2+parent_y)
		}
		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		portion := rect.Rect(0, self.IconIndex*self.IconHeight, self.IconWidth, self.IconHeight)
		rect_ := rect.Rect(self.PosX+parent_x, self.PosY+parent_y, self.Width, self.Height) // DIFF
		surface.Blit(self.Parent.GetCanvasHWND(), self.ImgSurf, &rect_, &portion)
	}

}

func (self *MultiIconItem) Draw() {
	if self.Parent == nil {
		fmt.Println("Error: MultiIconItem Draw Parent nil")
		return
	}
	parent_x, parent_y := self.Parent.Coord()

	if self.Label != nil {
		//		lab_x,lab_y := self.Label.Coord()
		lab_w, lab_h := self.Label.Size()
		if self.Align == ALIGN["VCenter"] {
			self.Label.NewCoord(self.PosX-lab_w/2+parent_x, self.PosY+self.Height/2+6+parent_y)
		} else if self.Align == ALIGN["HLeft"] {
			self.Label.NewCoord(self.PosX+self.Width/2+3+parent_x, self.PosY-lab_h/2+parent_y)
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {

		portion := rect.Rect(0, self.IconIndex*self.IconHeight, self.IconWidth, self.IconHeight)

		surface.Blit(self.Parent.GetCanvasHWND(),
			self.ImgSurf, draw.MidRect(self.PosX+parent_x, self.PosY+parent_y,
				self.Width, self.Height, Width, Height), &portion)
	}
}
