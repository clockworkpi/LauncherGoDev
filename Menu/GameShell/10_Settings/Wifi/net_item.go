package Wifi

import (
  "fmt"
  //"strconv"
  //"strings"
  
  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/ttf"
  "github.com/cuu/gogame/color"
  "github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/rect"
  "github.com/cuu/gogame/surface"
  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
  
	
)
var NetItemDefaultHeight = 30

type NetItemMultiIcon struct {
	UI.MultiIconItem
	CanvasHWND *sdl.Surface      // self._Parent._CanvasHWND
	Parent     UI.WidgetInterface //
}

func NewNetItemMultiIcon() *NetItemMultiIcon{
  p := &NetItemMultiIcon{}
  p.IconIndex = 0
  p.IconWidth = 18
  p.IconHeight = 18  
  
  p.Width  = 18
  p.Height = 18
  return p
}

func (self *NetItemMultiIcon) Draw() {
	_,h_ := self.Parent.Size()	
	dest_rect := rect.Rect(self.PosX,self.PosY+(h_-self.Height)/2, self.Width,self.Height)
	area_rect := rect.Rect(0,self.IconIndex*self.IconHeight,self.IconWidth,self.IconHeight)
	surface.Blit(self.CanvasHWND,self.ImgSurf,&dest_rect,&area_rect)
		
}

type NetItemIcon struct {
	UI.IconItem
	CanvasHWND *sdl.Surface
	Parent UI.WidgetInterface
}

func NewNetItemIcon() *NetItemIcon {
	p := &NetItemIcon{}
	p.Width = 18
	p.Height = 18
	return p
}

func (self *NetItemIcon) Draw() {
	_,h_ := self.Parent.Size()

	dest_rect := rect.Rect(self.PosX,self.PosY+(h_-self.Height)/2,self.Width,self.Height)

	surface.Blit(self.CanvasHWND,self.ImgSurf,&dest_rect,nil)
	
}


type NetItem struct {
	UI.Widget

	Bssid string //eg: 50:3A:A0:51:18:3C
	Essid string //eg: MERCURY_EB88

	ip string
	Encrypt   string // WPA2
	Channel   string //'10'
	Signal    int16 // -67
	Mode      string // Master or AdHoc
	Parent    *WifiList
	IsActive  bool
	Password  string 
	Labels map[string]UI.LabelInterface
	Icons  map[string]UI.IconItemInterface
	Fonts  map[string]*ttf.Font
	FontObj *ttf.Font
	
  
}

func NewNetItem() *NetItem {
	p := &NetItem{}

	p.Height = NetItemDefaultHeight 
	
	p.Labels = make(map[string]UI.LabelInterface)
	p.Icons = make( map[string]UI.IconItemInterface)
	p.Fonts = make(map[string]*ttf.Font)
	
	return p
}

func (self *NetItem) SetActive( act bool) {
	self.IsActive = act
}

func (self *NetItem) UpdateStrenLabel() { //  ## sig_str should be 'number',eg:'-70'
    
	if _, ok := self.Labels["stren"]; ok {
	  self.Labels["stren"].SetText( fmt.Sprintf("%d",self.CalcWifiQuality()) )
	}
	
}

func (self *NetItem) Init(is_active bool) {

  //strenstr := "quality"
  //gap := 7
  
  the_main_screen := self.Parent.GetScreen()
  
  if is_active {
	self.SetActive(is_active)
  }

  
  essid_label := UI.NewLabel()
  essid_label.PosX = 36
  essid_label.CanvasHWND = self.Parent.GetCanvasHWND()

  essid_  := ""
	
  if len(self.Essid) > 19 {
		essid_ = self.Essid[:20]
  }else {
		essid_ = self.Essid
  }

	if len(essid_) == 0 {
		essid_ = self.Bssid
	}
	
	if len(essid_) == 0 {
		essid_ = EMPTY_NETWORK
	}
	
	//fmt.Println("essid: ",essid_, len(essid_))
	
  essid_label.Init(essid_, self.FontObj,nil)
	
  self.Labels["essid"] = essid_label

  stren_label := UI.NewLabel()
  stren_label.CanvasHWND = self.Parent.GetCanvasHWND()

	stren_l := fmt.Sprintf("%%%d ",self.CalcWifiQuality())
	if len(stren_l) == 0 {
		stren_l = "%%0"
	}
  stren_label.Init(stren_l, self.FontObj,nil)
  stren_label.PosX = self.Width - 23 - stren_label.Width-2

  self.Labels["stren"] = stren_label

  lock_icon := NewNetItemIcon()
  lock_icon.ImgSurf = UI.MyIconPool.GetImgSurf("lock")
  lock_icon.CanvasHWND = self.Parent.GetCanvasHWND()
  lock_icon.Parent = self // WidgetInterface
  self.Icons["lock"] = lock_icon

  done_icon := NewNetItemIcon()
  done_icon.ImgSurf = UI.MyIconPool.GetImgSurf("done")
  done_icon.CanvasHWND = self.Parent.GetCanvasHWND()
  done_icon.Parent = self

  self.Icons["done"] = done_icon

  nimt := NewNetItemMultiIcon()
  nimt.ImgSurf = the_main_screen.TitleBar.Icons["wifistatus"].GetImgSurf()
  nimt.CanvasHWND = self.Parent.GetCanvasHWND()
  nimt.Parent = self // WidgetInterface

  self.Icons["wifistatus"] = nimt
}


func (self *NetItem) Connect() {

	
}

func (self *NetItem) CalcWifiQuality() int {

  qua := 0
  qua = 2 * (int(self.Signal) + 100)
  
  return qua
}

func (self *NetItem) CalcWifiStrength() int {
  
  segs := [][]int{ []int{-2,-1}, []int{0,25}, []int{25,50}, []int{50,75},[]int{75,100}}
  stren_number :=  self.CalcWifiQuality()
  ge := 0
  if stren_number == 0 {
    return ge
  }
    
  for i,v := range segs {
    if stren_number >= v[0] && stren_number <= v[1] {
      ge = i
      break
    }
  }

  return ge
  
}

func (self *NetItem) Draw() {
  for i,v := range self.Labels {
		x_,_ := v.Coord()
		_,h_  := v.Size()
		self.Labels[i].NewCoord(x_,self.PosY+(self.Height - h_)/2)
		self.Labels[i].Draw()
  }

  if self.IsActive == true {
		self.Icons["done"].NewCoord(14,self.PosY)
		self.Icons["done"].Draw()
  }

  /*
  if self.Encrypt != "Unsecured" {
	w_,_ := self.Labels["stren"].Size()
	self.Icons["lock"].NewCoord(self.Width -23 - w_ -2 - 18, self.PosY)
    self.Icons["lock"].Draw()
  }
  */
  //the_main_screen := self.Parent.GetScreen()
  ge := self.CalcWifiStrength()
  if ge > 0 {
    self.Icons["wifistatus"].SetIconIndex(ge)
		self.Icons["wifistatus"].NewCoord(self.Width-23,self.PosY)
		self.Icons["wifistatus"].Draw()
  }else {
		self.Icons["wifistatus"].SetIconIndex(0)
		self.Icons["wifistatus"].NewCoord(self.Width-23,self.PosY)
		self.Icons["wifistatus"].Draw()
  }

  draw.Line(self.Parent.GetCanvasHWND(),
		&color.Color{169,169,169,255},
		self.PosX,self.PosY+self.Height-1,
		self.PosX+self.Width,self.PosY+self.Height-1,
		1)
	
}
