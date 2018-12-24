package PowerOptions

import (
  "fmt"
  "io/ioutil"
  //"path/filepath"
  //"strings"
  
  "github.com/veandco/go-sdl2/ttf"


  "github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/time"
  
  "github.com/clockworkpi/LauncherGoDev/sysgo"

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
    self.Parent.(*PowerOptionsPage).Icons["done"].NewCoord(self.Parent.(*PowerOptionsPage).Width-30,self.PosY+5)
    self.Parent.(*PowerOptionsPage).Icons["done"].Draw()
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


type InfoPage struct {
  UI.Page
  ListFontObj *ttf.Font
  
  Time1 int
  Time2 int
  Time3 int

  AList map[string]map[string]string
}

func NewInfoPage() *InfoPage {
  
  p := &InfoPage{}
  
  p.ListFontObj = UI.MyLangManager.TrFont("varela15")
  p.FootMsg = [5]string{"Nav","","","Back",""}
  
  p.Time1 = 40
  p.Time2 = 120 // 120 secs
  p.Time3 = 300 // 5 minutes
  
  p.AList = make(map[string]map[string]string)
  
  return p
}


func (self *InfoPage) ConvertSecToMin( secs int) string {
  sec_str := ""
  min_str := ""
  
  if secs > 60 {
    m := int(secs/60)
    s := secs % 60
    
    if m > 1 {
      min_str = fmt.Sprintf("%d "+UI.MyLangManager.Tr("minutes")+ " ", m)
    }else {
      min_str = fmt.Sprintf("%d "+UI.MyLangManager.Tr("minute")+ " ", m)
    }
    
    if s == 1 {
      sec_str = fmt.Sprintf("%d "+UI.MyLangManager.Tr("second"), s)
    }else if s > 1 {
      sec_str = fmt.Sprintf("%d "+UI.MyLangManager.Tr("seconds"), s)
    }
    
  }else if secs <= 60 && secs > 0 {
    if secs > 1 {
      sec_str = fmt.Sprintf("%d "+UI.MyLangManager.Tr("seconds"), secs)
    }else {
      sec_str = fmt.Sprintf("%d "+UI.MyLangManager.Tr("second"), secs)
    }
    
  
  }else if secs == 0 {
    sec_str = UI.MyLangManager.Tr("Never")
  }
  
  
  return min_str + sec_str

}

func (self *InfoPage) RefreshList() {
  
  self.AList["time1"]["value"] = self.ConvertSecToMin(self.Time1)
  self.AList["time2"]["value"] = self.ConvertSecToMin(self.Time2)
  self.AList["time3"]["value"] = self.ConvertSecToMin(self.Time3)
  
  i := 0 
  for _ ,v := range self.AList {
    self.MyList[i].(*PageListItem).SetSmallText( v["value"] )
    i+=1
  }

}

func (self *InfoPage) GenList() {

  time1 := make(map[string]string)
  
  time1["key"] = "time1"
  
  if self.Time1 == 0 {
    time1["value"] = UI.MyLangManager.Tr("Never")
  }else {
    time1["value"] = fmt.Sprintf("%d secs",self.Time1)
  }
  time1["label"] = "Screen dimming"
  
  time2 := make(map[string]string)
  time2["key"] = "time2"
  
  if self.Time2 == 0 {
    time2["value"] = UI.MyLangManager.Tr("Never")
  }else {
    time2["value"] = fmt.Sprintf("%d secs",self.Time2)
  }
  
  time2["label"] = "Screen OFF"
  
  time3 := make(map[string]string)
  time3["key"] = "time3"
  
  if self.Time3 == 0 {
    time3["value"] = UI.MyLangManager.Tr("Never")
  }else {
    time3["value"] = fmt.Sprintf("%d secs",self.Time3)
  }
  
  time3["label"] = "Power OFF"
  
  
  self.AList["time1"] = time1
  self.AList["time2"] = time2
  self.AList["time3"] = time3  
  
  
  self.MyList = nil
  
  start_x := 0 
  start_y := 0 
  
  
  i:=0
  for _,v := range self.AList {
    
    li := NewPageListItem()
    li.Parent = self
    li.PosX = start_x
    li.PosY = start_y + i*UI.DefaultInfoPageListItemHeight
    li.Width = UI.Width
    li.Fonts["normal"] = self.ListFontObj
    li.Fonts["small"]  = UI.MyLangManager.TrFont("varela12")
    
    if v["label"] != "" {
      li.Init( v["label"])
    }else {
      li.Init(v["key"])
    }
    
    li.Flag = v["key"]
    
    li.SetSmallText(v["value"])
    
    self.MyList = append(self.MyList,li)
  
  }
}

func (self *InfoPage) Init() {

  if self.Screen != nil {
    if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
      self.CanvasHWND = self.Screen.CanvasHWND
    }
  }
  
  self.PosX = self.Index*self.Screen.Width 
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height

  ps := NewListPageSelector()
  ps.Parent = self
  self.Ps = ps
  self.PsIndex = 0
        
  self.GenList()
}

func (self *InfoPage) Click() {
  
  if self.PsIndex >= len(self.MyList) {
    return
  }
  
  cur_li := self.MyList[self.PsIndex]
  
  fmt.Println(cur_li.(*PageListItem).Flag )

}

func (self *InfoPage) OnLoadCb() {
  
  self.RefreshList()
  
}


func (self *InfoPage) KeyDown(ev *event.Event) {


	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
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

func (self *InfoPage) Draw() {
  self.ClearCanvas()
  self.Ps.Draw()
  for _,v := range self.MyList{
    v.Draw()
  }
}


type PowerOptionsPage struct {
  UI.Page
  
  ListFont *ttf.Font
  BGwidth int
  BGheight int
  
  AList map[string]map[string]string
  DrawOnce bool
  InfoPage *InfoPage
  Scroller *UI.ListScroller
  Icons map[string]UI.IconItemInterface
}


func NewPowerOptionsPage() *PowerOptionsPage {
  
  p := &PowerOptionsPage{}
  
  p.BGwidth = UI.Width 
  p.BGheight = UI.Height - 24 - 20
  
  p.AList = make(map[string]map[string]string)
  p.Icons = make(map[string]UI.IconItemInterface)

  p.ListFont = UI.Fonts["notosanscjk15"]
  p.FootMsg =  [5]string{"Nav","","Detail","Back","Select"}
  
  return p
}

func (self *PowerOptionsPage) GenList() {
  
  self.MyList = nil
  
  start_x := 0
  start_y := 0 
  
  last_height :=0
  
  supersaving := make(map[string]string)
  supersaving["key"] = "super"
  supersaving["label"] = "Power saving"
  supersaving["value"] = "supersaving"

  powersaving := make(map[string]string)
  powersaving["key"] = "saving"
  powersaving["label"] = "Balanced"
  powersaving["value"] = "powersaving"

  performance := make(map[string]string)
  performance["key"] = "performance"
  performance["label"] = "Performance"
  performance["value"] = "performance"

  server_saving := make(map[string]string)
  server_saving["key"] = "server"
  server_saving["label"] = "Server"
  server_saving["value"] = "server"
  
  self.AList["supersaving"] = supersaving
  self.AList["powersaving"] = powersaving
  self.AList["server"]      = server_saving
  self.AList["performance"] = performance
  
  
  for _,u := range [4]string{"supersaving","powersaving","server","performance"} {
    
    v := self.AList[u]
    
    li := NewPageListItem()
    li.Parent = self
    li.PosX = start_x
    li.PosY = start_y + last_height
    li.Width  = UI.Width
    
    li.Fonts["normal"] = self.ListFont
    li.Active = false
    li.Value = v["value"]
    
    if v["label"] != "" {
      li.Init(v["label"])
    }else {
      li.Init(v["key"])
    }
    
    last_height += li.Height
    
    self.MyList = append(self.MyList,li)
  
  }
}


func (self *PowerOptionsPage) Init() {
  
  if self.Screen != nil {
    if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
      self.HWND = self.Screen.CanvasHWND
      self.CanvasHWND = surface.Surface(self.Screen.Width,self.Screen.Height)
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

  self.InfoPage = NewInfoPage()
  self.InfoPage.Screen = self.Screen
  self.InfoPage.Name   = "Power option detail"
  self.InfoPage.Init()
  
}

func (self *PowerOptionsPage) Click() {
  if len(self.MyList) == 0 {
    return
  }
  
  cur_li := self.MyList[self.PsIndex].(*PageListItem)
  
  if cur_li.Active == true {
    return
  }
  
  for i,_ := range self.MyList {
    self.MyList[i].(*PageListItem).Active = false
  }
  
  cur_li.Active = true
  
  fmt.Println(cur_li.Value)
  
  d := []byte(cur_li.Value)
  err := ioutil.WriteFile("sysgo/.powerlevel", d, 0644)
  if err != nil {
    fmt.Println(err)
  }
  
  sysgo.CurPowerLevel = cur_li.Value
  
  if sysgo.CurPowerLevel == "supersaving" {
    UI.System("sudo iw wlan0 set power_save on >/dev/null")
  }else{
    UI.System("sudo iw wlan0 set power_save off >/dev/null")
  }
  
  self.Screen.MsgBox.SetText("Applying")
  self.Screen.MsgBox.Draw()
  self.Screen.SwapAndShow()  
  
  event.Post(UI.POWEROPT,"")
  
  time.BlockDelay(1000)
  
  self.Screen.Draw()
  self.Screen.SwapAndShow()  
}

func (self *PowerOptionsPage) OnLoadCb() {
  
  self.PosY = 0
  self.DrawOnce = false
  
  for i,_ := range self.MyList{
    if self.MyList[i].(*PageListItem).Value == sysgo.CurPowerLevel {
      self.MyList[i].(*PageListItem).Active = true
    }
  }
}

func (self *PowerOptionsPage) KeyDown(ev *event.Event) {
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  
  if ev.Data["Key"] == UI.CurKeys["B"] {
    self.Click()
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
  
  if ev.Data["Key"]  == UI.CurKeys["Y"] {
    cur_li := self.MyList[self.PsIndex].(*PageListItem)
    
    time1 := sysgo.PowerLevels[cur_li.Value].Dim
    time2 := sysgo.PowerLevels[cur_li.Value].Close
    time3 := sysgo.PowerLevels[cur_li.Value].PowerOff
     
            
    self.InfoPage.Time1 = time1
    self.InfoPage.Time2 = time2
    self.InfoPage.Time3 = time3
    
            
    self.Screen.PushPage(self.InfoPage)
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
}

func (self *PowerOptionsPage) Draw() {
  self.ClearCanvas()
  if len(self.MyList) == 0 {
    return
  }
  
  if len(self.MyList) * UI.DefaultInfoPageListItemHeight > self.Height {
    
    self.Ps.(*ListPageSelector).Width = self.Width - 11
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
    self.Ps.(*ListPageSelector).Width = self.Width
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


