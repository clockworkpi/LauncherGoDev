package Bluetooth

import (

  bleapi "github.com/muka/go-bluetooth/api"
  "github.com/muka/go-bluetooth/bluez/profile"
  "github.com/cuu/LauncherGoDev/sysgo/UI"
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
  
  Channel string //'10'
  Stren   string //19%
  
  Icons map[string]UI.IconItemInterface
  Labels map[string]UI.LabelInterface
  
  IsActive bool 
  FontObj *ttf.Font
  RSSI int // 0
  MacAddr  string // 
  Parent *BluetoothPage
  Path string ///org/bluez/hci0/dev_34_88_5D_97_FF_26
  Props *profile.Device1Properties
  Device  *bleapi.Device
}


func NewNetItem() *NetItem {
  p:= &NetItem{}
  
	p.Height = NetItemDefaultHeight 
	
	p.Labels = make(map[string]UI.LabelInterface)
	p.Icons = make( map[string]UI.IconItemInterface)  
  
  return p

}


func (self *NetItem) SetActive(act bool) {
  self.IsActive = act

}


func (self *NetItem) Init( _label string) {
  
  self.MacAddr = self.Props.Address
  self.SetActive(self.Props.Connected)
  
  
  name_label := UI.NewLabel()
  name_label.PosX = 12
  
  name_label.CanvasHWND = self.Parent.CanvasHWND
  
  mac_addr := self.MacAddr
  if len(self.Props.Name) > 3 {
    mac_addr = self.Props.Name
  }
  
  self.RSSI = int(self.Props.RSSI)
  
  name_label.Init(mac_addr,self.FontObj,nil)
  
  self.Labels["mac_addr"] = name_label
  
  done_icon := NewNetItemIcon()
	done_icon.ImgSurf = UI.MyIconPool.GetImgSurf("done")
	done_icon.CanvasHWND = self.Parent.GetCanvasHWND()
	done_icon.Parent = self

	self.Icons["done"] = done_icon

}

func (self *NetItem) Connect() {
  
  if self.Device != nil {
    self.Device.Connect()
  }
}


func (self *NetItem) Draw() {
  for k,v := range self.Labels {
    x,y := v.Coord()
    _,h := v.Size()
    self.Labels[k].NewCoord(x, self.PosY+(self.Height - h)/2)
    self.Labels[k].Draw()
  }
  
  if self.IsActive {
    self.Icons["done"].NewCoord(UI.Width-22, self.PosY)
    self.Icons["done"].Draw()
  }

  draw.Line(self.Parent.CanvasHWND,&color.Color{169,169,169,255},
          self.PosX,self.PosY+self.Height-1,
          self.PosX+self.Width,self.PosY+self.Height-1,
          1)
  
}
