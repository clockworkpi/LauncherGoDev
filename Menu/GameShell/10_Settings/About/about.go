package main

import (
	"fmt"
	"strconv"
	
	"strings"
	"os/exec"
	
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
	
	"github.com/cuu/LauncherGo/sysgo/UI"

	
)

type InfoPageListItem struct{
	UI.Widget
	Labels map[string]UI.LabelInterface
	Icons  map[string]UI.IconItemInterface
	Fonts  map[string]*ttf.Font

	Parent UI.PageInterface

	Flag   string
}

func NewInfoPageListItem() *InfoPageListItem {
	i := &InfoPageListItem{}
	i.Labels = make(map[string]UI.LabelInterface)
	i.Icons = make( map[string]UI.IconItemInterface)
	i.Fonts = make(map[string]*ttf.Font)

	i.Height = 20
	i.Width  = 0
	
	return i
}

func (self *InfoPageListItem) Init(text string) {
	l := UI.NewLabel()
	l.PosX = 10
	l.SetCanvasHWND(self.Parent.GetCanvasHWND())
	l.Init(text,self.Fonts["normal"],nil)

	self.Labels["Text"] = l
	
}

func (self *InfoPageListItem) SetSmallText( text string) {
	l := UI.NewMultiLabel()
	l.SetCanvasHWND(self.Parent.GetCanvasHWND())
	l.Init(text,self.Fonts["small"],nil)

	self.Labels["Small"] = l

	_,h_ := self.Labels["Small"].Size()
	if h_>= self.Height {
		self.Height = h_ + 10
	}

}

func (self *InfoPageListItem) Draw() {
	x_,_ := self.Labels["Text"].Coord()
	self.Labels["Text"].NewCoord(x_,self.PosY)
	self.Labels["Text"].Draw()

	if _, ok := self.Labels["Small"]; ok {
		w_,_ := self.Labels["Text"].Size()
		self.Labels["Small"].NewCoord(w_+16,self.PosY)
		self.Labels["Small"].Draw()
	}
	
}


type AboutPage struct {
	UI.Page
	AList map[string]map[string]string
	ListFontObj  *ttf.Font
	Scrolled int
	BGwidth int
	BGheight int
	DrawOnce bool
	Scroller *UI.ListScroller

	MyList []*InfoPageListItem
	Icons map[string]UI.IconItemInterface
}

func NewAboutPage() *AboutPage {
	p := &AboutPage{}
	
	p.FootMsg = [5]string{"Nav","","","Back",""}

	p.AList = make(map[string]map[string]string)

	p.BGwidth = 320
	p.BGheight = 300
	p.DrawOnce = false

	p.MyList = make([]*InfoPageListItem,0)

	p.ListFontObj = UI.Fonts["varela13"]

	p.Index = 0
	
  p.Icons = make(map[string]UI.IconItemInterface)
	return p
	
}

func (self *AboutPage) Uname() {
	out := make(map[string]string)

	out["key"] = "uname"
	out["label"] = "Kernel:"

	out_bytes, err := exec.Command("uname","-srmo").Output()
	if err != nil {
		fmt.Println(err)
		out["value"] = ""
	}
	
	out_str := strings.Trim(string(out_bytes), "\t\n")
	
	out["value"]= out_str

	self.AList["uname"] = out
}


func (self *AboutPage) CpuMhz() {
  
  lines, err := UI.ReadLines("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
  UI.ShowErr(err)

  mhz ,err := strconv.ParseInt(lines[0], 10, 64)
	UI.ShowErr(err)
  mhz_float := float64(mhz)/1000.0
  
  out := make(map[string]string)
  out["key"] = "cpuscalemhz"
  out["label"]="CPU Mhz:"
  out["value"] =  strconv.FormatFloat(mhz_float, 'f', 2, 64)
  
  self.AList["cpuscalemhz"] = out
  
}

func (self *AboutPage) CpuInfo() {
  last_processor := 0
  
}

func (self *AboutPage) MemInfo() {
  lines, err := UI.ReadLines("/proc/meminfo")
  UI.ShowErr(err)	
  
  for _,line := range lines {
    if strings.HasPrefix(line,"MemTotal") {
      parts := strings.Split(line,":")
      kb := strings.Replace(parts[1],"KB","",-1)
      kb = strings.TrimSpace(kb)
      kb_int,_ := strconv.ParseInt(kb,10,0)
      
      kb_float := float64(kb_int)/1000.0
      memory := make(map[string]string)
      memory["key"] = "memory" 
      memory["label"] = "Memory:"
      memory["value"] = strconv.FormatFloat(kb_float,'f',2,64) + " MB"
      self.AList["memory"] = memory 
      break
    }
  }  
}

func (self *AboutPage) GenList() {
	self.MyList = nil
	self.MyList = make([]*InfoPageListItem,0)
  
  start_x  := 0
  start_y  := 10
  last_height := 0 
  
  for _,u := range ( []string{"processor","armcores","cpuscalemhz","features","memory","uname"} ) {
    if val, ok := self.AList[u]; ok {
      
			li := NewInfoPageListItem()
			li.Parent = self
			li.PosX = start_x
			li.PosY = start_y + last_height
			li.Width = UI.Width
			li.Fonts["normal"] = self.ListFontObj
			li.Fonts["small"] = UI.Fonts["varela12"]      
      
      if self.AList[u]["label"] != "" {
        li.Init( self.AList[u]["label"]  )
      }else {
        li.Init( self.AList[u]["key"])
      }
      
			li.Flag = val["key"]
      li.SetSmallText(val["value"])
      last_height += li.Height
			
      self.MyList = append(self.MyList,li)

    }
  } 
	
}


func (self *AboutPage) Init() {

	if self.Screen != nil {
		if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
			self.HWND = self.Screen.CanvasHWND
			self.CanvasHWND = surface.Surface(self.Screen.Width,self.BGheight)
		}

		self.PosX = self.Index * self.Screen.Width
		self.Width = self.Screen.Width
		self.Height = self.Screen.Height

    bgpng := UI.NewIconItem()
    bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("about_bg")
    bgpng.MyType = UI.ICON_TYPES["STAT"]
    bgpng.Parent = self
    bgpng.Adjust(0,0,self.BGwidth,self.BGheight,0)
    
    self.Icons["bg"] = bgpng


		self.CpuInfo()
    self.MemInfo()
    self.CpuMhz()
		self.Uname()
    
		self.GenList()

		self.Scroller = UI.NewListScroller()
		
		self.Scroller.Parent = self
		self.Scroller.PosX = self.Width - 10
		self.Scroller.PosY = 2
		self.Scroller.Init()
		self.Scroller.SetCanvasHWND(self.HWND)
		
	}
}

func (self *AboutPage) ScrollDown() {
	dis := 10
	if UI.Abs(self.Scrolled) < ( self.BGheight - self.Height)/2 + 50 {
		self.PosY -= dis
		self.Scrolled -= dis
	}
}

func (self *AboutPage) ScrollUp() {
	dis := 10
	if self.PosY < 0 {
		self.PosY += dis
		self.Scrolled += dis
	}
}

func (self *AboutPage) OnLoadCb() {
	self.Scrolled = 0
	self.PosY     = 0
	self.DrawOnce = false
}

func (self *AboutPage) OnReturnBackCb() {
	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()
}


func (self *AboutPage) KeyDown( ev *event.Event) {
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
	
}


func (self *AboutPage) Draw() {
	if self.DrawOnce == false {

		self.ClearCanvas()

    self.Icons["bg"].NewCoord(self.Width/2, self.Height/2 + (self.BGheight - UI.Height)/2 + self.Screen.TitleBar.Height)
    self.Icons["bg"].Draw()
    
		for _,v := range self.MyList {
			v.Draw()
		}
		
		self.DrawOnce = true
	}

	if self.HWND != nil {
		surface.Fill(self.HWND, &color.Color{255,255,255,255})

		rect_ := rect.Rect(self.PosX,self.PosY,self.Width,self.Height)
		surface.Blit(self.HWND,self.CanvasHWND,&rect_, nil)
    
		self.Scroller.UpdateSize(self.BGheight,UI.Abs(self.Scrolled)*3)
		self.Scroller.Draw()
		
	}
}


/******************************************************************************/
type AboutPlugin struct {
	UI.Plugin
	Page UI.PageInterface
}


func (self *AboutPlugin) Init( main_screen *UI.MainScreen ) {
	self.Page = NewAboutPage()
	self.Page.SetScreen( main_screen)
	self.Page.SetName("About")
	self.Page.Init()
}

func (self *AboutPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
		main_screen.PushPage(self.Page)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ AboutPlugin








