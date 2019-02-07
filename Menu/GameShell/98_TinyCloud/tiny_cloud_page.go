package TinyCloud

import (
  "fmt"
  
	"github.com/veandco/go-sdl2/ttf"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
  
	"github.com/cuu/gogame/color"
	
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"

)

type TinyCloudLabel struct {
  Key string
  Lable string
  Font  *ttf.Font
  Color  *color.Color

}

type TinyCloudPage struct {
  UI.Page
  ListFontObj *ttf.Font  
  URLColor *color.Color
  TextColor *color.Color
  Labels map[string]UI.LabelInterface
  Icons  map[string]UI.IconItemInterface

  Coords map[string]UI.Coord
  IP string 
  
  PngSize map[string]UI.Plane
}

func NewTinyCloudPage() *TinyCloudPage {
  p := &TinyCloudPage{}
  p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]
	
	p.FootMsg = [5]string{"Nav.","","","Back",""}  
  
  p.Coords = make(map[string]UI.Coord)
  p.PngSize = make(map[string]UI.Plane)
  
  p.URLColor = UI.MySkinManager.GiveColor("URL")
  p.TextColor = UI.MySkinManager.GiveColor("Text")
  p.ListFontObj = UI.MyLangManager.TrFont("varela13")
  
  p.Labels = make(map[string]UI.LabelInterface)
  
  p.Icons  = make(map[string]UI.IconItemInterface)

  return p
}

func (self *TinyCloudPage) OnLoadCb() {
  self.PosY = 0
}

func (self *TinyCloudPage) SetCoords() {
  self.Coords["forID"] = UI.Coord{15,11}
  
  self.Coords["forKey"] = UI.Coord{71,self.Coords["forID"].Y}
  
  self.Coords["key_and_pass"] = UI.Coord{36, self.Coords["forID"].Y}
  
  self.Coords["forssh"] = UI.Coord{self.Coords["forID"].X,36}
     
  self.Coords["ssh_addr"] = UI.Coord{self.Coords["forID"].X,54}
 
  self.Coords["forwin"]  = UI.Coord{self.Coords["forID"].X,80}
  
  self.Coords["samba_games"] = UI.Coord{ self.Coords["forID"].X,97}
   
  self.Coords["samba_music"] = UI.Coord{ self.Coords["samba_games"].X,115}

  self.Coords["for_airplay"] = UI.Coord{ self.Coords["forID"].X,140}
  
  self.Coords["airplay_name"] = UI.Coord{ 68,self.Coords["for_airplay"].Y}

  self.Coords["for-usb-eth"] = UI.Coord{ self.Coords["forID"].X,163}
  
  self.Coords["usb-eth-addr"] = UI.Coord{ 112,self.Coords["for-usb-eth"].Y}
  
  self.Coords["bg"]  =  UI.Coord{  self.Width/2,self.Height/2 }
  
  self.Coords["online"] = UI.Coord{ 266, 99 }

}

func (self *TinyCloudPage) SetLabels() {
  if self.Screen.DBusManager.IsWifiConnectedNow() {
    self.IP = self.Screen.DBusManager.GetWifiIP()
    fmt.Printf("TinyCould : %s\n",self.IP)
  }else {
    self.IP = "xxx.xxx.xxx.xxx"
  }
  
  
  labels := []*TinyCloudLabel{
    &TinyCloudLabel{"forssh","For ssh and scp:",self.ListFontObj,self.TextColor},
    &TinyCloudLabel{"ssh_addr",fmt.Sprintf("ssh cpi@%s",self.IP), self.ListFontObj,self.URLColor},
    &TinyCloudLabel{"forwin", "For Windows network:",    self.ListFontObj, self.TextColor},
    &TinyCloudLabel{"samba_games", fmt.Sprintf("\\\\%s\\games", self.IP), self.ListFontObj,self.URLColor},
    &TinyCloudLabel{"samba_music", fmt.Sprintf("\\\\%s\\music" , self.IP), self.ListFontObj,self.URLColor},
    &TinyCloudLabel{"forID",      "ID:",            self.ListFontObj, self.TextColor},
    &TinyCloudLabel{"forKey",     "Key:",           self.ListFontObj, self.TextColor},
    &TinyCloudLabel{"key_and_pass", "cpi",          self.ListFontObj, self.URLColor},
    &TinyCloudLabel{"for_airplay", "Airplay:",      self.ListFontObj, self.TextColor},
    &TinyCloudLabel{"airplay_name","clockworkpi",   self.ListFontObj, self.URLColor},
    &TinyCloudLabel{"for-usb-eth","USB-Ethernet:",  self.ListFontObj, self.TextColor},
    &TinyCloudLabel{"usb-eth-addr","192.168.10.1",  self.ListFontObj, self.URLColor},
  }
  
  for _,v := range labels {
    l := UI.NewLabel()
    l.SetCanvasHWND(self.CanvasHWND)
    l.Init(v.Lable,v.Font,nil)
    l.SetColor(v.Color)
    self.Labels[v.Key]  = l 
  }
  
  self.SetCoords()

}

func (self *TinyCloudPage) Init() {
  if self.Screen == nil {
    panic("No Screen")
  }
  
  if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil{
    self.HWND = self.Screen.CanvasHWND
    self.CanvasHWND = surface.Surface(self.Screen.Width,self.Screen.Height)
  }
  
  self.PosX = self.Index*self.Screen.Width 
  self.Width = self.Screen.Width 
  self.Height = self.Screen.Height  
  
  self.PngSize["bg"] = UI.Plane{253,114}
  self.PngSize["online"] = UI.Plane{75,122}
  
  bgpng := UI.NewIconItem()
  bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("needwifi_bg")
  bgpng.MyType = UI.ICON_TYPES["STAT"]
  bgpng.Parent = self
  bgpng.Adjust(0,0,self.PngSize["bg"].W,self.PngSize["bg"].H,0)
  
  self.Icons["bg"] = bgpng
 
  onlinepng := UI.NewIconItem() 
  onlinepng.ImgSurf = UI.MyIconPool.GetImgSurf("online")
  onlinepng.MyType = UI.ICON_TYPES["STAT"]
  onlinepng.Parent = self
  onlinepng.Adjust(0,0,self.PngSize["online"].W, self.PngSize["online"].H,0)
  
  self.Icons["online"] = onlinepng

  self.SetLabels()
    
}

func (self *TinyCloudPage) KeyDown( ev *event.Event ) {
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  return
}

func (self *TinyCloudPage) Draw() {
  self.ClearCanvas()
  if self.Screen.DBusManager.IsWifiConnectedNow() {
    self.Icons["online"].NewCoord(self.Coords["online"].X, self.Coords["online"].Y)
    self.Icons["online"].Draw()    
    
    
    for k,_ := range self.Labels{
      if _ ,ok :=  self.Coords[k]; ok {
        self.Labels[k].NewCoord( self.Coords[k].X, self.Coords[k].Y)
        self.Labels[k].Draw()
      }
    }
    
    self.Labels["key_and_pass"].NewCoord( 103,self.Coords["key_and_pass"].Y)
    self.Labels["key_and_pass"].Draw()
  }else {
    self.Icons["bg"].NewCoord(self.Coords["bg"].X, self.Coords["bg"].Y)
    self.Icons["bg"].Draw()
    
    self.Labels["for-usb-eth"].NewCoord(self.Coords["for-usb-eth"].X+55, self.Coords["for-usb-eth"].Y)
    self.Labels["for-usb-eth"].Draw()
                
    self.Labels["usb-eth-addr"].NewCoord(self.Coords["usb-eth-addr"].X+55, self.Coords["usb-eth-addr"].Y)
    self.Labels["usb-eth-addr"].Draw()                
                

  
  }
  
  if self.HWND != nil {
    surface.Fill(self.HWND,UI.MySkinManager.GiveColor("white"))
    rect_ := rect.Rect(self.PosX,self.PosY,self.Width,self.Height)
    surface.Blit(self.HWND,self.CanvasHWND,&rect_,nil)
  }
}
