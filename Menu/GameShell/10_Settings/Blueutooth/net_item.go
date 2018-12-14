package Bluetooth

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
  
  Path string ///org/bluez/hci0/dev_34_88_5D_97_FF_26
  
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


func (self *NetItem) Init(path string ) {
  


}
