package UI

import (
	"log"
	"strconv"
	"bufio"
	"strings"
	
	"github.com/veandco/go-sdl2/sdl"
	
	"github.com/cuu/gogame/surface"
	
	"github.com/itchyny/volume-go"

	"../../sysgo"
)


var TitleBar_BarHeight = 24

type TitleBarIconItem struct {
	MultiIconItem
	Parent *TitleBar
}

func NewTitleBarIconItem() *TitleBarIconItem {
	m := &TitleBarIconItem{}

	return m

}

func (self *TitleBarIconItem) Draw() {
	parent_x,parent_y := self.Parent.PosX,self.Parent.PosY
	
	if self.Label != nil {
//		lab_x,lab_y := self.Label.Coord()
		lab_w,lab_h:= self.Label.Size()
		if self.Align == ALIGN["VCenter"] {
			self.Label.NewCoord( self.PosX - lab_w/2 + parent_x,        self.PosY + self.Height/2+6 + parent_y)
		}else if self.Align == ALIGN["HLeft"] {
			self.Label.NewCoord( self.PosX + self.Width/2+3 + parent_x, self.PosY - lab_h/2 + parent_y )
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		
		portion := rect.Rect(0,self.IconIndex*self.IconHeight,self.IconWidth,self.IconHeight)
		
		surface.Blit(self.Parent.GetCanvasHWND(),
			self.ImgSurf,draw.MidRect(self.PosX + parent_x, self.PosY + parent_y,
			self.Width,self.Height, Width, Height),&portion)
	}
}


type TitleBar struct {

	PosX int
	PosY int
	Width int
	Height int
	BarHeight int
	LOffset int
	ROffset int
	Icons map[string]IconItemInterface
	IconWidth
	IconHeight
	BorderWidth
	CanvasHWND *sdl.Surface
	HWND       *sdl.Surface
	Title string
	InLowBackLight int
	SkinManager *SkinManager //set by MainScreen
	DBusManager *DBusInterface
	
	icon_base_path string /// SkinMap("gameshell/titlebar_icons/")
}


func NewTitleBar() *TitleBar {
	t := &TitleBar{}

	
	t.BorderWidth = 1

	t.BarHeight = TitleBar_BarHeight
	t.Height = t.BarHeight + t.BorderWidth

	t.Width = Width

	t.IconWidth = 18
	t.IconHeight = 18

	t.LOffset = 3
	t.ROffset = 3
	
	t.Icons = make(map[string]IconItemInterface)
	
	t.icon_base_path  = SkinMap("gameshell/titlebar_icons/")
}

func (t *TitleBar) RoundRobinCheck {
	
}

func (t *TitleBar) UpdateWifiStrength() {
	
}

func (t *TitleBar) GetWifiStrength(stren string) int {
	segs := [][]int{ []int{-2,-1}, []int{0,25}, []int{25,50}, []int{50,75},int{75,100}}
	stren_number,_ :=  strconv.Atoi( stren )
	ge := 0
	if stren_number == 0 {
		return ge
	}
	
	for i,v in range segs {
		if stren_number >= v[0] && stren_number <= v[1] {
			ge = i
			break
		}
	}

	return ge	
}

func (self *TitleBar) SyncSoundVolume() {
	
  vol, err := volume.GetVolume()
  if err != nil {
    log.Fatalf("get volume failed: %+v", err)
		vol = 0
  }
  fmt.Printf("current volume: %d\n", vol)

	snd_segs := [][]int{ []int{0,10}, []int{10,30}, []int{30,70},[]int{70,100} }
	ge := 0

	for i,v in range snd_segs {
		if vol >= v[0] && vol <= v[1] {
			ge = i
			break
		}
	}

	self.Icons["soundvolume"].SetIconIndex(ge)
	self.Icons["sound"] = self.Icons["soundvolume"]
	// 
}

func (t *TitleBar) SetSoundVolume(vol int) {
	//pass
}

func (self *TitleBar) CheckBatteryStat() {
	bat_segs:= [][]int{[]int{0,6},[]int{7,15},[]int{16,20},[]int{21,30},[]int{31,50},[]int{51,60},[]int{61,80},[]int{81,90},[]int{91,100}}
	
	file, err := os.Open( sysgo.Battery )
	if err != nil {
		fmt.Println("Could not open file ", sysgo.Battery)
		self.Icons["battery"] = self.Icons["battery_unknown"]
		return
	}

	defer file.Close()

	bat_uevent := make([string]string)
	
  scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines) 

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line," ")
		pis := strings.Split(line,"=")
		if len(pis) > 1 {
			bat_uevent[ pis[0] ] = pis[1]
		}
	}

	cur_cap := 0
	
	if val, ok := bat_uevent["POWER_SUPPLY_CAPACITY"]; ok {
		cur_cap = strings.Atoi(val)
	}else {
		cur_cap = 0
	}

	cap_ge := 0

	for i,v in range bat_segs {
		if cur_cap >= v[0] && cur_cap <= v[1] {
			cap_ge = i
			break
		}
	}

	if val, ok := bat_uevent["POWER_SUPPLY_STATUS"]; ok {
		if val == "Charging" {
			self.Icons["battery_charging"].SetIconIndex(cap_ge)
			self.Icons["battery"] = self.Icons["battery_charging"]
		}else {
			self.Icons["battery_charging"].SetIconIndex(cap_ge)
			self.Icons["battery"] = self.Icons["battery_discharging"]	
		}
	}
	
}

func (self *TitleBar) SetBatteryStat( bat int) {
	
}

func (self *TitleBar) Init(main_screen *MainScreen) {

	start_x := 0

	self.CanvasHWND = surface.Surface(self.Width,self.Height)
	self.HWND = main_screen.HWND
	self.SkinManager = main_screen.SkinManager
	self.DBusManager = main_screen.DBusManager
	
	icon_wifi_status := NewTitleBarIconItem()

	icon_wifi_status.MyType = ICON_TYPES["STAT"]
	icon_wifi_status.ImageName = self.icon_base_path+"wifi.png"
	icon_wifi_status.Parent = self

	icon_wifi_status.Adjust(start_x+self.IconWidth+5,self.IconHeight/2+(self.BarHeight-self.IconHeight)/2,self.IconWidth,self.IconHeight,0)

	self.Icons["wifistatus"] = icon_wifi_status

	battery_charging := NewTitleBarIconItem()
	battery_charging.MyType = ICON_TYPES["STAT"]
	battery_charging.Parent = self
	battery_charging.ImageName = self.icon_base_path+"withcharging.png"
	battery_charging.Adjust(start_x+self.IconWidth+self.IconWidth+8,self.IconHeight/2+(self.BarHeight-self.IconHeight)/2,self.IconWidth,self.IconHeight,0)

	self.Icons["battery_charging"] = battery_charging

	battery_discharging := NewTitleBarIconItem()
	battery_discharging.MyType = ICON_TYPES["STAT"]
	battery_discharging.Parent = self
	battery_discharging.ImageName = self.icon_base_path+"without_charging.png"
	battery_discharging.Adjust(start_x+self.IconWidth+self.IconWidth+8,self.IconHeight/2+(self.BarHeight-self.IconHeight)/2,self.IconWidth,self.IconHeight,0)

	self.Icons["battery_discharging"] = battery_discharging

	battery_unknown  := NewTitleBarIconItem()
	battery_unknown.MyType = ICON_TYPES["STAT"]
	battery_unknown.Parent = self
	battery_unknown.ImageName = self.icon_base_path+"battery_unknown.png"
	battery_unknown.Adjust(start_x+self.IconWidth+self.IconWidth+8,self.IconHeight/2+(self.BarHeight-self.IconHeight)/2,self.IconWidth,self.IconHeight,0)
	
	self.Icons["battery_unknown"] = battery_unknown

	self.CheckBatteryStat()

	sound_volume := NewTitleBarIconItem()
	sound_volume.MyType = ICON_TYPES["STAT"]
	sound_volume.Parent = self
	sound_volume.ImageName = self.icon_base_path+"soundvolume.png"
	sound_volume.Adjust(start_x+self.IconWidth+self.IconWidth+8,self.IconHeight/2+(self.BarHeight-self.IconHeight)/2,self.IconWidth,self.IconHeight,0)

	self.Icons["soundvolume"] = sound_volume

	self.SyncSoundVolume()

	round_corners := NewTitleBarIconItem()
	round_corners.IconWidth = 10
	round_corners.IconHeight = 10
	
	round_corners.MyType = ICON_TYPES["STAT"]
	round_corners.Parent = self
	round_corners.ImgSurf = MyIconPool.GetImageSurf["roundcorners"]
	round_corners.Adjust(0,0,10,10,0)
	
	self.Icons["round_corners"] = round_corners

	if is_wifi_connected_now() {
		print("wifi is connected")
		print( wifi_strength())
	}
}
