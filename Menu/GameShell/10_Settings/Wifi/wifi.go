package Wifi
//wifi_list.py

import (
    "fmt"
    //"strconv"
    "strings"
    "os"
    "os/exec"
    gotime "time"

    //"github.com/godbus/dbus"

    "github.com/veandco/go-sdl2/ttf"

    "github.com/cuu/gogame/surface"
    "github.com/cuu/gogame/font"
    "github.com/cuu/gogame/color"
    "github.com/cuu/gogame/event"
    "github.com/cuu/gogame/time"
    "github.com/cuu/gogame/rect"
    "github.com/cuu/gogame/draw"
    "github.com/clockworkpi/LauncherGoDev/sysgo"
    "github.com/clockworkpi/LauncherGoDev/sysgo/UI"

    wifi "github.com/mark2b/wpa-connect"
)

const EMPTY_NETWORK = "00:00:00:00:00:00"

type WifiDisconnectConfirmPage struct {
  UI.ConfirmPage
  Parent *WifiInfoPage
}

func cmdEnv() []string {
    return []string{"LANG=C", "LC_ALL=C"}
}

func execCmd(cmdArgs []string) ([]byte, error) {
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
    cmd.Env = append(os.Environ(), cmdEnv()...)
    out, err := cmd.Output()
    if err != nil {
        err = fmt.Errorf(`failed to execute "%v" (%+v)`, strings.Join(cmdArgs, " "), err)
    }
    return out, err
}


func NewWifiDisconnectConfirmPage() *WifiDisconnectConfirmPage {
  p := &WifiDisconnectConfirmPage{}  
  p.ListFont = UI.Fonts["veramono20"]
  p.FootMsg = [5]string{"Nav","","","Cancel","Yes"}

  p.ConfirmText ="Confirm Disconnect?"
  return p
}

func (self *WifiDisconnectConfirmPage) KeyDown(ev *event.Event ) {

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

  }
}

type WifiInfoPage struct {
  UI.Page
  ListFontObj  *ttf.Font
  Bss *wifi.BSS

  AList map[string]map[string]string
  NetworkId  string

  MyList []UI.ListItemInterface

  DisconnectConfirmPage *WifiDisconnectConfirmPage //child page 
  Parent *WifiList
}

func NewWifiInfoPage() *WifiInfoPage {
  p := &WifiInfoPage{}
  p.FootMsg = [5]string{"Nav","Disconnect","","Back",""}

  p.ListFontObj = UI.Fonts["varela15"]

  p.AList = make(map[string]map[string]string)

  p.NetworkId = EMPTY_NETWORK
  return p

}

func (self *WifiInfoPage) GetWirelessIP() string {

  return "0.0.0.0"
}

func (self *WifiInfoPage) GenList() {
    var cur_network_id string
    self.MyList = nil
    self.MyList = make([]UI.ListItemInterface,0)

    cur_network_id = EMPTY_NETWORK

    if self.NetworkId != EMPTY_NETWORK {
      self.AList["ip"]["value"] = "Not Connected"
      if cur_network_id == self.NetworkId {
        var ip string
        ip = self.GetWirelessIP()
        if len(ip) > 0 {
          self.AList["ip"]["value"]=ip
        }
      }

      self.AList["bssid"]["value"] = self.Parent.CurBssid
    }

    start_x := 0 
    start_y := 0 
    i := 0
    for k,_ := range self.AList {
      li := UI.NewInfoPageListItem()
      li.Parent = self
      li.PosX = start_x
      li.PosY = start_y + i * li.Height//default is 30
      li.Width = UI.Width
      li.Fonts["normal"] = self.ListFontObj
      li.Fonts["small"]  = UI.Fonts["varela12"]

      if self.AList[k]["label"] != "" {
        li.Init(self.AList[k]["label"])
      }else {
        li.Init(self.AList[k]["key"])
      }

      li.Flag = self.AList[k]["key"]

      li.SetSmallText(self.AList[k]["value"])
      self.MyList = append(self.MyList,li)
      i+=1
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

  bssid := make(map[string]string) // bssid = {}
  bssid["key"] = "bssid"
  bssid["label"] = "BSSID"
  bssid["value"] = ""

  self.AList["ip"] = ip
  self.AList["bssid"] = bssid

  self.GenList()

  self.DisconnectConfirmPage = NewWifiDisconnectConfirmPage()
  self.DisconnectConfirmPage.Screen = self.Screen
  self.DisconnectConfirmPage.Name   = "Confirm Disconnect"
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
  x,y := cur_li.Coord()
  if x < 0 {
    for i:=0;i<len(self.MyList);i++ {
      _,h := self.MyList[i].Size()
      x,y  = self.MyList[i].Coord()
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
  x,y  := cur_li.Coord()
  _,h  := cur_li.Size()

  if y + h > self.Height {
    for i:=0;i<len(self.MyList);i++ {
      _,h = self.MyList[i].Size()
      x,y = self.MyList[i].Coord()
      self.MyList[i].NewCoord(x, y - h)
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
 
  if len(ip) >  6 {
    self.Screen.PushPage(self.DisconnectConfirmPage)
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }else {
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

func (self *WifiInfoPage) KeyDown(ev *event.Event ) {

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

    for i:=0;i<len(self.MyList);i++ {
      self.MyList[i].Draw()
    }
}


type WifiListSelector struct{
  UI.PageSelector
    BackgroundColor *color.Color

    Parent *WifiList
}

func NewWifiListSelector() *WifiListSelector {
p := &WifiListSelector{}
   p.BackgroundColor = &color.Color{131,199,219,255} //SkinManager().GiveColor('Front')

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

    rect_ := rect.Rect(x,y,self.Width,h)
    draw.AARoundRect(self.Parent.CanvasHWND,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  }
}

type WifiListMessageBox struct{
  UI.Label
  Parent *WifiList
}

func NewWifiListMessageBox() *WifiListMessageBox{
p := &WifiListMessageBox{}
   p.Color = &color.Color{83,83,83,255}
   return p
}


func (self *WifiListMessageBox) Draw() {
  my_text := font.Render(self.FontObj,self.Text,true,self.Color,nil)

  w := surface.GetWidth(my_text)
  h := surface.GetHeight(my_text)

  x := (self.Parent.Width - w )/2
  y := (self.Parent.Height - h)/2

  padding := 10

  white := &color.Color{255,255,255,255}
  black := &color.Color{0,  0,  0,  255}

  rect_ := rect.Rect(x-padding,y-padding,w+padding*2,h+padding*2)

  draw.Rect(self.CanvasHWND,white,&rect_,0)
  draw.Rect(self.CanvasHWND,black,&rect_,1)

  rect_2 := rect.Rect(x,y,w,h)
  surface.Blit(self.CanvasHWND,my_text,&rect_2,nil)
  my_text.Free()
}

//---------WifiList---------------------------------
type BlockCbFunc func()

type WifiList struct{
  UI.Page
  WifiPassword string
  Connecting  bool
  Scanning    bool

  ShowingMessageBox bool
  MsgBox   *WifiListMessageBox
  ConnectTry  int

  BlockingUI  bool
  BlockCb     BlockCbFunc

  LastStatusMsg string
  Scroller  *UI.ListScroller
  ListFontObj  *ttf.Font

  InfoPage   *WifiInfoPage

  MyList []*NetItem 
  CurBssid  string
  CurIP     string
  CurSig    string
}

func NewWifiList() *WifiList {
  p:= &WifiList{}
  p.ListFontObj = UI.Fonts["notosanscjk15"]
  p.FootMsg = [5]string{"Nav.","Scan","Info","Back","Enter"}

  return p
}

func (self *WifiList) ShowBox(msg string ) {
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
  self.MyList = nil
  self.MyList = make([]*NetItem,0)

  start_x := 0 
  start_y := 0


  var is_active bool
  var li_idx int
  li_idx = 0
 
  self.WifiScanStarted() 
  if bssList, err := wifi.ScanManager.Scan(); err == nil {

    for _, bss := range bssList {
        is_active = false
        print(bss.SSID, bss.Signal, bss.KeyMgmt)
        ni := NewNetItem()
        ni.Parent = self
        ni.PosX   = start_x
        ni.PosY   = start_y + li_idx* NetItemDefaultHeight
        ni.Width  = UI.Width
        ni.FontObj = self.ListFontObj
        ni.Essid = bss.SSID
        ni.Bssid = bss.BSSID
        ni.Init(is_active)
        self.MyList = append(self.MyList,ni)
      
        li_idx++
    }
  }

  self.WifiScanFinished()

  self.PsIndex = 0
}

func (self *WifiList) Disconnect() {
  self.Connecting = false
  wpa_cli_disconnect := []string{"wpa_cli","disconnect",self.CurBssid}
  //out, err := execCmd(getVolumeCmd())
  execCmd( wpa_cli_disconnect )

  
}

func (self *WifiList) ShutDownConnecting() {
   
  self.Connecting= false
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

  self.UpdateNetList(-1,[]string{}, true,false)

  self.Scanning= false
  self.HideBox()

  self.BlockingUI = false
  fmt.Println("dbus says scan finished")

}

func (self *WifiList) WifiScanStarted( ) {
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
//----------------------------------------------------------------------------------

func (self *WifiList) UpdateNetList(state int,info []string ,force_check bool,firstrun bool) { //force_check default ==false, firstrun default == false 

  if force_check == true  {
    self.GenNetworkList()
    self.SaveNetworkList()
  }
  
}

func (self *WifiList) ConfigWireless(password string) {
  
  ssid := self.MyList[self.PsIndex].Essid

  self.ShowBox("Connecting...")

  if conn, err := wifi.ConnectManager.Connect(ssid, password, gotime.Second * 20); err == nil {
	fmt.Println("Connected", conn.NetInterface, conn.SSID, conn.IP4.String(), conn.IP6.String())
    self.CurBssid = self.MyList[self.PsIndex].Bssid
    self.MyList[self.PsIndex].Password = password
    self.CurIP    = conn.IP4.String();
   
  } else {
	  fmt.Println(err)
  }


  //self.UpdateStatus()


}

func (self *WifiList) GetWirelessIP() string {

  cli := fmt.Sprintf( "ip -4 addr show %s | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}'",sysgo.WifiDev)
  out := UI.System(cli)
  
  return out
  
}

func (self *WifiList) GetWirelessEncrypt(network_id int) []map[string]string {
  return nil

}

func (self *WifiList) ScrollUp() {
  if len(self.MyList) == 0 {
    return
  }

  self.PsIndex -= 1
  if self.PsIndex < 0 {
    self.PsIndex=0
  }

  cur_ni := self.MyList[self.PsIndex]//*NetItem
  if cur_ni.PosY < 0 {
    for i:=0;i<len(self.MyList);i++ {
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
  if cur_ni.PosY + cur_ni.Height > self.Height {
    for i:=0;i<len(self.MyList);i++ {
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
  password_inputed := strings.Join(APIOBJ.PasswordPage.Textarea.MyWords,"")
  if self.Screen.DBusManager.IsWifiConnectedNow() == false {
    self.ConfigWireless(password_inputed)
  }
}

func (self *WifiList) OnReturnBackCb() {

}

func (self *WifiList) KeyDown( ev *event.Event  ) {
  if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {

      self.ShutDownConnecting()
      self.ShowBox("ShutDownConnecting...")
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
    self.CurBssid = self.MyList[self.PsIndex].Bssid

    if self.MyList[self.PsIndex].IsActive == true {
      var ip string
      ip = self.GetWirelessIP()
      self.ShowBox(ip)
    }else {
      self.Screen.PushCurPage()
      self.Screen.SetCurPage(APIOBJ.PasswordPage)

      thepass := self.MyList[self.PsIndex].Password

      fmt.Println("APIOBJ.PasswordPage.SetPassword ", thepass,len(thepass))
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


      self.Screen.PushPage(self.InfoPage)
      self.Screen.Draw()
      self.Screen.SwapAndShow()
  }

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
  msgbox.Init(" ",UI.Fonts["veramono12"],nil)
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

    self.UpdateNetList(-1,[]string{}, true,true) // self.UpdateNetList(force_check=True,firstrun=True)

    self.Scroller = UI.NewListScroller()
    self.Scroller.Parent = self
    self.Scroller.PosX = 2
    self.Scroller.PosY = 2
    self.Scroller.Init()

    self.InfoPage = NewWifiInfoPage()
    self.InfoPage.Screen = self.Screen
    self.InfoPage.Name = "Wifi info"

    self.InfoPage.Init()

}

func (self *WifiList) Draw() {
  self.ClearCanvas()

  if len(self.MyList) == 0 {
    return
  }

  self.Ps.Draw()

  for _,v := range self.MyList {
    v.Draw()
  }  

  self.Scroller.UpdateSize( len(self.MyList)*NetItemDefaultHeight, self.PsIndex*NetItemDefaultHeight)
  self.Scroller.Draw()

}


