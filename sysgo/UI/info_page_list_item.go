package UI

import (
	"github.com/veandco/go-sdl2/ttf"
	
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	
)

type InfoPageListItem struct {
  ListItem
  Flag string 
  ReadOnly bool
}

func NewInfoPageListItem() *InfoPageListItem {
  p := &InfoPageListItem{}
  p.Height = 30
  p.ReadOnly = false
	p.Labels = make(map[string]LabelInterface)
	p.Icons  = make( map[string]IconItemInterface)
	p.Fonts  = make(map[string]*ttf.Font)
  
  return p
}

func (self *InfoPageListItem) SetSmallText(text string) {
	l := NewLabel()
	l.PosX = 40
	l.SetCanvasHWND(self.Parent.GetCanvasHWND())
	l.Init(text,self.Fonts["small"],nil)
	self.Labels["Small"] = l
}

func (self *InfoPageListItem) Init(text string ) {
	l := NewLabel()
	l.PosX = 10
	l.SetCanvasHWND(self.Parent.GetCanvasHWND())
	l.Init(text,self.Fonts["normal"],nil)
	self.Labels["Text"] = l
}

func (self *InfoPageListItem) Draw() {
  if self.ReadOnly == true {
    self.Labels["Text"].SetColor(&color.Color{130,130,130,255}  ) //SkinManager().GiveColor("ReadOnlyText")
  }else {
    self.Labels["Text"].SetColor(&color.Color{83,83,83,255} ) // SkinManager().GiveColor("Text")
  }
  
  x,_ := self.Labels["Text"].Coord()
  w,h := self.Labels["Text"].Size()
  
  self.Labels["Text"].NewCoord( x + self.PosX, self.PosY + (self.Height - h)/2 )
  
  self.Labels["Text"].Draw()
  
  self.Labels["Text"].NewCoord(x, self.PosY + (self.Height - h)/2 )
  
  if _, ok := self.Labels["Small"]; ok {
    x,_ = self.Labels["Small"].Coord()
    w,h = self.Labels["Small"].Size()
    
    self.Labels["Small"].NewCoord( self.Width - w - 5 , self.PosY + (self.Height - h)/2 )
    self.Labels["Small"].Draw()
    
  }
  
  canvas_ := self.Parent.GetCanvasHWND()
  draw.Line(canvas_, &color.Color{169,169,169,255}, self.PosX, self.PosY+self.Height -1,self.PosX + self.Width, self.PosY+self.Height -1 ,1)

}

