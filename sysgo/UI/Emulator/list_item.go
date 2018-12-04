package Emulator

import (
  "path/filepath"
  
	"github.com/cuu/gogame/surface"
  "github.com/cuu/LauncherGoDev/sysgo/UI"

)

type EmulatorPageInterface {
  UI.PageInterface
  GetIcons() map[string]UI.IconItemInterface
}


type ListItemIcon struct {
  UI.IconItem

}

func NewListItemIcon() *ListItemIcon {
  p := &ListItemIcon{}
  p.MyType = UI.ICON_TYPES["EXE"]

	p.Align = UI.ALIGN["VCenter"]
  
  p.Width = 18
  p.Height = 18
  
  return p
}

func (self *ListItemIcon) Draw() {
  _,h := self.Parent.Size()
  
  rect_ := rect.Rect(self.PosX,self.PosY+(h-self.Height)/2,self.Width,self.Height)
  
  surface.Blit(self.CanvasHWND, self.ImgSurf,&rect_,nil)
}

/// [..] [.] 
type HierListItem struct {
  UI.ListItem
  MyType int
  Path  string
  Active bool
  Playing bool
}

var HierListItemDefaultHeight = 32

func NewHierListItem() *HierListItem {
  p := &HierListItem{}
  p.Labels = make(map[string]LabelInterface)
	p.Icons  = make( map[string]IconItemInterface)
	p.Fonts  = make(map[string]*ttf.Font)
  
  p.MyType = UI.ICON_TYPES["EXE"]
	p.Height = HierListItemDefaultHeight
	p.Width  = 0
  
  return p
}

func (self *HierListItem) IsFile() bool {
  if self.MyType == UI.ICON_TYPES["FILE"] {
    return true
  }
  
  return false
}


func (self *HierListItem) IsDir() bool {
  if self.MyType == UI.ICON_TYPES["DIR"] {
    return true
  }
  
  return false
}


func (self *HierListItem) Init(text string) {
  l := UI.NewLabel()
  l.PosX = 20
  l.SetCanvasHWND(self.Parent.GetCanvasHWND())
  
  if self.IsDir() == true || self.IsFile() == true {
    self.Path = text
  }
  
  label_text := filepath.Base(text)
  ext:= filepath.Ext(text)
  if ext != "" {
    alias_file := strings.Replace(text,ext,"",-1) + ".alias"
    
    if UI.FileExists(alias_file) == true {
      b, err := ioutil.ReadFile(alias_file) 
      if err != nil {
        fmt.Print(err)
      }else {
        label_text = string(b)
      }
    }
    
  }
  
  if self.IsDir() == true {
    l.Init(label_text, self.Fonts["normal"])
  }else {
    l.Init(label_text,self.Fonts["normal"])
  }
  
  self.Labels["Text"] = l
}

func (self *HierListItem) Draw() {
  
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
  
  
  /*
  w,h := self.Parent.Icons["sys"].Size()
  
  if self.IsDir() == true && self.Path != "[..]" {
    self.Parent.Icons["sys"].IconIndex = 0
    self.Parent.Icons["sys"].NewCoord(self.PosX+12,self.PosY+(self.Height-h)/2+h/2)
    self.Parent.Icons["sys"].Draw()
  }
  
  if self.IsFile() == true {
    self.Parent.Icons["sys"].IconIndex = 1
    self.Parent.Icons["sys"].NewCoord(self.PosX+12,self.PosY+(self.Height-h)/2+h/2)
    self.Parent.Icons["sys"].Draw()
  }
  */
  
  draw.Line(self.Parent.GetCanvasHWND(),&color.Color{169,169,169,255},
    self.PosX,self.PosY+self.Height-1,self.PosX+self.Width,self.PosY+self.Height-1,1)
  
}

type EmulatorListItem struct {
  HierListItem
  Parent EmulatorPageInterface
}

func NewEmulatorListItem() *EmulatorListItem {
  p := &EmulatorListItem{}
  p.Labels = make(map[string]LabelInterface)
	p.Icons  = make( map[string]IconItemInterface)
	p.Fonts  = make(map[string]*ttf.Font)
  
  p.MyType = UI.ICON_TYPES["EXE"]
	p.Height = 32
	p.Width  = 0  
  return p
}

func (self *EmulatorListItem) Draw() {
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
  
  parent_icons := self.Parent.GetIcons()
  w,h := parent_cons["sys"].Size()
  
  if self.IsDir() == true && self.Path != "[..]" {
    parent_icons["sys"].IconIndex = 0
    parent_icons["sys"].NewCoord(self.PosX+12,self.PosY+(self.Height-h)/2+h/2)
    parent_icons["sys"].Draw()
  }
  
  if self.IsFile() == true {
    parent_icons["sys"].IconIndex = 1
    parent_icons["sys"].NewCoord(self.PosX+12,self.PosY+(self.Height-h)/2+h/2)
    parent_icons["sys"].Draw()
  }
  
  draw.Line(self.Parent.GetCanvasHWND(),&color.Color{169,169,169,255},
    self.PosX,self.PosY+self.Height-1,self.PosX+self.Width,self.PosY+self.Height-1,1)

}


