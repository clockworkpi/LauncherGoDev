package Wifi
//wifi_list.py

import (
  "fmt"
  "strconv"
  "strings"
  gotime "time"
  
  "github.com/godbus/dbus"
  
  "github.com/veandco/go-sdl2/ttf"
  
  "github.com/cuu/gogame/surface"
  "github.com/cuu/gogame/font"
  "github.com/cuu/gogame/color"
  "github.com/cuu/gogame/event"
  "github.com/cuu/gogame/time"
  "github.com/cuu/gogame/rect"
  "github.com/cuu/gogame/draw"
	"github.com/cuu/LauncherGoDev/sysgo/UI"
  "github.com/cuu/LauncherGoDev/sysgo/DBUS"
  
  
  "github.com/cuu/LauncherGoDev/sysgo/wicd/misc"
  
)

type WifiDisconnectConfirmPage struct {
    UI.ConfirmPage
    Parent *WifiInfoPage
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
  
  if ev.Data["key"] == UI.CurKeys["B"] {
    self.SnapMsg("Disconnecting...")
    self.Screen.Draw()
    self.Screen.SwapAndShow()
    
    DBUS.DBusHandler.Daemon.Method("Disconnect")
    
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
  
  MyList []UI.ListItemInterface
  
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
  self.MyList = make([]UI.ListItemInterface,0)
  
  cur_network_id = -2
  
  if self.NetworkId != -1 {
    self.AList["ip"]["value"] = "Not Connected"
    self.Wireless.Get( self.Wireless.Method("GetIwconfig"), &iwconfig)
    self.Wireless.Get( self.Wireless.Method("GetCurrentNetworkID",iwconfig), &cur_network_id)
    if cur_network_id == self.NetworkId {
      var ip string 
      self.Wireless.Get( self.Wireless.Method("GetWirelessIP",""), &ip)
      
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
  self.Wireless.Get( self.Wireless.Method("GetWirelessIP",""), &ip)
  
  if cur_network_id == self.NetworkId  && len(ip) >  1 {
    self.Screen.PushPage(self.DisconnectConfirmPage)
    self.Screen.Draw()
    self.Screen.SwapAndShow()
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
  self.Wireless.Get( self.Wireless.Method("GetWirelessIP",""), &ip)
  
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
  
  var num_of_networks int
  var cur_signal_strength int
  var cur_network_id int // -1 or 0-n
  var iwconfig string  
  var wireless_ip string 
  
  var is_active bool
  
  self.Wireless.Get(self.Wireless.Method("GetNumberOfNetworks"),&num_of_networks)
  
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
func (self *WifiList) WifiDbusScanFinishedSig(body []interface{}) {
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

func (self *WifiList) WifiDbusScanStarted(body []interface{} ) {
  if self.Screen.CurrentPage != self {
    return
  }
  
  self.Scanning = true
  self.ShowBox("Wifi scanning...")
  self.BlockingUI = true
  fmt.Println("dbus says start scan")
}


func (self *WifiList) DbusDaemonStatusChangedSig(body []interface{}) {
	var state int
	var info []dbus.Variant

	err := dbus.Store(body,&state,&info)

	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(state," ", info)
	}
  
  var info_str []string 
  for _,v := range info {
    info_str = append(info_str, v.String())
  } 
  
  self.UpdateNetList(state,info_str,false,false)
  if len(info_str) > 0 {
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
}

func (self *WifiList) DbusConnectResultsSent(body []interface{}) {
  var ret_val string
	err := dbus.Store(body,&ret_val)

	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(ret_val)
	}  
  
  self.Connecting = false
  self.BlockingUI = false
  if self.BlockCb != nil {
    self.BlockCb()
    self.BlockCb = nil
  }
  
  self.Screen.FootBar.ResetNavText()
}

//----------------------------------------------------------------------------------

func (self *WifiList) UpdateNetList(state int,info []string ,force_check bool,firstrun bool) { //force_check default ==false, firstrun default == false 
  if self.Daemon == nil {
    return
  }
    
  type status struct {
    State int
    Trash  []string
  }
  
  var mystatus status
  
  if state == -1 {
    
    self.Daemon.Get(self.Daemon.Method("GetConnectionStatus"),&mystatus)
    fmt.Println("state ",mystatus.State)
    fmt.Println("Trash ",mystatus.Trash)
    
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
    
    status_msg = fmt.Sprintf("%s: %s", essid,stat)
    
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
    
    self.Wireless.Get( self.Wireless.Method("GetWirelessIP",""), &ip)
    
    if self.CheckForWireless(iwconfig,ip,"") == true { // self.CheckForWireless(iwconfig,self._Wireless.GetWirelessIP(''),None)
      return true
    }else {
      fmt.Println("not Connected")
      return true
    }
  }
  
  return true
  
}

//set_status == "" not used
func (self *WifiList) CheckForWireless(iwconfig string, wireless_ip string , set_status string ) bool { 
  if len(wireless_ip) == 0 {
    return false 
  }
  
  var network string 
  self.Wireless.Get(self.Wireless.Method("GetCurrentNetwork",iwconfig),&network)
  
  if len(network) == 0 {
    return false
  }
  
  var sig_display_type int
  var strength int
  
  strength = -1
  
  self.Daemon.Get(self.Daemon.Method("GetSignalDisplayType"),&sig_display_type)
  
  if sig_display_type == 0 {
    self.Wireless.Get(self.Wireless.Method("GetCurrentSignalStrength",iwconfig),&strength)
  }else {
    self.Wireless.Get(self.Wireless.Method("GetCurrentDBMStrength",iwconfig),&strength)
  }
  
  if strength == -1 {
    return false 
  }
  
  var strength_str string
  
  self.Daemon.Get(self.Daemon.Method("FormatSignalForPrinting",strength),&strength_str)
  
  fmt.Printf("Connected to %s at %s (IP: %s)\n",network,strength_str,wireless_ip)
  
  return true
  
}

func (self *WifiList) ConfigWireless(password string) {
  netid := self.PsIndex
  fmt.Println(netid, " ", password)
  
  /*
  self.Wireless.Method("SetWirelessProperty",netid,"dhcphostname","GameShell")
  self.Wireless.Method("SetWirelessProperty",netid,"ip","None")
  self.Wireless.Method("SetWirelessProperty",netid,"dns_domain","None")
  self.Wireless.Method("SetWirelessProperty",netid,"gateway","None")
  self.Wireless.Method("SetWirelessProperty",netid,"use_global_dns",0)
  self.Wireless.Method("SetWirelessProperty",netid,"netmask","None")
  self.Wireless.Method("SetWirelessProperty",netid,"usedhcphostname",0) ## set 1 to use hostname above
  self.Wireless.Method("SetWirelessProperty",netid,"bitrate","auto")
  self.Wireless.Method("SetWirelessProperty",netid,"allow_lower_bitrates",0)
  self.Wireless.Method("SetWirelessProperty",netid,"dns3","None")
  self.Wireless.Method("SetWirelessProperty",netid,"dns2","None")
  self.Wireless.Method("SetWirelessProperty",netid,"dns1","None")
  self.Wireless.Method("SetWirelessProperty",netid,"use_settings_globally",0)
  self.Wireless.Method("SetWirelessProperty",netid,"use_static_dns",0)
  self.Wireless.Method("SetWirelessProperty",netid,"search_domain","None")
  */
  
  self.Wireless.Method("SetWirelessProperty",netid,"enctype","wpa-psk")
  self.Wireless.Method("SetWirelessProperty",netid,"apsk",password)
  self.Wireless.Method("SetWirelessProperty",netid,"automatic",1)
  
  self.ShowBox("Connecting...")
  
  self.MyList[netid].Connect()
  
  fmt.Println("after connect")
  self.UpdateStatus()
  

}

func (self *WifiList) GetWirelessEncrypt(network_id int) []map[string]string {
  var results []map[string]string 
  
  activeID := -1
  var enc_type string 
  
  for i,v := range self.EncMethods {
    enc_type = ""
    self.Wireless.Get(self.Wireless.Method("GetWirelessProperty",network_id,"encryption_method"),&enc_type)
    enc_type = strings.ToLower(enc_type)
    if enc_type != "" && v.Type == enc_type {
      activeID = i
      break
    }
  }
  
  if activeID == -1 {
    return results
  }
  
  required_fields := self.EncMethods[activeID].Required
  for _,field := range required_fields {
    if len(field) != 2 {
      continue
    }
    text := strings.Replace(strings.ToLower(field[1])," ","_",-1)
    
    var value string
    
    self.Wireless.Get(self.Wireless.Method("GetWirelessProperty",network_id,field[0]),&value)
    
    kv_map := make(map[string]string)
    kv_map[text] = value
    
    results = append(results,kv_map)
    
/*
         """
        [{'preshared_key': 'blah blah blah',},]

        or nothing 
        [{'identity': "",},{'password': "",},]

        """
 */    
    
  }

  optional_fields := self.EncMethods[activeID].Optional
  for _,field := range optional_fields {
    if len(field) != 2 {
      continue
    }
    text := strings.Replace(strings.ToLower(field[1])," ","_",-1)
    
    var value string
    
    self.Wireless.Get(self.Wireless.Method("GetWirelessProperty",network_id,field[0]),&value)
    
    kv_map := make(map[string]string)
    kv_map[text] = value
    
    results = append(results,kv_map)
  }  
  
  return results
  
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

func (self *WifiList) OnReturnBackCb() {
  password_inputed := strings.Join(APIOBJ.PasswordPage.Textarea.MyWords,"")
  if self.Screen.DBusManager.IsWifiConnectedNow() == false {
    self.ConfigWireless(password_inputed)
  }
}

func (self *WifiList) KeyDown( ev *event.Event  ) {
  if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
    if self.Wireless != nil {
      var wireless_connecting bool
      self.Wireless.Get(self.Wireless.Method("CheckIfWirelessConnecting"),&wireless_connecting)
      
      if wireless_connecting == true {
        self.ShutDownConnecting()
        self.ShowBox("ShutDownConnecting...")
        self.BlockingUI = true
        self.BlockCb = self.AbortedAndReturnToUpLevel
        
      }else {
        self.AbortedAndReturnToUpLevel()
      }
    }else {
      self.HideBox()
      self.ReturnToUpLevelPage()
      self.Screen.Draw()
      self.Screen.SwapAndShow()
    }
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
    
    wicd_wireless_encrypt_pwd := self.GetWirelessEncrypt(self.PsIndex) 
    fmt.Println("wicd_wireless_encrypt_pwd  ", wicd_wireless_encrypt_pwd)
    
    if self.MyList[self.PsIndex].IsActive == true {
      var ip string
      self.Wireless.Get(self.Wireless.Method("GetWirelessIP",""),&ip)
      self.ShowBox(ip)
    }else {
      self.Screen.PushCurPage()
      self.Screen.SetCurPage(APIOBJ.PasswordPage)
      
      thepass := ""
      for _,v := range wicd_wireless_encrypt_pwd { //[]map[string]string
        if _, ok := v["preshared_key"]; ok {
          if len(v["preshared_key"]) > 0 {
            thepass = v["preshared_key"]
          }
        }
      }
      
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
    
    self.InfoPage.NetworkId = self.PsIndex
    self.InfoPage.Wireless  = self.Wireless
    self.InfoPage.Daemon    = self.Daemon
    
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
  
  self.EncMethods = misc.LoadEncryptionMethods(false) //# load predefined templates from /etc/wicd/...
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


