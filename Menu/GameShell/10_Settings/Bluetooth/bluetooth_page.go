package Bluetooth

import (
  "fmt"
  "os"
  "log"
  "strings"
  //"errors"
	gotime "time"
  "github.com/fatih/structs"
  
  "github.com/veandco/go-sdl2/ttf"

  "github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/surface"
  "github.com/cuu/gogame/rect"
  "github.com/cuu/gogame/event"
  "github.com/cuu/gogame/time"
  "github.com/cuu/gogame/color"
  "github.com/cuu/gogame/font"

  //"github.com/godbus/dbus"
  bleapi "github.com/muka/go-bluetooth/api"
  //"github.com/muka/go-bluetooth/bluez"
 // "github.com/muka/go-bluetooth/bluez/profile"
  "github.com/muka/go-bluetooth/bluez/profile/device"
  "github.com/muka/go-bluetooth/bluez/profile/adapter"

	logrus "github.com/sirupsen/logrus"
  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

func showDeviceInfo(dev *device.Device1) {
  if dev == nil {
    return
  }
  props, err := dev.GetProperties()
  if err != nil {
	  fmt.Printf("%s: Failed to get properties: %s\n", dev.Path, err.Error())
	  return
  }
  fmt.Printf("name=%s addr=%s rssi=%d\n", props.Name, props.Address, props.RSSI)
}


type BleForgetConfirmPage struct {
  
  UI.ConfirmPage


}

func NewBleForgetConfirmPage() *BleForgetConfirmPage {
  p := &BleForgetConfirmPage{}
  
  p.ListFont = UI.Fonts["veramono20"]
  p.FootMsg = [5]string{"Nav","","","Cancel","Yes"}
  
  p.ConfirmText = "Confirm Forget?"
  p.ConfirmPage.ConfirmText = p.ConfirmText
  
  return p
}

func (self *BleForgetConfirmPage) KeyDown(ev *event.Event) {

  if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  } 

  if ev.Data["Key"] == UI.CurKeys["B"] {
    self.SnapMsg("Deleting")
    self.Screen.Draw()
    self.Screen.SwapAndShow()
    
    
    time.BlockDelay(400)
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()    

  }
}

func (self *BleForgetConfirmPage) Draw() {
  self.DrawBG()
  for _,v := range self.MyList{
    v.Draw()
  }
}


type BleInfoPageSelector struct {
  UI.InfoPageSelector
  
}

func NewBleInfoPageSelector() *BleInfoPageSelector{
  p := &BleInfoPageSelector{}
  p.Width = UI.Width
  p.BackgroundColor = &color.Color{131,199,219,255} //SkinManager().GiveColor('Front')
  
  return p 
}

func (self *BleInfoPageSelector) Draw() {
  
  idx := self.Parent.GetPsIndex()
  mylist := self.Parent.GetMyList()

  if idx < len(mylist) {
    _,y := mylist[idx].Coord()
    _,h := mylist[idx].Size()
    
    x := self.PosX+2
    self.PosY = y+1
    self.Height = h-3
    
    canvas_ := self.Parent.GetCanvasHWND()
    rect_   := rect.Rect(x,self.PosY,self.Width-4, self.Height)

    draw.AARoundRect(canvas_,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  }
}


type BleInfoPage struct {
  UI.Page
  
  ListFontObj *ttf.Font
  ListSmFontObj *ttf.Font
  ListSm2FontObj *ttf.Font
    
  AList map[string]interface{}
  
  Scroller *UI.ListScroller
  ConfirmPage1 *BleForgetConfirmPage
  MyDevice  *device.Device1 // from NetItem-> from BluetoothPage
  Props    *device.Device1Properties
  Path      string
}


func NewBleInfoPage() *BleInfoPage {
  p :=&BleInfoPage{}
  
  p.FootMsg = [5]string{"Nav","Disconnect","Forget","Back","" }  
  
  p.ListFontObj = UI.Fonts["varela15"]
  p.ListSmFontObj = UI.Fonts["varela12"]
  p.ListSm2FontObj = UI.Fonts["varela11"]
  
  return p
}

func (self *BleInfoPage) Init() {
  
  if self.Screen != nil {
    if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
      self.CanvasHWND = self.Screen.CanvasHWND
    }
  }
     
  self.PosX = self.Index*self.Screen.Width 
  self.Width = self.Screen.Width //  equals to screen width
  self.Height = self.Screen.Height
        
  ps := NewBleInfoPageSelector()
  ps.Parent = self
  self.Ps = ps
  self.PsIndex = 0
  
  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = 2
  self.Scroller.PosY = 2
  self.Scroller.Init()
        
  self.ConfirmPage1 = NewBleForgetConfirmPage()
  self.ConfirmPage1.Screen = self.Screen
  self.ConfirmPage1.Name   = "Confirm Forget"
  self.ConfirmPage1.Init()   

}

func (self *BleInfoPage) GenList() {
  
  self.AList = structs.Map(self.Props) //map[string]interface{}
  
  self.MyList = nil
  
  self.PsIndex = 0
  
  start_x := 0 
  start_y := 0
  
  i := 0
  skip_arrays := []string{"ManufacturerData","AdvertisingFlags","ServiceData"}
  
  for k,v := range self.AList {
    
    skip2 := false
    for _,u := range skip_arrays {
      if strings.HasPrefix(k,u) {
        skip2 = true
        break
      }
    }
    if skip2 {
      continue
    }
    
    li := UI.NewInfoPageListItem()
    li.Parent = self
    li.PosX   = start_x
    li.PosY   = start_y +i*NetItemDefaultHeight
    li.Width  = UI.Width
    
    li.Fonts["normal"] = self.ListFontObj
    if k =="UUIDs" {
      li.Fonts["small"] = self.ListSm2FontObj
    }else{
      li.Fonts["small"] = self.ListSmFontObj
    }
    
    li.Init(k)
    li.Flag = k
    
    sm_text := ""
    if k == "UUIDs" {
      if len(v.([]string))> 1 {
        sm_text = v.([]string)[0]
      }else{
        sm_text = "<empty>"
      }
    }else {
      sm_text = fmt.Sprintf("%v",v)
    }
    
    if len(sm_text) > 20 {
      sm_text = sm_text[:20]
    }
    li.SetSmallText(sm_text)
    li.PosX = 2
    
    self.MyList = append(self.MyList,li)
    i+=1
  }

}

func (self *BleInfoPage) ScrollUp() {
  if len(self.MyList) == 0 {
    return
  }
  
  self.PsIndex -= 1
  
  if self.PsIndex < 0 {
    self.PsIndex = 0
  }
  
  cur_li := self.MyList[self.PsIndex]
  
  x,y := cur_li.Coord()
  
  if y < 0 {
    for i,v := range self.MyList {
      x,y = v.Coord()
      _,h := v.Size()
      self.MyList[i].NewCoord(x,y+h)
    }
  }

}

func (self *BleInfoPage) ScrollDown() {
  if len(self.MyList) == 0 {
    return
  }

  self.PsIndex += 1
  
  if self.PsIndex >= len(self.MyList) {
    self.PsIndex = len(self.MyList)-1
  }
  
  cur_li := self.MyList[self.PsIndex]
  
  x,y := cur_li.Coord()
  _,h := cur_li.Size()
  
  if y + h >  self.Height {
    for i,v := range self.MyList {
      x,y = v.Coord()
      _,h = v.Size()
      self.MyList[i].NewCoord(x,y-h)
    }
  } 

}


func (self *BleInfoPage) TryToForget() {
  //muka Adapter1 RemoveDevice  Path
  
  adapter,err := bleapi.GetAdapter(adapterID)
  if err == nil {
    self.Screen.MsgBox.SetText("Forgeting")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow()    
  
  
    err = adapter.RemoveDevice(self.MyDevice.Path())
    if err != nil {
      fmt.Println("BleInfoPage TryToForget: ",err)
    }
    
    time.BlockDelay(400)
    
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()  
    
  }else {
  
    fmt.Println("BleInfoPage TryToForget GetAdapter: ",err)
  }
  
}

func (self *BleInfoPage) TryToDisconnect() {
  
  is_connected,_ := self.MyDevice.GetConnected();

  if is_connected {
  
    self.Screen.FootBar.UpdateNavText("Disconnecting")
    self.Screen.MsgBox.SetText("Disconnecting")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow()
    
    self.MyDevice.Disconnect()
    
    time.BlockDelay(350)
    
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  
    self.Screen.FootBar.ResetNavText()
  }

}

func (self *BleInfoPage) Click() {
  if self.PsIndex >= len(self.MyList) {
    return
  }
  
  
  cur_li := self.MyList[self.PsIndex]
  
  fmt.Println(cur_li.(*UI.InfoPageListItem).Flag)
  

}

func (self *BleInfoPage) OnLoadCb() {
  if self.Props.Connected == true {
    self.FootMsg[1] = "Disconnect"
  }else {
    self.FootMsg[1] = ""
  }
  
  self.GenList()
}


func (self *BleInfoPage) KeyDown(ev *event.Event) {

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
  if ev.Data["Key"] == UI.CurKeys["Enter"]{
    self.Click()
  }
  if ev.Data["Key"] == UI.CurKeys["X"] {
    self.TryToDisconnect()
  }
  if ev.Data["Key"] == UI.CurKeys["Y"] {
    self.TryToForget()
  }
  
}

func (self *BleInfoPage) Draw() {
  if len(self.MyList) == 0 {
    return
  }

  self.ClearCanvas()
  
  if len(self.MyList) * UI.DefaultInfoPageListItemHeight > self.Height {
    self.Ps.(*BleInfoPageSelector).Width = self.Width - 10
    self.Ps.(*BleInfoPageSelector).PosX = 9
    self.Ps.Draw()
    
    for _,v := range self.MyList {
      v.Draw()
    }
    
    self.Scroller.UpdateSize(len(self.MyList)*UI.DefaultInfoPageListItemHeight, 
                            self.PsIndex*UI.DefaultInfoPageListItemHeight)
    self.Scroller.Draw()
    
  }else {
    self.Ps.(*BleInfoPageSelector).Width = self.Width 
    self.Ps.Draw()
    for _,v := range self.MyList {
      v.Draw()
    }
  }
}

type BleListMessageBox struct {
  UI.Label
  Parent UI.PageInterface
}

func NewBleListMessageBox() *BleListMessageBox {
  p := &BleListMessageBox{}
  p.Color = &color.Color{83,83,83,255}
  
  return p
}

func (self *BleListMessageBox) Draw() {
  
	my_text := font.Render(self.FontObj,self.Text, true, self.Color, nil)
  w := surface.GetWidth(my_text)
  h := surface.GetHeight(my_text)
  
  pw,ph := self.Parent.Size()
  
  x := (pw-w)/2
  y := (ph-h)/2
  
  padding := 10
  
  rect_ := rect.Rect(x-padding,y-padding,w+padding*2,h+padding*2)
  
  draw.Rect(self.CanvasHWND,&color.Color{255,255,255,255}, &rect_,0)
  draw.Rect(self.CanvasHWND,&color.Color{0,  0,  0,  255}, &rect_,1)
  
  rect2_ := rect.Rect(x,y,w,h)
  surface.Blit(self.CanvasHWND,my_text,&rect2_,nil)
  my_text.Free()

}

type BluetoothPage struct{
  UI.Page
   
  Devices []*device.Device1 
  
  BlePassword string 
  Connecting bool
  Scanning  bool
  
  
  ListFontObj *ttf.Font
  Scroller *UI.ListScroller
  InfoPage *BleInfoPage
  
  PrevState  int
  
  ShowingMessageBox bool
  MsgBox  *BleListMessageBox
  ConnectTry int
  
  //BlockCb ??
  
  LastStatusMsg string
  ADAPTER_DEV string // == adapterID
  
  Offline  bool
  
  Leader *BluetoothPlugin
}

func NewBluetoothPage() *BluetoothPage {
  p := &BluetoothPage{}
  
  p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10
	p.Align = UI.ALIGN["SLeft"]
  
  p.ADAPTER_DEV = adapterID
  
  p.FootMsg = [5]string { "Nav","Scan","Info","Back","TryConnect" }
  
  p.ListFontObj = UI.Fonts["notosanscjk15"]
  
  return p 
}

func (self *BluetoothPage) ShowBox(msg string) {
  self.MsgBox.Text = msg
  self.ShowingMessageBox = true
  self.Screen.Draw()
  self.MsgBox.Draw()
  self.Screen.SwapAndShow()  
}

func (self *BluetoothPage) HideBox() {
  self.Draw()
  self.ShowingMessageBox = false
  self.Screen.SwapAndShow()
}

func (self *BluetoothPage) Init() {
  self.PosX = self.Index * self.Screen.Width
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height  
    
  self.CanvasHWND = self.Screen.CanvasHWND

  ps := NewBleInfoPageSelector()
  ps.Parent = self
  ps.Width = UI.Width - 12
        
  self.Ps = ps
  self.PsIndex = 0
        
  msgbox := NewBleListMessageBox()
  msgbox.CanvasHWND = self.CanvasHWND
  msgbox.Init(" ",UI.Fonts["veramono12"],nil)
  msgbox.Parent = self
        
  self.MsgBox = msgbox     

  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = 2
  self.Scroller.PosY = 2
  self.Scroller.Init()
	
  self.InfoPage = NewBleInfoPage()
  self.InfoPage.Screen = self.Screen
  self.InfoPage.Name   = "BluetoothInfo"
  self.InfoPage.Init()
  
}


func (self *BluetoothPage) AbortedAndReturnToUpLevel() {

  self.HideBox()
  self.Screen.FootBar.ResetNavText()
  self.ReturnToUpLevelPage()
  self.Screen.Draw()
  self.Screen.SwapAndShow()

}

func (self *BluetoothPage) TryConnect() {
  
  if self.PsIndex >= len(self.MyList) {
    return
  }
  
  cur_li := self.MyList[self.PsIndex]
  
  if cur_li.(*NetItem).Props.Connected {
  
    return
  }
  
  self.Screen.FootBar.UpdateNavText("Connecting")
  self.ShowBox("Connecting")
  
  self.Leader.PairPage.DevObj = cur_li.(*NetItem).Device
  
  err := cur_li.(*NetItem).Device.Pair()
  if err != nil {
    fmt.Println(err)
    err_msg := ""
    s := err.Error()
    err_msg = "Pair error"
    if strings.Contains(s,"ConnectionAttemptFailed") {
      err_msg = "Page Timeout"
    }
    if strings.Contains(s,"NoReply") {
      err_msg = "NoReply,Cancelling"
      dev1 := cur_li.(*NetItem).Device
      dev1.CancelPairing()
      
    }
    if strings.Contains(s,"Exists") {
      err_msg = "Already Exists"
      adapter,err := bleapi.GetAdapter(adapterID)
      if err == nil {
        err = adapter.RemoveDevice(cur_li.(*NetItem).Device.Path())
        if err != nil {
          fmt.Println(err)
        }
      }else {
        fmt.Println(err)
      }
    }
    
    self.Leader.PairPage.PairErrorCb( err_msg )
    self.Leader.PairPage.DevObj= nil
    
  }else{
    self.Leader.PairPage.PairOKCb()
    dev1 := cur_li.(*NetItem).Device
	err = dev1.SetTrusted(true)
	if err != nil {
	  fmt.Println(err)
    }
    cur_li.(*NetItem).Device.Connect()
  }
  
  self.HideBox()
  self.Screen.FootBar.ResetNavText()
}

//GetDevices returns a list of bluetooth discovered Devices
func (self *BluetoothPage) GetDevices() ([]*device.Device1, error) {

  adapter,err := bleapi.GetAdapter(adapterID)
  if err != nil {
    return nil,err
  }
  
  list, err := adapter.GetDevices()
  return list,err
}

func (self *BluetoothPage) RefreshDevices() {
  
  // sync the cached devices 
  self.Devices = self.Devices[:0]
  
  devices, err := self.GetDevices()
	if err != nil {
		panic(err)
		os.Exit(1)
	}
  
  self.Devices  = devices
  
}


func (self *BluetoothPage) GenNetworkList() {
  self.MyList = nil
  
  start_x := 0 
  start_y := 0 
	
  for i, v := range self.Devices { // v == bleapi.Device
  
  	props, err := v.GetProperties()
    if err != nil {
      log.Fatalf("%s: Failed to get properties: %s", v.Path, err.Error())
      return
    }

    ni := NewNetItem()
    ni.Parent = self
    
    ni.PosX = start_x
    ni.PosY  = start_y + i*NetItemDefaultHeight
    ni.Width  = UI.Width
    ni.FontObj = self.ListFontObj
    ni.Props  = props
    ni.Parent = self
    ni.Device = v
    if props.Name != "" {
      ni.Init(props.Name)
    }else {
      ni.Init(props.Address)
    }
    
    self.MyList = append(self.MyList,ni)
    
  }
  
  self.PsIndex = 0
}


func (self *BluetoothPage) Rescan() {

	if self.Scanning == true {
		self.ShowBox("Bluetooth scanning")
		self.Screen.FootBar.UpdateNavText("Scanning")
	}
	
  a, err := adapter.GetAdapter(adapterID)
  if err != nil {
    fmt.Println(err)
    return
  }

  discovery, cancel, err := bleapi.Discover(a, nil)
  if err != nil {
    fmt.Println(err)
  }

  defer cancel()

  wait := make(chan error)
	self.Scanning  = true
	self.ShowBox("Bluetooth scanning")
  self.Screen.FootBar.UpdateNavText("Scanning")
	
  go func() {
    for dev := range discovery {
      if dev == nil {
        return
      }
      wait <- nil
    }
  }()

  go func() {
    sleep := 5
    gotime.Sleep(gotime.Duration(sleep) * gotime.Second)
    logrus.Debugf("Discovery timeout exceeded (%ds)", sleep)
    wait <- nil
  }()

  err = <-wait
  if err != nil {
    fmt.Println(err)
  }
	self.Scanning = false
	self.HideBox()
	self.Screen.FootBar.ResetNavText()
}

func (self *BluetoothPage) OnLoadCb() {
  self.Offline = false
  
  if self.Screen.TitleBar.InAirPlaneMode == false {
    out := UI.System("hcitool dev | grep hci0 |cut -f3")
    if len(out) < 17 {
      self.Offline = true
      fmt.Println("Bluetooth OnLoadCb ,can not find hci0 alive,try to reboot")
    }else {
			self.Rescan()
      self.RefreshDevices()
      self.GenNetworkList()
    }
  }else {
    self.Offline = true
  }
  
}

func (self *BluetoothPage) ScrollUp() {

  if len(self.MyList) == 0 {
    return
  }

  self.PsIndex -= 1
  if self.PsIndex < 0 {
    self.PsIndex=0
  }
  
  cur_ni := self.MyList[self.PsIndex]//*NetItem
  if cur_ni.(*NetItem).PosY < 0 {
    for i:=0;i<len(self.MyList);i++ {
      self.MyList[i].(*NetItem).PosY += self.MyList[i].(*NetItem).Height
    }
  }
}

func (self *BluetoothPage) ScrollDown() {
  if len(self.MyList) == 0 {
    return 
  }
  
  self.PsIndex += 1
  if self.PsIndex >= len(self.MyList) {
    self.PsIndex = len(self.MyList) - 1
  }
  
  cur_ni := self.MyList[self.PsIndex]
  if cur_ni.(*NetItem).PosY + cur_ni.(*NetItem).Height > self.Height {
    for i:=0;i<len(self.MyList);i++ {
      self.MyList[i].(*NetItem).PosY -= self.MyList[i].(*NetItem).Height
    }
  }
}

func (self *BluetoothPage) KeyDown(ev *event.Event) {

  if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
    if self.Offline == true {
      self.AbortedAndReturnToUpLevel()
      return
    }
    
    a, nil := bleapi.GetAdapter(adapterID)
    err := a.StopDiscovery()
    if err != nil {
      fmt.Println(err)
    }
    
    self.HideBox()
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()  
    
    self.Screen.FootBar.ResetNavText()
  }
  
  if ev.Data["Key"]  == UI.CurKeys["Up"] {
  
    self.ScrollUp()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  if ev.Data["Key"]  == UI.CurKeys["Down"] {
  
    self.ScrollDown()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  if ev.Data["Key"]  == UI.CurKeys["X"] {
    if self.Offline == false{
      self.Rescan()
    }
  }
  
  if ev.Data["Key"]  == UI.CurKeys["Y"] {
    if len(self.MyList) == 0 {
      return
    }
    if self.Offline == true {
      return
    }
    self.InfoPage.Props    = self.MyList[self.PsIndex].(*NetItem).Props
    self.InfoPage.Path     = self.MyList[self.PsIndex].(*NetItem).Path
    self.InfoPage.MyDevice = self.MyList[self.PsIndex].(*NetItem).Device
    
    self.Screen.PushPage(self.InfoPage)
    self.Screen.Draw()
    self.Screen.SwapAndShow()        
  }
  
  if ev.Data["Key"] == UI.CurKeys["B"] {
    if self.Offline == false {
      self.TryConnect()
    }
  }
}

func (self *BluetoothPage) Draw() {
  
  self.ClearCanvas()
  
  if len(self.MyList) == 0 {
    return
  }
  
  
  if len(self.MyList) * NetItemDefaultHeight > self.Height {
    self.Ps.(*BleInfoPageSelector).Width  = self.Width - 11
    self.Ps.Draw()
    
    for _,v := range self.MyList {
      v.Draw()
    }
    
    self.Scroller.UpdateSize(len(self.MyList)*NetItemDefaultHeight,self.PsIndex*NetItemDefaultHeight)
    self.Scroller.Draw()
    
  }else {
    self.Ps.(*BleInfoPageSelector).Width = self.Width
    self.Ps.Draw()

    for _,v := range self.MyList {
      v.Draw()
    }    
  
  }
 
}
