package UI

import (
  //"fmt"
//	"github.com/veandco/go-sdl2/ttf"

//	"github.com/cuu/gogame/surface"
//	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/draw"
	
)

type InfoPageSelector struct {
  PageSelector
  BackgroundColor *color.Color
}

func NewInfoPageSelector() *InfoPageSelector {
  p := &InfoPageSelector{}
  
  p.Width = Width
  p.BackgroundColor = &color.Color{131,199,219,255} //SkinManager().GiveColor('Front')
  
  return p
  
}

func (self *InfoPageSelector) AnimateDraw(x2, y2 int) {
  //pass
}

func (self *InfoPageSelector) Draw() {
  idx := self.Parent.GetPsIndex()
  mylist := self.Parent.GetMyList()
 
  if idx < len(mylist) {
    _,y := mylist[idx].Coord()
    _,h := mylist[idx].Size()
    
    
    self.PosY = y+1
    self.Height = h-3
    
    canvas_ := self.Parent.GetCanvasHWND()
    rect_   := rect.Rect(self.PosX,self.PosY,self.Width-4, self.Height)
    
    draw.AARoundRect(canvas_,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  }
}

 
