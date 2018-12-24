package TimeZone

import (
  //"fmt"
  //"strings"
  //"io/ioutil"
  "path/filepath"
  "github.com/veandco/go-sdl2/ttf"
  
  //"github.com/veandco/go-sdl2/sdl"
	//"github.com/cuu/gogame/surface"
  //"github.com/cuu/gogame/rect"
 	"github.com/cuu/gogame/color"

	"github.com/cuu/gogame/draw"
  "github.com/cuu/LauncherGoDev/sysgo/UI"

)

var TimeZoneListPageListItemDefaultHeight = 30

type TimeZoneListPageInterface interface {  
  UI.PageInterface
  GetMapIcons() map[string]UI.IconItemInterface
  
}

type TimeZoneListPageListItem struct {
  UI.HierListItem
  Parent TimeZoneListPageInterface
}

func NewTimeZoneListPageListItem() *TimeZoneListPageListItem {
  p := &TimeZoneListPageListItem{}
  p.Labels = make(map[string]UI.LabelInterface)
	p.Icons  = make( map[string]UI.IconItemInterface)
	p.Fonts  = make(map[string]*ttf.Font)
  
  p.MyType = UI.ICON_TYPES["EXE"]
	p.Height = TimeZoneListPageListItemDefaultHeight
	p.Width  = 0
 
  
  return p
}

func (self *TimeZoneListPageListItem) Init(text string) {
  l := UI.NewLabel()
  l.PosX = 20

  l.SetCanvasHWND(self.Parent.GetCanvasHWND())
  
  if self.IsDir() == true || self.IsFile() == true {
    self.Path = text
  }
  
  label_text := filepath.Base(text)
  
  if self.IsDir() == true {
    l.Init(label_text, self.Fonts["normal"],nil)
  }else {
    l.Init(label_text,self.Fonts["normal"],nil)
  }
  
  self.Labels["Text"] = l
}

func (self *TimeZoneListPageListItem) Draw() {
  x,y := self.Labels["Text"].Coord()
  _,h := self.Labels["Text"].Size()
  
  
  if self.Path != "[..]" {
    self.Labels["Text"].NewCoord(23,y)
    
  }else {
    self.Labels["Text"].NewCoord(3,y)
  }
  
  x,y = self.Labels["Text"].Coord()
  self.Labels["Text"].NewCoord(x, self.PosY + (self.Height-h)/2)
  
  self.Labels["Text"].Draw()
  
  parent_icons := self.Parent.GetMapIcons()
  _,h = parent_icons["sys"].Size()
  
  if self.IsDir() == true && self.Path != "[..]" {
    parent_icons["sys"].SetIconIndex (0)
    parent_icons["sys"].NewCoord(self.PosX+12,self.PosY+(self.Height-h)/2+h/2)
    parent_icons["sys"].Draw()
  }
  
  if self.IsFile() == true {
    parent_icons["sys"].SetIconIndex(1)
    parent_icons["sys"].NewCoord(self.PosX+12,self.PosY+(self.Height-h)/2+h/2)
    parent_icons["sys"].Draw()
  }
  
  draw.Line(self.Parent.GetCanvasHWND(),&color.Color{169,169,169,255},
    self.PosX,self.PosY+self.Height-1,self.PosX+self.Width,self.PosY+self.Height-1,1)

}


