package main

import (
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
	
}

func NewAboutPage() *AboutPage {
	p := &HelloWorldPage{}
	
	p.FootMsg = [5]string{"Nav.","","","Back",""}

	p.AList = make(map[string]map[string]string)

	p.BGwidth = 320
	p.BGheight = 300
	p.DrawOnce = false

	p.MyList = make([]*InfoPageListItem,0)

	p.ListFontObj = UI.Fonts["varela13"]

	p.Index = 0
	
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
	
}

func (self *AboutPage) CpuInfo() {
	
}

func (self *AboutPage) MemInfo() {
	
}

func (self *AboutPage) GenList() {

	
}

