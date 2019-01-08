package GateWay

import (
  "fmt"
  //"io/ioutil"
  //"path/filepath"
  "strings"
  "strconv"
  "runtime"
  "github.com/veandco/go-sdl2/ttf"
  //"github.com/mitchellh/go-homedir"

  "github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/time"

  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"

)

type ListPageSelector struct {
  UI.InfoPageSelector
}

func NewListPageSelector() *ListPageSelector {
  
  p := &ListPageSelector{}
  
  p.Width = UI.Width
  p.BackgroundColor = &color.Color{131,199,219,255} //SkinManager().GiveColor('Front')
  
  return p

}


func (self *ListPageSelector) Draw() {
  
  idx := self.Parent.GetPsIndex()
  mylist := self.Parent.GetMyList()
 
  if idx < len(mylist) {
    x,y := mylist[idx].Coord()
    _,h := mylist[idx].Size()
    
    self.PosX = x+2
    self.PosY = y+1
    self.Height = h-3
    
    canvas_ := self.Parent.GetCanvasHWND()
    rect_   := rect.Rect(self.PosX,self.PosY,self.Width-4, self.Height)
    
    draw.AARoundRect(canvas_,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  } 
}

type PageListItem struct {
  UI.InfoPageListItem
  
  Active bool
  Value string
}


func NewPageListItem() *PageListItem {
  
  p := &PageListItem{}
  p.Height = UI.DefaultInfoPageListItemHeight
  p.ReadOnly = false
	p.Labels = make(map[string]UI.LabelInterface)
	p.Icons  = make( map[string]UI.IconItemInterface)
	p.Fonts  = make(map[string]*ttf.Font)  

  return p
}


func (self *PageListItem) Draw() {
    
  x,_ := self.Labels["Text"].Coord()
  w,h := self.Labels["Text"].Size()
  
  self.Labels["Text"].NewCoord( x, self.PosY + (self.Height - h)/2 )
  
  
  if self.Active == true {
    self.Parent.(*GateWayPage).Icons["done"].NewCoord(self.Parent.(*GateWayPage).Width-30,self.PosY+5)
    self.Parent.(*GateWayPage).Icons["done"].Draw()
  }
  
  self.Labels["Text"].SetBold(self.Active)
  self.Labels["Text"].Draw()
  
  
  
  if _, ok := self.Labels["Small"]; ok {
    x,_ = self.Labels["Small"].Coord()
    w,h = self.Labels["Small"].Size()
    
    self.Labels["Small"].NewCoord( self.Width - w - 10 , self.PosY + (self.Height - h)/2 )
    self.Labels["Small"].Draw()
    
  }
  
  canvas_ := self.Parent.GetCanvasHWND()
  draw.Line(canvas_, &color.Color{169,169,169,255}, 
    self.PosX, self.PosY+self.Height -1,
    self.PosX + self.Width, self.PosY+self.Height -1 ,1)
  
}

type GateWayPage struct {
  UI.Page
  
  ListFont *ttf.Font
  
  BGwidth  int 
  BGheight  int
  
  DrawOnce bool
  
  Scroller *UI.ListScroller

  Icons map[string]UI.IconItemInterface
}


func NewGateWayPage() *GateWayPage {
  p := &GateWayPage{}
  
  p.ListFont = UI.Fonts["notosanscjk15"]
  p.FootMsg = [5]string{"Nav","","Clear All","Back","Select"}
  
  p.BGwidth = UI.Width 
  p.BGheight = UI.Height - 24 - 20
  
  p.Icons = make(map[string]UI.IconItemInterface)
  
  return p
}

func (self *GateWayPage) GenList() {

  self.MyList = nil
  
  start_x := 0 
  start_y := 0 
  last_height := 0
  
  var drivers  = [][2]string{[2]string{"usb0","USB Ethernet"},
                             [2]string{"wlan0","Wi-Fi"}}
  

  for _,u := range drivers {
    li := NewPageListItem()
    li.Parent = self
    li.PosX   = start_x
    li.PosY   = start_y + last_height
    li.Width  = UI.Width
    li.Fonts["normal"] = self.ListFont
    li.Active = false
    li.Value  = u[0]
    li.Init(u[1])
    last_height += li.Height
    self.MyList = append(self.MyList,li)
  }

}

func (self *GateWayPage) Init() {

  if self.Screen != nil {
    if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
      self.HWND = self.Screen.CanvasHWND
      self.CanvasHWND = surface.Surface( self.Screen.Width,self.Screen.Height )
    }
  }
  
  self.PosX = self.Index*self.Screen.Width 
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height
  
  done := UI.NewIconItem()
  done.ImgSurf = UI.MyIconPool.GetImgSurf("done")
  done.MyType = UI.ICON_TYPES["STAT"]
  done.Parent = self
  
  self.Icons["done"] = done
  
  ps := NewListPageSelector()
  ps.Parent = self
  
  self.Ps = ps
  self.PsIndex = 0
  
  self.GenList()
  
  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = self.Width - 10
  self.Scroller.PosY = 2
  self.Scroller.Init()
  self.Scroller.SetCanvasHWND(self.HWND)
  
}

func (self *GateWayPage) Click() {
  
  if len(self.MyList) == 0 {
    return
  }
  
  if self.PsIndex >= len(self.MyList) {
    self.PsIndex = len(self.MyList) -1 
  }
  
  cur_li := self.MyList[self.PsIndex].(*PageListItem)
  if cur_li.Active == true {
    out := UI.System( "sudo ip route | grep default | cut -d \" \" -f3"  )
    if len(out) > 7 {
      self.Screen.MsgBox.SetText(strings.Trim(out,"\r\n "))
      self.Screen.MsgBox.Draw()
      self.Screen.SwapAndShow()      
    }
    return
  }
  
  if strings.Contains(runtime.GOARCH,"arm") == true {
    for i,_ := range self.MyList {
      self.MyList[i].(*PageListItem).Active = false
    }
    cur_li.Active = self.ApplyGateWay(cur_li.Value)
    self.Screen.MsgBox.SetText("Applying")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow()
    
    time.BlockDelay(1000)
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }else {
    self.Screen.MsgBox.SetText("Do it in GameShell")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow()
  }

}

func (self *GateWayPage) ClearAllGateways() {
  self.Screen.MsgBox.SetText("Cleaning up")
  self.Screen.MsgBox.Draw()
  self.Screen.SwapAndShow()
  
  UI.System("sudo ip route del 0/0")
  time.BlockDelay(800)
  
  for i,_ := range self.MyList {
    self.MyList[i].(*PageListItem).Active = false
  }
   
  self.Screen.Draw()
  self.Screen.SwapAndShow()
  
}

func (self *GateWayPage) ApplyGateWay( gateway string ) bool {
  UI.System("sudo ip route del 0/0")
  if gateway == "usb0" {
    out := UI.System("sudo ifconfig usb0 | grep inet | tr -s \" \"| cut -d \" \" -f3")
    if len(out) > 7 {
      if strings.Contains(out,"error") == false {
        parts := strings.Split(out,".")
        if len(parts) == 4 { // IPv4
          tmp,err := strconv.Atoi(parts[3])
          if err == nil {
            tmp = tmp +1 
            if tmp > 255 {
              tmp = 255
            }
            
            parts[3] = strconv.Itoa(tmp)
            ipaddress := strings.Join(parts,".")
            UI.System( fmt.Sprintf("sudo route add default gw %s",ipaddress)  )
            return true
          }
        }
      }
    }
  }else { // wlan0
    if self.Screen.DBusManager.IsWifiConnectedNow() == true {
      UI.System("sudo dhclient wlan0")
      return true
    }else {
      self.Screen.MsgBox.SetText("Wi-Fi is not connected")
      self.Screen.MsgBox.Draw()
      self.Screen.SwapAndShow()    
      return false
    }
    
  }
  
  return false
}

func (self *GateWayPage) OnLoadCb() {
  
  self.PosY = 0
  self.DrawOnce = false
  
  thedrv := ""
  if strings.Contains(runtime.GOARCH,"arm") == true {
    out := UI.System("sudo ip route | grep default")
    if len(out) > 7 {
      if strings.Contains(out,"usb0") {
        thedrv = "usb0"
      }else if strings.Contains(out,"wlan0") {
        thedrv = "wlan0"
      }
    }
  }
    
  for i, _ := range self.MyList {
    self.MyList[i].(*PageListItem).Active = false
  }
  
  if thedrv != "" {
    for i, v := range self.MyList {
      if strings.Contains( v.(*PageListItem).Value, thedrv) {
        self.MyList[i].(*PageListItem).Active = true
      }
    }  
  }
}

func (self *GateWayPage) KeyDown(ev *event.Event ) {
  
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  
  if ev.Data["Key"] == UI.CurKeys["B"] {
    self.Click()
  }
  
  if ev.Data["Key"] == UI.CurKeys["Y"] {
    self.ClearAllGateways()
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
  
}

func (self *GateWayPage) Draw() {
  
  self.ClearCanvas()
  if len(self.MyList) == 0 {
    return
  }
  
  if len(self.MyList) * UI.DefaultInfoPageListItemHeight > self.Height {
    
    self.Ps.(*ListPageSelector).Width  = self.Width - 11
    self.Ps.Draw()
    
    for _,v := range self.MyList {
      if v.(*PageListItem).PosY > self.Height + self.Height/2 {
        break
      }
      
      if v.(*PageListItem).PosY < 0 {
        continue
      }
      
      v.Draw()
    
    }
    
    self.Scroller.UpdateSize( len(self.MyList)*UI.DefaultInfoPageListItemHeight,
                            self.PsIndex*UI.DefaultInfoPageListItemHeight)
    self.Scroller.Draw()
    
  }else {
    self.Ps.(*ListPageSelector).Width  = self.Width
    self.Ps.Draw()
    for _,v := range self.MyList {
      if v.(*PageListItem).PosY > self.Height + self.Height/2 {
        break
      }
      
      if v.(*PageListItem).PosY < 0 {
        continue
      }
      
      v.Draw()
    
    }
  }
  
  if self.HWND != nil {
    surface.Fill(self.HWND, &color.Color{255,255,255,255})
    rect_ := rect.Rect(self.PosX,self.PosY,self.Width,self.Height)
    surface.Blit(self.HWND,self.CanvasHWND,&rect_,nil)
  }
}

