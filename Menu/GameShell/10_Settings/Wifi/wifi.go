package main
//wifi_list.py

import (
  gotime "time"
  "github.com/cuu/gogame/surface"
  "github.com/cuu/gogame/font"
  "github.com/cuu/gogame/color"
  "github.com/cuu/gogame/event"
  "github.com/cuu/gogame/time"
	"github.com/cuu/LauncherGo/sysgo/UI"
  "github.com/cuu/LauncherGo/sysgo/DBUS"
  
  "github.com/cuu/LaucherGo/sysgo/wicd/misc"
  
)

type WifiDisconnectConfirmPage struct {
    UI.ConfirmPage
    Parent *WifiInfoPage
}

func NewWifiDisconnectConfirmPage() *WifiDisconnectConfirmPage {
  p := &WifiDisconnectConfirmPage{}  
  p.ConfirmText ="Confirm Disconnect?"
  return p
}

func (self *WifiDisconnectConfirmPage) KeyDown(ev *event.Event ) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  
  if ev.Data["key"] == UI.CurKeys["B"] {
    self.SnapMsg("Disconnecting...")
    self.Screen.Draw()
    self.Screen.SwapAndShow()
    
    self.Parent.Daemon.Disconnect()
    
    time.BlockDelay(400)
    
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
    
  }
}

type WifiInfoPage struct {
  UI.Page
  ListFontObj  *ttf.Font
  
  Wireless *DBUS.DbusInterface
  Daemon   *DBUS.DbusInterface
  AList map[string]map[string]string
  NetworkId  int
  
  MyList []*UI.ListItemInterface
  
  DisconnectConfirmPage *WifiDisconnectConfirmPage //child page 
}

func NewWifiInfoPage() *WifiInfoPage {
  p := &WifiInfoPage{}
  p.FootMsg = [5]string{"Nav","Disconnect","","Back",""}
  
  p.ListFontObj = UI.Fonts["varela15"]
  
  p.AList = make(map[string]map[string]string)

  p.NetworkId = -1
  return p
  
}

func (self *WifiInfoPage) GenList() {
  var iwconfig string
  var cur_network_id int
  self.MyList = nil
  self.MyList = make([]*UI.ListItemInterface,0)
  
  cur_network_id = -2
  
  if self.NetworkId != -1 {
    self.AList["ip"]["value"] = "Not Connected"
    self.Wireless.Get( self.Wireless.Method("GetIwconfig"), &iwconfig)
    self.Wireless.Get( self.Wireless.Method("GetCurrentNetworkID",iwconfig), &cur_network_id)
    if cur_network_id == self.NetworkId {
      var ip string 
      self.Wireless.Get( self.Wireless.Method("GetWirelessIP",''), &ip)
      
      if len(ip) > 0 {
        self.AList["ip"]["value"]=ip
      }
    }
    var bssid string
    self.Wireless.Get( self.Wireless.Method("GetWirelessProperty",self.NetworkId,"bssid"),&bssid)
    
    self.AList["bssid"]["value"] = bssid
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
  var iwconfig string
  var cur_network_id int 
  var ip string 
  self.Wireless.Get( self.Wireless.Method("GetIwconfig"), &iwconfig)
  self.Wireless.Get( self.Wireless.Method("GetCurrentNetworkID",iwconfig), &cur_network_id)  
  self.Wireless.Get( self.Wireless.Method("GetWirelessIP",''), &ip)
  
  if cur_network_id == self.NetworkId  && len(ip) >  1 {
    self.Screen.PushPage(self.DisconnectConfirmPage)
    self.Screen.Draw()
    self.SwapAndShow()
  }else {
    return
  }
}

func (self *WifiInfoPage) OnLoadCb() {
  var iwconfig string
  var cur_network_id int 
  var ip string 
  self.Wireless.Get( self.Wireless.Method("GetIwconfig"), &iwconfig)
  self.Wireless.Get( self.Wireless.Method("GetCurrentNetworkID",iwconfig), &cur_network_id)  
  self.Wireless.Get( self.Wireless.Method("GetWirelessIP",''), &ip)
  
  if cur_network_id == self.NetworkId && len(ip) > 1 {
    self.FootMsg[1]="Disconnect"
  }else {
    self.FootMsg[1] = ""
  }
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
  if idx < len(self.Parent.WirelessList) {
    x := self.Parent.WirelessList[idx].PosX + 11
    y := self.Parent.WirelessList[idx].PosY + 1
    h := self.Parent.WirelessList[idx].Height - 3
    
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
  
  return p
}


func (self *WifiListMessageBox) Draw() {
  my_text := font.Render(self.FontObj,self.Text,true,self.Color)
  
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
  
}

//---------WifiList---------------------------------
type BlockCbFunc func()

type WifiList struct{
  UI.Page
  Wireless *DBUS.DbusInterface
  Daemon   *DBUS.DbusInterface
  WifiPassword string
  Connecting  bool
  Scanning    bool
  
  PrevWicdState  int
  
  ShowingMessageBox bool
  MsgBox   *WifiListMessageBox
  ConnectTry  int
  
  BlockingUI  bool
  BlockCb     BlockCbFunc
  
  LastStatusMsg string
  EncMethods []*misc.CurType
  Scroller  *UI.ListScroller
  ListFontObj  *ttf.Font
  
  InfoPage   *WifiInfoPage
  
  MyList []*NetItem 
  
}

func NewWifiList() *WifiList {
  p:= &WifiList{}
  p.PrevWicdState = -1
  
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
  
  var num_of_networks int
  var cur_signal_strength int
  var cur_network_id int // -1 or 0-n
  var iwconfig string  
  var wireless_ip string 
  
  var is_active bool
  
  self.Wireless.Method("GetNumberOfNetworks"),&num_of_networks)
  
  for network_id:=0;network_id< num_of_networks;network_id++ {
    is_active = false
    
    self.Wireless.Get(self.Wireless.Method("GetCurrentSignalStrength",""), &cur_signal_strength)
    self.Wireless.Get(self.Wireless.Method("GetIwconfig"),&iwconfig)
    self.Wireless.Get(self.Wireless.Method("GetCurrentNetworkID",iwconfig),&cur_network_id)
    
    if cur_signal_strength != 0 && cur_network_id == network_id {
      self.Wireless.Get(self.Wireless.Method("GetWirelessIP",""),&wireless_ip)
      if wireless_ip != "" {
        is_active = true
      }
    }
    
    ni := NewNetItem()
    ni.Parent = self
    ni.PosX   = start_x
    ni.PosY   = start_y + network_id* NetItemDefaultHeight
    ni.Width  = UI.Width
    ni.FontObj = self.ListFontObj
    ni.Init(network_id, is_active)
    self.MyList = append(self.MyList,ni)
    
  }
  self.PsIndex = 0
}

func (self *WifiList) Disconnect() {
  self.Connecting = false
  self.Daemon.Method("Disconnect")
}

func (self *WifiList) ShutDownConnecting() {
  fmt.Println("Shutdownconnecting...", self.ConnectTry)
  self.Daemon.Method("CancelConnect")
  self.Daemon.Method("SetForcedDisconnect",true)
  self.Connecting= false
}

func (self *WifiList) Rescan(sync bool) { // sync default should be false
  fmt.Println("start Rescan")
  if self.Wireless!= nil {
    self.Wireless.Method("Scan",sync)
  }
}

// dbus signal functions
func (self *WifiList) DbusScanFinishedSig() {
}

func (self *WifiList) DbusScanStarted() {

}


func (self *WifiList) UpdateNetList(state int,info []string ,force_check bool,firstrun bool) { //force_check default ==false, firstrun default == false 
  if self.Daemon == nil {
    return
  }
  
  var state_ int
  var trash []string
  
  if state == -1 {
    self.Daemon.Get(self.Daemon.Method("GetConnectionStatus"),&state_,&trash)
    fmt.Println("state ",state_)
    fmt.Println("Trash ",trash)
  }
  
  if force_check == true || self.PrevWicdState != state {
    self.GenNetworkList()
  }
  
  if len(info) > 0 {
    if len(info) > 3 {
      _id,_ := strconv.Atoi(info[3])
      if _id < len(self.MyList) {
        self.MyList[_id].UpdateStrenLabel(info[2])
      }
    }
  }
  
  self.PrevWicdState = state

}

func (self *WifiList) SetConnectingStatus(fast bool) bool { // default fast == false
  
  var wireless_connecting bool
  var iwconfig string 
  
  var essid string
  var stat  string
  var status_msg string 
  
  self.Wireless.Get(self.Wireless.Method("CheckIfWirelessConnecting"),&wireless_connecting)
  
  
  if wireless_connecting == true {
    if fast == false {
      self.Wireless.Get(self.Wireless.Method("GetIwconfig"),&iwconfig)
    }else {
      iwconfig=""
    }
    
    self.Wireless.Get(self.Wireless.Method("GetCurrentNetwork",iwconfig),&essid)
    
    err := self.Wireless.Get(self.Wireless.Method("CheckWirelessConnectingMessage"),&stat) // wicd will return False or stat message,False is a boolean,stat is string
    if err != nil {
      return false
    }
    
    status_msg := fmt.Sprintf("%s: %s", essid,stat)
    
    if self.LastStatusMsg != status_msg {
      fmt.Printf("%s: %s\n",essid,stat)
      self.LastStatusMsg = status_msg
      
      self.ShowBox(self.LastStatusMsg)
      
      self.Screen.FootBar.UpdateNavText(self.LastStatusMsg)
      UI.SwapAndShow()
      
    }
    
    return true
    
  }else {
    self.Connecting=false
    return self.Connecting
  }
  
  return false
}

func (self *WifiList) UpdateStatus() bool {
  fmt.Println("UpdateStatus")
  var wireless_connecting bool
  var fast bool 
  
  self.Wireless.Get(self.Wireless.Method("CheckIfWirelessConnecting"),&wireless_connecting)
  
  self.Daemon.Get(self.Daemon.Method("NeedsExternalCalls"),&fast)
  
  fast = !fast
  
  self.Connecting = wireless_connecting
  
  if self.Connecting  == true {
    go func() {
      for {
        gotime.Sleep(250 * gotime.Millisecond)
        ret := self.SetConnectingStatus(fast)
        if ret == false {
          break
        }
      }
    }()
  }else {
    
    var iwconfig string 
    var ip string 
    if fast == false {
      self.Wireless.Get(self.Wireless.Method("GetIwconfig"),&iwconfig)
    }else {
      iwconfig = ""
    }
    
    self.Wireless.Get( self.Wireless.Method("GetWirelessIP",''), &ip)
    
    if self.CheckForWireless(iwconfig,ip,"") == true { // self.CheckForWireless(iwconfig,self._Wireless.GetWirelessIP(''),None)
      return true
    }else {
      fmt.Println("not Connected")
      return true
    }
  }
  
  return true
  
}

func (self *WifiList) DbusDaemonStatusChangedSig(state int,info []string) {
  
}

func (self *WifiList) DbusConnectResultsSent( result string) {
  
}

//set_status == "" not used
func (self *WifiList) CheckForWireless(iwconfig string, wireless_ip string , set_status string ) { 
  
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
  msgbox.Init(" ",UI.Fonts["veramono12"])
  msgbox.Parent = self
  
  self.MsgBox = msgbox
  
  self.EncMethods = misc.LoadEncryptionMethods() //# load predefined templates from /etc/wicd/...
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
  
  self.UpdateNetList(true,true) // self.UpdateNetList(force_check=True,firstrun=True)
  
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


