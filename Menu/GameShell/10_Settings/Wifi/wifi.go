package Wifi

//wifi_list.py

import (
	"fmt"
	"strconv"
	"strings"
	//"os"
	// "os/exec"
	// gotime "time"
	"log"
	//"github.com/godbus/dbus"
	
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	
	"github.com/veandco/go-sdl2/ttf"

	"github.com/clockworkpi/LauncherGoDev/sysgo"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/time"

	wifi "github.com/cuu/wpa-connect"
)

const EMPTY_NETWORK = "00:00:00:00:00:00"

type WifiDisconnectConfirmPage struct {
	UI.ConfirmPage
	Parent *WifiInfoPage
}

func NewWifiDisconnectConfirmPage() *WifiDisconnectConfirmPage {
	p := &WifiDisconnectConfirmPage{}
	p.ListFont = UI.Fonts["veramono20"]
	p.FootMsg = [5]string{"Nav", "", "", "Cancel", "Yes"}

	p.ConfirmText = "Confirm Disconnect?"
	return p
}

func (self *WifiDisconnectConfirmPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["B"] {
		fmt.Println("Disconnecting..")
		self.SnapMsg("Disconnecting...")
		self.Screen.Draw()
		self.Screen.SwapAndShow()

		self.Parent.Parent.Disconnect()

		time.BlockDelay(400)

		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
		self.Parent.Parent.Rescan(false)
	}
}

type WifiInfoPage struct {
	UI.Page
	ListFontObj *ttf.Font
	Bss         *wifi.BSS

	AList  map[string]map[string]string
	ESSID  string
	BSSID  string
	MyList []UI.ListItemInterface

	DisconnectConfirmPage *WifiDisconnectConfirmPage //child page
	Parent                *WifiList
}

func NewWifiInfoPage() *WifiInfoPage {
	p := &WifiInfoPage{}
	p.FootMsg = [5]string{"Nav", "Disconnect", "", "Back", ""}

	p.ListFontObj = UI.Fonts["varela15"]

	p.AList = make(map[string]map[string]string)

	p.BSSID = ""
	p.ESSID = ""
	return p

}

func (self *WifiInfoPage) GenList() {

	self.MyList = nil
	self.MyList = make([]UI.ListItemInterface, 0)

	if self.BSSID != "" {
		self.AList["ip"]["value"] = "Not Connected"
		if self.BSSID == self.Parent.CurBssid {
			var ip string
			ip = self.Parent.GetWirelessIP()
			if len(ip) > 0 {
				self.AList["ip"]["value"] = ip
			}
		} else {
			fmt.Println(self.BSSID)
		}

		self.AList["ssid"]["value"] = self.ESSID
	}

	start_x := 0
	start_y := 0
	i := 0
	for k, _ := range self.AList {
		li := UI.NewInfoPageListItem()
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y + i*li.Height //default is 30
		li.Width = UI.Width
		li.Fonts["normal"] = self.ListFontObj
		li.Fonts["small"] = UI.Fonts["varela12"]

		if self.AList[k]["label"] != "" {
			li.Init(self.AList[k]["label"])
		} else {
			li.Init(self.AList[k]["key"])
		}

		li.Flag = self.AList[k]["key"]

		li.SetSmallText(self.AList[k]["value"])
		self.MyList = append(self.MyList, li)
		i += 1
	}

}

func (self *WifiInfoPage) Init() {
	if self.Screen != nil {
		if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
			self.CanvasHWND = self.Screen.CanvasHWND
		}
	}
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	ps := UI.NewInfoPageSelector()
	ps.Parent = self
	ps.PosX = 2
	self.Ps = ps
	self.PsIndex = 0

	ip := make(map[string]string) // ip = {}
	ip["key"] = "ip"
	ip["label"] = "IP"
	ip["value"] = "Not Connected"

	ssid := make(map[string]string) // ssid = {}
	ssid["key"] = "ssid"
	ssid["label"] = "SSID"
	ssid["value"] = ""

	self.AList["ip"] = ip
	self.AList["ssid"] = ssid

	self.DisconnectConfirmPage = NewWifiDisconnectConfirmPage()
	self.DisconnectConfirmPage.Screen = self.Screen
	self.DisconnectConfirmPage.Name = "Confirm Disconnect"
	self.DisconnectConfirmPage.Parent = self

	self.DisconnectConfirmPage.Init()

}

func (self *WifiInfoPage) ScrollUp() {
	if len(self.MyList) == 0 {
		return
	}

	self.PsIndex -= 1

	if self.PsIndex < 0 {
		self.PsIndex = 0
	}
	cur_li := self.MyList[self.PsIndex]
	x, y := cur_li.Coord()
	if x < 0 {
		for i := 0; i < len(self.MyList); i++ {
			_, h := self.MyList[i].Size()
			x, y = self.MyList[i].Coord()
			self.MyList[i].NewCoord(x, y+h)
		}
	}
}

func (self *WifiInfoPage) ScrollDown() {
	if len(self.MyList) == 0 {
		return
	}

	self.PsIndex += 1
	if self.PsIndex >= len(self.MyList) {
		self.PsIndex = len(self.MyList) - 1
	}

	cur_li := self.MyList[self.PsIndex]
	x, y := cur_li.Coord()
	_, h := cur_li.Size()

	if y+h > self.Height {
		for i := 0; i < len(self.MyList); i++ {
			_, h = self.MyList[i].Size()
			x, y = self.MyList[i].Coord()
			self.MyList[i].NewCoord(x, y-h)
		}
	}
}

func (self *WifiInfoPage) Click() {
	/*
	   cur_li = self._MyList[self._PsIndex]
	   print(cur_li._Flag)
	*/
}

func (self *WifiInfoPage) TryDisconnect() {
	var ip string

	ip = self.Parent.GetWirelessIP()

	if len(ip) > 6 {
		self.Screen.PushPage(self.DisconnectConfirmPage)
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	} else {
		fmt.Println("WifiInfoPage TryDisconnect can not get IP,maybe you are offline")
		return
	}
}

func (self *WifiInfoPage) OnLoadCb() {

	/*
	  self.FootMsg[1]="Disconnect"
	  self.FootMsg[1] = ""
	*/

	self.GenList()
}

func (self *WifiInfoPage) OnReturnBackCb() {

	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()

}

func (self *WifiInfoPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["Up"] {
		self.ScrollUp()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["Down"] {
		self.ScrollDown()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["Enter"] {
		self.Click()
	}

	if ev.Data["Key"] == UI.CurKeys["X"] {
		self.TryDisconnect()
	}
}

func (self *WifiInfoPage) Draw() {
	self.ClearCanvas()
	self.Ps.Draw()

	for i := 0; i < len(self.MyList); i++ {
		self.MyList[i].Draw()
	}
}

type WifiListSelector struct {
	UI.PageSelector
	BackgroundColor *color.Color

	Parent *WifiList
}

func NewWifiListSelector() *WifiListSelector {
	p := &WifiListSelector{}
	p.BackgroundColor = &color.Color{131, 199, 219, 255} //SkinManager().GiveColor('Front')

	return p
}

func (self *WifiListSelector) Draw() {
	idx := self.Parent.PsIndex
	if idx < len(self.Parent.MyList) {
		x := self.Parent.MyList[idx].PosX + 11
		y := self.Parent.MyList[idx].PosY + 1
		h := self.Parent.MyList[idx].Height - 3

		self.PosX = x
		self.PosY = y
		self.Height = h

		rect_ := rect.Rect(x, y, self.Width, h)
		draw.AARoundRect(self.Parent.CanvasHWND, &rect_, self.BackgroundColor, 4, 0, self.BackgroundColor)
	}
}

type WifiListMessageBox struct {
	UI.Label
	Parent *WifiList
}

func NewWifiListMessageBox() *WifiListMessageBox {
	p := &WifiListMessageBox{}
	p.Color = &color.Color{83, 83, 83, 255}
	return p
}

func (self *WifiListMessageBox) Draw() {
	my_text := font.Render(self.FontObj, self.Text, true, self.Color, nil)

	w := surface.GetWidth(my_text)
	h := surface.GetHeight(my_text)

	x := (self.Parent.Width - w) / 2
	y := (self.Parent.Height - h) / 2

	padding := 10

	white := &color.Color{255, 255, 255, 255}
	black := &color.Color{0, 0, 0, 255}

	rect_ := rect.Rect(x-padding, y-padding, w+padding*2, h+padding*2)

	draw.Rect(self.CanvasHWND, white, &rect_, 0)
	draw.Rect(self.CanvasHWND, black, &rect_, 1)

	rect_2 := rect.Rect(x, y, w, h)
	surface.Blit(self.CanvasHWND, my_text, &rect_2, nil)
	my_text.Free()
}

//---------WifiList---------------------------------
type BlockCbFunc func()

type WifiList struct {
	UI.Page
	WifiPassword string
	Connecting   bool
	Scanning     bool

	ShowingMessageBox bool
	MsgBox            *WifiListMessageBox
	ConnectTry        int

	BlockingUI bool
	BlockCb    BlockCbFunc

	LastStatusMsg string
	Scroller      *UI.ListScroller
	ListFontObj   *ttf.Font

	InfoPage *WifiInfoPage

	MyList   []*NetItem
	CurEssid string ///SomeWifi
	CurBssid string //00:00:00:00:00:00
	CurIP    string
	CurSig   string
}

func NewWifiList() *WifiList {
	p := &WifiList{}
	p.ListFontObj = UI.Fonts["notosanscjk15"]
	p.FootMsg = [5]string{"Nav.", "Scan", "Info", "Back", "Enter"}

	return p
}

func (self *WifiList) ShowBox(msg string) {
	self.MsgBox.Text = msg
	self.ShowingMessageBox = true
	self.Screen.Draw()
	self.MsgBox.Draw()
	self.Screen.SwapAndShow()

}

func (self *WifiList) HideBox() {
	self.Draw()
	self.ShowingMessageBox = false
	self.Screen.SwapAndShow()
}

func (self *WifiList) GenNetworkList() {
	self.MyList = self.MyList[:0]

	start_x := 0
	start_y := 0

	var is_active bool
	var li_idx int
	li_idx = 0

	self.WifiScanStarted()

	if bssList, err := GsScanManager.Scan(); err == nil {

		self.CurEssid = GsScanManager.GetCurrentSSID()
		self.CurBssid = GsScanManager.GetCurrentBSSID()

		for _, bss := range bssList {
			is_active = false
			fmt.Println(bss.SSID, " ", bss.BSSID, " ", bss.Signal, bss.KeyMgmt)
			ni := NewNetItem()
			ni.Parent = self
			ni.PosX = start_x
			ni.PosY = start_y + li_idx*NetItemDefaultHeight
			ni.Width = UI.Width
			ni.FontObj = self.ListFontObj
			ni.Essid = bss.SSID
			ni.Bssid = bss.BSSID
			ni.Signal = bss.Signal

			if self.CurBssid == ni.Bssid {
				is_active = true
			}

			ni.Init(is_active)
			self.MyList = append(self.MyList, ni)

			li_idx++
		}
	}

	self.WifiScanFinished()

	self.PsIndex = 0
}

func (self *WifiList) Disconnect() {
	self.Connecting = false
  //nmcli -t -f NAME c show --active
	//nmcli con down
	cli := "nmcli -t -f NAME c show --active"
	out := UI.SystemTrim(cli)
	
	cli = fmt.Sprintf("sudo nmcli con down \"%s\"",out)

	out = UI.System(cli)
	log.Println(out)
	
	self.CurEssid = ""
	self.CurBssid = ""

}

func (self *WifiList) ShutDownConnecting() {

	self.Connecting = false
	self.Disconnect()
}

func (self *WifiList) Rescan(sync bool) { // sync default should be false
	fmt.Println("start Rescan")
	self.GenNetworkList()
}

// dbus signal functions
func (self *WifiList) WifiScanFinished() {
	if self.Screen.CurrentPage != self {
		return
	}

	self.ResetPageSelector()

	self.Scanning = false
	self.HideBox()

	self.BlockingUI = false
	fmt.Println("dbus says scan finished")

}

func (self *WifiList) WifiScanStarted() {
	if self.Screen.CurrentPage != self {
		return
	}

	self.Scanning = true
	self.ShowBox("Wifi scanning...")
	self.BlockingUI = true
	fmt.Println("dbus says start scan")
}

func (self *WifiList) SaveNetworkList() {
	
}

func (self *WifiList) SaveWifiPassword(essid,password string) {

	db, err := sql.Open("sqlite3", sysgo.SQLDB)
	if err != nil {
		log.Fatal(err)
		return 
	}
	defer db.Close()

	stmt, err := db.Prepare("select count(*) from wifi where essid = ?")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	var count string
	err = stmt.QueryRow(essid).Scan(&count)
	if err != nil {
		log.Println(err)
		count = "0"
	}

	cnt,_ := strconv.Atoi(count)
	if cnt > 0 {
		_,err = db.Exec("update wifi set pass= :pass where essid = :essid",sql.Named("pass",password),sql.Named("essid",essid))
		if err != nil {
			log.Println(err)
		}
	}else {
		_,err = db.Exec("insert into wifi(essid,pass) values(:essid,:pass)",sql.Named("essid",essid),sql.Named("pass",password))
		if err != nil {
			log.Println(err)
		}
	}

}

func (self *WifiList) LoadWifiPassword(essid string) string {
	db, err := sql.Open("sqlite3", sysgo.SQLDB)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer db.Close()

	password := ""
	stmt, err := db.Prepare("select pass from wifi where essid = ?")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}else {
		err = stmt.QueryRow(essid).Scan(&password)
		if err != nil {
			log.Println(err)
		}
	}
	return password
}
//----------------------------------------------------------------------------------

func (self *WifiList) UpdateNetList(state int, info []string, force_check bool, firstrun bool) { //force_check default ==false, firstrun default == false

	if force_check == true {
		self.GenNetworkList()
		self.SaveNetworkList()
	}

}

func (self *WifiList) UpdateListActive() {

	for i := 0; i < len(self.MyList); i++ {
		if self.MyList[i].Bssid == self.CurBssid {
			self.MyList[i].IsActive = true
		} else {
			self.MyList[i].IsActive = false
		}
	}
}

func (self *WifiList) ConfigWireless(password string) {

	ssid := self.MyList[self.PsIndex].Essid
	fmt.Println(ssid)
	fmt.Println(password)
	self.ShowBox("Connecting...")

	self.Connecting = true
	cli := fmt.Sprintf("sudo nmcli dev wifi connect %s password \"%s\"", ssid, password)
	log.Println(cli)
	out := UI.System(cli)
	log.Println(out)
	if strings.Contains(out, "successfully") {
		self.CurEssid = self.MyList[self.PsIndex].Essid
		self.CurBssid = self.MyList[self.PsIndex].Bssid
		self.MyList[self.PsIndex].Password = password
		self.SaveWifiPassword(ssid,password)
		self.ShowBox("Connected")
	} else {
		self.ShowBox("Wifi connect error")
		self.CurEssid = ""
		self.CurBssid = ""
	}

	self.Connecting = false

	self.UpdateListActive()

}

func (self *WifiList) GetWirelessIP() string {

	cli := fmt.Sprintf("ip -4 addr show %s | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}'", sysgo.WifiDev)
	out := UI.SystemTrim(cli)

	return out

}

func (self *WifiList) ScrollUp() {
	if len(self.MyList) == 0 {
		return
	}

	self.PsIndex -= 1
	if self.PsIndex < 0 {
		self.PsIndex = 0
	}

	cur_ni := self.MyList[self.PsIndex] //*NetItem
	if cur_ni.PosY < 0 {
		for i := 0; i < len(self.MyList); i++ {
			self.MyList[i].PosY += self.MyList[i].Height
		}
	}
}

func (self *WifiList) ScrollDown() {
	if len(self.MyList) == 0 {
		return
	}

	self.PsIndex += 1
	if self.PsIndex >= len(self.MyList) {
		self.PsIndex = len(self.MyList) - 1
	}

	cur_ni := self.MyList[self.PsIndex]
	if cur_ni.PosY+cur_ni.Height > self.Height {
		for i := 0; i < len(self.MyList); i++ {
			self.MyList[i].PosY -= self.MyList[i].Height
		}
	}

}

func (self *WifiList) AbortedAndReturnToUpLevel() {
	self.HideBox()
	self.Screen.FootBar.ResetNavText()
	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()
}

func (self *WifiList) OnKbdReturnBackCb() {
	password_inputed := strings.Join(APIOBJ.PasswordPage.Textarea.MyWords, "")
	fmt.Println("Password inputed: ", password_inputed)
	
	if len(password_inputed) > 4 {
		self.ConfigWireless(password_inputed)
	}else {
		log.Println("wifi password length too short ",len(password_inputed))
	}
}

func (self *WifiList) OnReturnBackCb() {
	//fmt.Println("return back")
}

func (self *WifiList) KeyDown(ev *event.Event) {
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {

		//self.ShutDownConnecting()
		//self.ShowBox("ShutDownConnecting...")
		self.AbortedAndReturnToUpLevel()

	}

	if ev.Data["Key"] == UI.CurKeys["Up"] {
		self.ScrollUp()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["Down"] {
		self.ScrollDown()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["Enter"] { // enter to set password,enter is B on GM
		if len(self.MyList) == 0 {
			return
		}
		if self.MyList[self.PsIndex].IsActive == true {
			var ip string
			ip = self.GetWirelessIP()
			self.ShowBox(ip)
		} else {
			self.Screen.PushCurPage()
			self.Screen.SetCurPage(APIOBJ.PasswordPage)

			thepass := self.LoadWifiPassword(self.MyList[self.PsIndex].Essid)

			fmt.Println("APIOBJ.PasswordPage.SetPassword ", thepass, len(thepass))
			APIOBJ.PasswordPage.SetPassword(thepass)

			self.Screen.Draw()
			self.Screen.SwapAndShow()

		}
	}

	if ev.Data["Key"] == UI.CurKeys["X"] {
		self.Rescan(false)
	}

	if ev.Data["Key"] == UI.CurKeys["Y"] {
		if len(self.MyList) == 0 {
			return
		}
		self.InfoPage.BSSID = self.MyList[self.PsIndex].Bssid
		self.InfoPage.ESSID = self.MyList[self.PsIndex].Essid
		self.Screen.PushPage(self.InfoPage)
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

}

func (self *WifiList) OnLoadCb() {

	ip := self.GetWirelessIP()
	if len(ip) < 6 {
		self.CurEssid = ""
		self.CurBssid = ""
		self.CurIP = ip
	}
	self.Rescan(false)
}

func (self *WifiList) Init() {

	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.CanvasHWND = self.Screen.CanvasHWND

	ps := NewWifiListSelector()
	ps.Parent = self
	ps.Width = UI.Width - 12

	self.Ps = ps
	self.PsIndex = 0

	msgbox := NewWifiListMessageBox()
	msgbox.CanvasHWND = self.CanvasHWND
	msgbox.Init(" ", UI.Fonts["veramono12"], nil)
	msgbox.Parent = self

	self.MsgBox = msgbox

	/*
	   {
	   'fields': [],
	   'name': 'WPA 1/2 (Passphrase)',
	   'optional': [],
	   'protected': [
	   ['apsk', 'Preshared_Key'],
	   ],
	   'required': [
	   ['apsk', 'Preshared_Key'],
	   ],
	   'type': 'wpa-psk',
	   },
	*/

	self.Scroller = UI.NewListScroller()
	self.Scroller.Parent = self
	self.Scroller.PosX = 2
	self.Scroller.PosY = 2
	self.Scroller.Init()

	self.InfoPage = NewWifiInfoPage()
	self.InfoPage.Screen = self.Screen
	self.InfoPage.Name = "Wifi info"
	self.InfoPage.Parent = self
	self.InfoPage.Init()

}

func (self *WifiList) Draw() {
	self.ClearCanvas()

	if len(self.MyList) == 0 {
		return
	}

	self.Ps.Draw()

	for _, v := range self.MyList {
		v.Draw()
	}

	self.Scroller.UpdateSize(len(self.MyList)*NetItemDefaultHeight, self.PsIndex*NetItemDefaultHeight)
	self.Scroller.Draw()

}
