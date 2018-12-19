package Languages

import (
  "github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/rect"

  "github.com/cuu/LauncherGoDev/sysgo/UI"

)

type ListPageSelector struct {
  UI.InfoPageSelector
}

func NewListPageSelector() *ListPageSelector {
  
  p := &ListPageSelector{}
  
  p.Width = UI.Width
  p.BackgroundColor = &color.Color{131,199,219,255} //SkinManager().GiveColor('Front')
  
  return p

}


func (self *ListPageSelector) Draw() {
  
  idx := self.Parent.GetPsIndex()
  mylist := self.Parent.GetMyList()
 
  if idx < len(mylist) {
    x,y := mylist[idx].Coord()
    _,h := mylist[idx].Size()
    
    self.PosX = x+2
    self.PosY = y+1
    self.Height = h-3
    
    canvas_ := self.Parent.GetCanvasHWND()
    rect_   := rect.Rect(self.PosX,self.PosY,self.Width-4, self.Height)
    
    draw.AARoundRect(canvas_,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  } 
}

type PageListItem struct {
  UI.InfoPageListItem
  
  Active bool
  
}


func NewPageListItem() *PageListItem {
  
  p := &PageListItem{}
  p.Height = UI.DefaultInfoPageListItemHeight
  p.ReadOnly = false
	p.Labels = make(map[string]LabelInterface)
	p.Icons  = make( map[string]IconItemInterface)
	p.Fonts  = make(map[string]*ttf.Font)  

  return p
}


func (self *PageListItem) Draw() {
    
  x,_ := self.Labels["Text"].Coord()
  w,h := self.Labels["Text"].Size()
  
  self.Labels["Text"].NewCoord( x, self.PosY + (self.Height - h)/2 )
  
  
  if self.Active == true {
    self.Parent.(*LanguagesPage).Icons["done"].NewCoord(self.Parent.(*LanguagesPage).Width-30,self.PosY+5)
    self.Parent.(*LanguagesPage).Icons["done"].Draw()
  }
  
  self.Labels["Text"].SetBold(self.Active)
  self.Labels["Text"].Draw()
  
  
  
  if _, ok := self.Labels["Small"]; ok {
    x,_ = self.Labels["Small"].Coord()
    w,h = self.Labels["Small"].Size()
    
    self.Labels["Small"].NewCoord( self.Width - w - 10 , self.PosY + (self.Height - h)/2 )
    self.Labels["Small"].Draw()
    
  }
  
  canvas_ := self.Parent.GetCanvasHWND()
  draw.Line(canvas_, &color.Color{169,169,169,255}, 
    self.PosX, self.PosY+self.Height -1,
    self.PosX + self.Width, self.PosY+self.Height -1 ,1)
  
}

type LanguagesPage struct {
  UI.Page
  
  ListFont *ttf.Font
  
  BGwidth  int 
  BGheight  int
  
  DrawOnce bool
  
  Scroller *UI.ListScroller

  Icons map[string]UI.IconItemInterface
}


func NewLanguagesPage() *LanguagesPage {
  p := &LanguagesPage{}
  
  p.ListFont = UI.Fonts["notosanscjk15"]
  p.FootMsg = [5]string{"Nav","","","Back","Select"}
  
  p.BGwidth = UI.Width 
  p.BGheight = UI.Height - 24 - 20
  
  p.Icons = make(map[string]UI.IconItemInterface)
  
  return p
}

func (self *LanguagesPage) GenList() {

  self.MyList = nil
  
  start_x := 0 
  start_y := 0 
  last_height := 0
  
  file_paths,err := filepath.Glob("sysgo/langs/*.ini")//sorted
  
  if err == nil {
    for i,u := range file_paths {
      li := NewPageListItem()
      
    
    }
    
  }

}


