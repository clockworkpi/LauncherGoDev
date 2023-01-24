package UI

import (
	"context"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	gotime "time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/zyxar/argo/rpc"
	
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/itchyny/volume-go"

	"github.com/vjeantet/jodaTime"

	"github.com/clockworkpi/LauncherGoDev/sysgo"
)

var TitleBar_BarHeight = 24

type TitleBarIconItem struct {
	MultiIconItem
	Parent *TitleBar
}

func NewTitleBarIconItem() *TitleBarIconItem {
	m := &TitleBarIconItem{}
	m.IconIndex = 0
	m.IconWidth = 18
	m.IconHeight = 18
	m.Align = ALIGN["VCenter"]
	return m

}

func (self *TitleBarIconItem) Adjust(x, y, w, h, at int) {
	self.PosX = x
	self.PosY = y
	self.Width = w
	self.Height = h
	self.AnimationTime = at

	if self.Label != nil {
		self.Label.SetCanvasHWND(self.Parent.CanvasHWND)
	}

	self.CreateImgSurf()
	//	self.AdjustLinkPage()

}

func (self *TitleBarIconItem) Draw() {
	parent_x, parent_y := self.Parent.PosX, self.Parent.PosY

	if self.Label != nil {
		//		lab_x,lab_y := self.Label.Coord()
		lab_w, lab_h := self.Label.Size()
		if self.Align == ALIGN["VCenter"] {
			self.Label.NewCoord(self.PosX-lab_w/2+parent_x, self.PosY+self.Height/2+6+parent_y)
		} else if self.Align == ALIGN["HLeft"] {
			self.Label.NewCoord(self.PosX+self.Width/2+3+parent_x, self.PosY-lab_h/2+parent_y)
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {

		portion := rect.Rect(0, self.IconIndex*self.IconHeight, self.IconWidth, self.IconHeight)

		surface.Blit(self.Parent.CanvasHWND,
			self.ImgSurf, draw.MidRect(self.PosX+parent_x, self.PosY+parent_y,
				self.Width, self.Height, Width, Height), &portion)
	}
}

type TitleBar struct {
	Widget
	BarHeight   int
	LOffset     int
	ROffset     int
	Icons       map[string]IconItemInterface
	IconWidth   int
	IconHeight  int
	BorderWidth int
	CanvasHWND  *sdl.Surface
	HWND        *sdl.Surface
	Title       string

	InLowBackLight int
	InAirPlaneMode bool
	
	WifiStrength int

	SkinManager *SkinManager //set by MainScreen

	icon_base_path string /// SkinMap("gameshell/titlebar_icons/")

	MyTimeLocation *gotime.Location

	TitleFont *ttf.Font
	TimeFont  *ttf.Font
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

	t.icon_base_path = SkinMap("sysgo/gameshell/titlebar_icons/")

	t.TitleFont = Fonts["varela16"]
	t.TimeFont = Fonts["varela12"]

	t.InLowBackLight = -1
	t.WifiStrength = 0 
	return t

}

func (self *TitleBar) Redraw() {
	self.UpdateDownloadStatus()
	SwapAndShow()
}

func (self *TitleBar) UpdateDownloadStatus() {
	
	rpcc, err := rpc.New(context.Background(), sysgo.Aria2Url, "", gotime.Second, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if resp,err := rpcc.GetGlobalStat();err == nil {
		num_active,_ := strconv.Atoi(resp.NumActive)
		
		if num_active > 0 {
			self.Icons["dlstatus"].SetIconIndex(1)
		}else if num_active == 0 {
			self.Icons["dlstatus"].SetIconIndex(0)
		}
	}
	
	
	defer rpcc.Close()

}

func (self *TitleBar) RoundRobinCheck() {
	for {

		if self.InLowBackLight < 0 {
			self.CheckBatteryStat()
			self.CheckBluetooth()
			self.UpdateWifiStrength()
			self.UpdateDownloadStatus()
			SwapAndShow()

		} else if self.InLowBackLight >= 0 {
			self.InLowBackLight += 1

			if self.InLowBackLight > 10 {
				self.CheckBatteryStat()
				self.CheckBluetooth()
				self.UpdateWifiStrength()
				self.UpdateDownloadStatus()
				self.InLowBackLight = 0 // reset
			}

		}

		gotime.Sleep(3000 * gotime.Millisecond)

	}
}

func (self *TitleBar) IsWifiConnectedNow() bool {
	cli := fmt.Sprintf("ip -4 addr show %s | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}'", sysgo.WifiDev)
	out := System(cli)
	if len(out) > 7 {
		if strings.Contains(out, "not") {
			return false
		} else {
			return true
		}
	}

	return false

}

func (self *TitleBar) UpdateWifiStrength() {
	self.GetWifiStrength()
	self.Draw(self.Title)
}

func (self *TitleBar) GetWifiStrength() int {
	qua := 0

	cli := fmt.Sprintf("sudo iwgetid %s -r",sysgo.WifiDev)
	out := System(cli)
	if len(out) > 2{
		out = strings.TrimSuffix(out, "\n")
		cli = fmt.Sprintf("sudo nmcli -t -f SSID,SIGNAL dev wifi list | grep \"^%s:\" | cut -d : -f 2",out)
		out = System(cli)

		out = strings.TrimSuffix(out, "\n")
		qua,_ = strconv.Atoi(out)
	}

	segs := [][]int{[]int{-2, -1}, []int{0, 25}, []int{25, 50}, []int{50, 75}, []int{75, 100}}
	stren_number := qua
	ge := 0
	if stren_number == 0 {
		return ge
	}

	for i, v := range segs {
		if stren_number >= v[0] && stren_number <= v[1] {
			ge = i
			break
		}
	}
	self.WifiStrength = ge
	return ge
}

func (self *TitleBar) SyncSoundVolume() {

	vol, err := volume.GetVolume()
	if err != nil {
		log.Printf("TitleBar SyncSoundVolume get volume failed: %+v\n", err)
		vol = 0
	}
	fmt.Printf("TitleBar SyncSoundVolume current volume: %d\n", vol)

	snd_segs := [][]int{[]int{0, 10}, []int{10, 30}, []int{30, 70}, []int{70, 100}}
	ge := 0

	for i, v := range snd_segs {
		if vol >= v[0] && vol <= v[1] {
			ge = i
			break
		}
	}

	self.Icons["soundvolume"].SetIconIndex(ge)
	self.Icons["sound"] = self.Icons["soundvolume"]
	//
}

// for outside widget to update sound icon
func (self *TitleBar) SetSoundVolume(vol int) {

	snd_segs := [][]int{[]int{0, 10}, []int{10, 30}, []int{30, 70}, []int{70, 100}}
	ge := 0

	for i, v := range snd_segs {
		if vol >= v[0] && vol <= v[1] {
			ge = i
			break
		}
	}

	self.Icons["soundvolume"].SetIconIndex(ge)
	self.Icons["sound"] = self.Icons["soundvolume"]

}

func (self *TitleBar) CheckBatteryStat() {
	bat_segs := [][]int{[]int{0, 6}, []int{7, 15}, []int{16, 20}, []int{21, 30}, []int{31, 50}, []int{51, 60}, []int{61, 80}, []int{81, 90}, []int{91, 100}}


	if FileExists(sysgo.Battery) == false {
		return
	}

	file, err := os.Open(sysgo.Battery)
	if err != nil {
		fmt.Println("Could not open file ", sysgo.Battery)
		return
	}

	defer file.Close()

	bat_uevent := make(map[string]string)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		pis := strings.Split(line, "=")
		if len(pis) > 1 {
			bat_uevent[pis[0]] = pis[1]
		}
	}

	cur_cap := 0

	if val, ok := bat_uevent["POWER_SUPPLY_CAPACITY"]; ok {
		cur_cap, _ = strconv.Atoi(val)
	} else {
		cur_cap = 0
	}

	cap_ge := 0

	for i, v := range bat_segs {
		if cur_cap >= v[0] && cur_cap <= v[1] {
			cap_ge = i
			break
		}
	}

	if val, ok := bat_uevent["POWER_SUPPLY_STATUS"]; ok {
		if val == "Charging" {
			self.Icons["battery"].SetIconIndex(1+cap_ge)
		} else {
			self.Icons["battery"].SetIconIndex(1+9+cap_ge)
		}
	}

}

func (self *TitleBar) SetBatteryStat(bat int) {

}

func (self *TitleBar) CheckBluetooth() {

	out := System("hcitool dev | grep hci0 |cut -f3")

	if len(out) < 17 {
		fmt.Println("Titlebar CheckBluetooth: no bluetooth", out)
		self.Icons["bluetooth"].SetIconIndex(2)
		return
	} else {
		out = System("sudo rfkill list | grep hci0 -A 2 | grep yes")
		if len(out) > 10 {
			self.Icons["bluetooth"].SetIconIndex(1)
			return
		}
	}

	self.Icons["bluetooth"].SetIconIndex(0)

}

func (self *TitleBar) Init(main_screen *MainScreen) {

	start_x := 0

	self.CanvasHWND = surface.Surface(self.Width, self.Height)
	self.HWND = main_screen.HWND
	self.SkinManager = main_screen.SkinManager

	icon_wifi_status := NewTitleBarIconItem()

	icon_wifi_status.MyType = ICON_TYPES["STAT"]
	icon_wifi_status.ImageName = self.icon_base_path + "wifi.png"
	icon_wifi_status.Parent = self

	icon_wifi_status.Adjust(start_x+self.IconWidth+5, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2, self.IconWidth, self.IconHeight, 0)

	self.Icons["wifistatus"] = icon_wifi_status

	battery_unknown := NewTitleBarIconItem()
	battery_unknown.MyType = ICON_TYPES["STAT"]
	battery_unknown.Parent = self
	battery_unknown.ImageName = self.icon_base_path + "battery.png"
	battery_unknown.Adjust(start_x+self.IconWidth+self.IconWidth+8, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2, self.IconWidth, self.IconHeight, 0)

	self.Icons["battery"] = battery_unknown

	self.CheckBatteryStat()

	sound_volume := NewTitleBarIconItem()
	sound_volume.MyType = ICON_TYPES["STAT"]
	sound_volume.Parent = self
	sound_volume.ImageName = self.icon_base_path + "soundvolume.png"
	sound_volume.Adjust(start_x+self.IconWidth+self.IconWidth+8, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2, self.IconWidth, self.IconHeight, 0)

	self.Icons["soundvolume"] = sound_volume

	self.SyncSoundVolume()

	bluetooth := NewTitleBarIconItem()
	bluetooth.MyType = ICON_TYPES["STAT"]
	bluetooth.Parent = self
	bluetooth.ImageName = self.icon_base_path + "bluetooth.png"
	bluetooth.Adjust(start_x+self.IconWidth+self.IconWidth+8, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2, self.IconWidth, self.IconHeight, 0)

	self.Icons["bluetooth"] = bluetooth
	self.CheckBluetooth()

	round_corners := NewTitleBarIconItem()
	round_corners.IconWidth = 10
	round_corners.IconHeight = 10

	round_corners.MyType = ICON_TYPES["STAT"]
	round_corners.Parent = self
	round_corners.ImgSurf = MyIconPool.GetImgSurf("roundcorners")
	round_corners.Adjust(0, 0, 10, 10, 0)

	self.Icons["round_corners"] = round_corners

	dlstatus := NewTitleBarIconItem()
	dlstatus.MyType = ICON_TYPES["STAT"]
	dlstatus.Parent = self
	if FileExists(self.icon_base_path + "dlstatus18.png") {
		dlstatus.ImageName = self.icon_base_path + "dlstatus18.png"
	}
	dlstatus.Adjust(start_x+self.IconWidth+self.IconWidth+8, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2, self.IconWidth, self.IconHeight, 0)
	self.Icons["dlstatus"] = dlstatus

	self.UpdateDownloadStatus()
	
	if self.IsWifiConnectedNow() {
		print("wifi is connected")
	} else {

		cmd := "sudo rfkill list | grep yes | cut -d \" \" -f3" //make sure sudo rfkill needs no password
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			fmt.Printf("Failed to execute command: %s\n", cmd)
		} else {
			outs := strings.Split(string(out), "\n")
			if len(outs) > 0 && outs[0] == "yes" {
				self.InAirPlaneMode = true
			} else {
				self.InAirPlaneMode = false
			}
		}

	}

	self.UpdateTimeLocation()

}

func (self *TitleBar) ClearCanvas() {
	surface.Fill(self.CanvasHWND, self.SkinManager.GiveColor("TitleBg"))

	self.Icons["round_corners"].NewCoord(5, 5)
	self.Icons["round_corners"].SetIconIndex(0)
	self.Icons["round_corners"].Draw()

	self.Icons["round_corners"].NewCoord(self.Width-5, 5)
	self.Icons["round_corners"].SetIconIndex(1)
	self.Icons["round_corners"].Draw()

}

func (self *TitleBar) UpdateTimeLocation() {

	d, err := ioutil.ReadFile("/etc/localtime")
	if err != nil {
		return
	}

	self.MyTimeLocation, err = gotime.LoadLocationFromTZData("local", d)
	if err != nil {
		fmt.Println(err)
		self.MyTimeLocation = nil
	}
}

func (self *TitleBar) GetLocalTime() gotime.Time {
	if self.MyTimeLocation == nil {
		return gotime.Now()
	} else {
		return gotime.Now().In(self.MyTimeLocation)
	}
}

func (self *TitleBar) Draw(title string) {
	self.ClearCanvas()
	self.Title = title
	
	cur_time := jodaTime.Format("HH:mm", self.GetLocalTime())

	time_text_w, time_text_h := font.Size(self.TimeFont, cur_time)
	title_text_w, title_text_h := font.Size(self.TitleFont, self.Title)

	title_text_surf := font.Render(self.TitleFont, self.Title, true, self.SkinManager.GiveColor("Text"), nil)

	surface.Blit(self.CanvasHWND, title_text_surf, draw.MidRect(title_text_w/2+self.LOffset, title_text_h/2+(self.BarHeight-title_text_h)/2, title_text_w, title_text_h, Width, Height), nil)

	time_text_surf := font.Render(self.TimeFont, cur_time, true, self.SkinManager.GiveColor("Text"), nil)
	surface.Blit(self.CanvasHWND, time_text_surf, draw.MidRect(Width-time_text_w/2-self.ROffset, time_text_h/2+(self.BarHeight-time_text_h)/2, time_text_w, time_text_h, Width, Height), nil)


	start_x := Width - time_text_w - self.ROffset - self.IconWidth*3 // close to the time_text

	self.Icons["bluetooth"].NewCoord(start_x-self.IconWidth, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2)
	self.Icons["sound"].NewCoord(start_x, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2)
	self.Icons["battery"].NewCoord(start_x+self.IconWidth+self.IconWidth+8, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2)

	if self.IsWifiConnectedNow() == true {
		ge := self.WifiStrength
		//fmt.Println("wifi ge: ",ge)
		if ge > 0 {
			self.Icons["wifistatus"].SetIconIndex(ge)
			self.Icons["wifistatus"].NewCoord(start_x+self.IconWidth+5, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2)
			self.Icons["wifistatus"].Draw()
		} else {
			self.Icons["wifistatus"].SetIconIndex(0)
			self.Icons["wifistatus"].NewCoord(start_x+self.IconWidth+5, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2)
			self.Icons["wifistatus"].Draw()
		}
	} else {

		self.Icons["wifistatus"].SetIconIndex(0)

		self.Icons["wifistatus"].NewCoord(start_x+self.IconWidth+5, self.IconHeight/2+(self.BarHeight-self.IconHeight)/2)

		self.Icons["wifistatus"].Draw()
	}

	self.Icons["sound"].Draw()
	self.Icons["battery"].Draw()

	self.Icons["bluetooth"].Draw()

	draw.Line(self.CanvasHWND, self.SkinManager.GiveColor("Line"), 0, self.BarHeight, self.Width, self.BarHeight, self.BorderWidth)

	if self.HWND != nil {
		rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
		surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
	}

	title_text_surf.Free()
	time_text_surf.Free()
}
