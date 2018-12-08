package Wifi

import (
  "fmt"
  "strconv"
  "strings"
  
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
  "github.com/cuu/gogame/color"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/LauncherGoDev/sysgo/UI"
  
  "github.com/cuu/LauncherGoDev/sysgo/DBUS"
	
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

	dhcphostname string //"GameShell"
	ip string
	dns_domain string
	gateway string
	use_global_dns int // eg 0==False,1 == True
	netmask string
	usedhcphostname int
	bitrate string //"auto"
	dns3 string
	dns2 string
	dns1 string
	use_settings_globally int // 0
	use_static_dns int //eg: 1== True
	search_domain string

	Encrypt   string // WPA2
	Channel   string //'10'
	Stren     string // 19%
	NetId     int
	Mode      string // Master or AdHoc
	Parent    *WifiList
	IsActive  bool
	
	Labels map[string]UI.LabelInterface
	Icons  map[string]UI.IconItemInterface
	Fonts  map[string]*ttf.Font
	FontObj *ttf.Font
	
  Wireless *DBUS.DbusInterface
  Daemon   *DBUS.DbusInterface
  
}

func NewNetItem() *NetItem {
	p := &NetItem{}
	p.NetId = 0
	p.bitrate = "auto"
	p.dhcphostname = "GameShell"

	p.Height = NetItemDefaultHeight 
	
	p.Labels = make(map[string]UI.LabelInterface)
	p.Icons = make( map[string]UI.IconItemInterface)
	p.Fonts = make(map[string]*ttf.Font)
	
	return p
}

func (self *NetItem) SetActive( act bool) {
	self.IsActive = act
}

func (self *NetItem) UpdateStrenLabel( strenstr string) { //  ## strenstr should be 'number',eg:'90'

	self.Stren = strenstr

	if _, ok := self.Labels["stren"]; ok {
		self.Labels["stren"].SetText(self.Stren)
	}
	
}

func (self *NetItem) Init(i int,is_active bool) {

	var sig_display_type int
	strenstr := "quality"
	gap := 4
  
  the_main_screen := self.Parent.GetScreen()
  
  self.Wireless = self.Parent.Wireless
  self.Daemon   = self.Parent.Daemon
    
	self.Daemon.Get( self.Daemon.Method("GetSignalDisplayType"), &sig_display_type )
  
	if sig_display_type == 0 {
		strenstr = "quality"
		gap = 4 // Allow for 100%
	}else {
		strenstr = "strength"
		gap = 7 //  -XX dbm = 7
	}
  
	self.NetId = i
    
	tmp := 0
	self.Wireless.Get(self.Wireless.Method("GetWirelessProperty",self.NetId, strenstr),&tmp)
  tmp2 := ""
	self.Daemon.Get( self.Daemon.Method("FormatSignalForPrinting",tmp), &tmp2)
  
	self.Stren = tmp2

	self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetId,"essid"),&self.Essid)
	self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetId,"bssid"),&self.Bssid)
  
	check_enc := false
	self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetId,"encryption"),&check_enc)

	if check_enc == true {
		self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetId,"encryption_method"),&self.Encrypt)
	}else {
		self.Encrypt = "Unsecured"
	}

	self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetId,"mode"),&self.Mode)

	self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetId,"channel"),&self.Channel)

	theString := fmt.Sprintf("  %-*s %25s %9s %17s %6s %4s",gap,self.Stren,self.Essid,self.Encrypt,self.Bssid,self.Mode,
		self.Channel)
  
  
	if is_active {
		theString = ">> " + theString[1:]
		self.SetActive(is_active)
	}

  //fmt.Println(theString)
  
	essid_label := UI.NewLabel()
	essid_label.PosX = 36
	essid_label.CanvasHWND = self.Parent.GetCanvasHWND()

	essid_  := ""
	
	if len(self.Essid) > 19 {
		essid_ = self.Essid[:20]
	}else {
		essid_ = self.Essid
	}
 
	essid_label.Init(essid_, self.FontObj,nil)

	self.Labels["essid"] = essid_label

	stren_label := UI.NewLabel()
	stren_label.CanvasHWND = self.Parent.GetCanvasHWND()

	stren_label.Init(self.Stren, self.FontObj,nil)
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

	self.Wireless.Method("ConnectWireless",self.NetId)
	
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

	if self.Encrypt != "Unsecured" {
		w_,_ := self.Labels["stren"].Size()
		self.Icons["lock"].NewCoord(self.Width -23 - w_ -2 - 18, self.PosY)
		self.Icons["lock"].Draw()
	}

	stren_int,err := strconv.ParseInt(strings.Replace(self.Stren,"%","",-1),10,64)
	if err == nil {
    the_main_screen := self.Parent.GetScreen()
		ge := the_main_screen.TitleBar.GetWifiStrength(int(stren_int))
		if ge > 0 {
			self.Icons["wifistatus"].SetIconIndex(ge)
			self.Icons["wifistatus"].NewCoord(self.Width-23,self.PosY)
			self.Icons["wifistatus"].Draw()
		}else {
			self.Icons["wifistatus"].SetIconIndex(0)
			self.Icons["wifistatus"].NewCoord(self.Width-23,self.PosY)
			self.Icons["wifistatus"].Draw()
		}
	}

	draw.Line(self.Parent.GetCanvasHWND(),
		&color.Color{169,169,169,255},
		self.PosX,self.PosY+self.Height-1,
		self.PosX+self.Width,self.PosY+self.Height-1,
		1)
}
