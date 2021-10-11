package Airplane

import (
	//"fmt"
	//"io/ioutil"
	//"path/filepath"
	"strings"

	"github.com/veandco/go-sdl2/ttf"

	//"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/time"

	//"github.com/clockworkpi/LauncherGoDev/sysgo"

	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

type AirplanePage struct {
	UI.Page

	ListFontObj *ttf.Font

	BGwidth  int
	BGheight int

	Scrolled int
	Scroller *UI.ListScroller

	airwire_y    int //0
	dialog_index int //0

	Icons map[string]UI.IconItemInterface
}

func NewAirplanePage() *AirplanePage {
	p := &AirplanePage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.ListFontObj = UI.MyLangManager.TrFont("varela13")
	p.FootMsg = [5]string{"Nav", "Rescue", "", "Back", "Toggle"}

	p.BGwidth = UI.Width
	p.BGheight = UI.Height - 24 - 20

	p.Icons = make(map[string]UI.IconItemInterface)

	return p
}

func (self *AirplanePage) GenList() {
	self.MyList = nil

}

func (self *AirplanePage) Init() {
	if self.Screen != nil {
		if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
			self.HWND = self.Screen.CanvasHWND
			self.CanvasHWND = surface.Surface(self.Screen.Width, self.Screen.Height)
		}
	}

	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	airwire := UI.NewIconItem()
	airwire.ImgSurf = UI.MyIconPool.GetImgSurf("airwire")
	airwire.MyType = UI.ICON_TYPES["STAT"]
	airwire.Parent = self
	airwire.Adjust(0, 0, 5, 43, 0)
	self.Icons["airwire"] = airwire

	GS := UI.NewIconItem()
	GS.ImgSurf = UI.MyIconPool.GetImgSurf("GS")
	GS.MyType = UI.ICON_TYPES["STAT"]
	GS.Parent = self
	GS.Adjust(0, 0, 72, 95, 0)
	self.Icons["GS"] = GS

	DialogBoxs := UI.NewMultiIconItem()
	DialogBoxs.ImgSurf = UI.MyIconPool.GetImgSurf("DialogBoxs")
	DialogBoxs.MyType = UI.ICON_TYPES["STAT"]
	DialogBoxs.Parent = self
	DialogBoxs.IconWidth = 134
	DialogBoxs.IconHeight = 93
	DialogBoxs.Adjust(0, 0, 134, 372, 0)
	self.Icons["DialogBoxs"] = DialogBoxs

	self.GenList()

	self.Scroller = UI.NewListScroller()
	self.Scroller.Parent = self
	self.Scroller.PosX = self.Width - 10
	self.Scroller.PosY = 2
	self.Scroller.Init()
	self.Scroller.SetCanvasHWND(self.HWND)

}

func (self *AirplanePage) ScrollUp() {
	dis := 10

	if self.PosY < 0 {
		self.PosY += dis
		self.Scrolled += dis
	}
}

func (self *AirplanePage) ScrollDown() {
	dis := 10

	if UI.Abs(self.Scrolled) < (self.BGheight-self.Height)/2+0 {
		self.PosY -= dis
		self.Scrolled -= dis
	}

}

func (self *AirplanePage) ToggleModeAni() {

	out := UI.System("sudo rfkill list | grep yes | cut -d \" \" -f3")

	if strings.Contains(out, "yes") {
		data := self.EasingData(0, 43)

		for _, v := range data {
			self.airwire_y -= v
			self.dialog_index = 2
			time.BlockDelay(40)
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}

		UI.System("sudo rfkill unblock all")
		self.Screen.TitleBar.InAirPlaneMode = false

	} else {
		data := self.EasingData(0, 43)
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 { // reverse data
			data[i], data[j] = data[j], data[i]
		}
		for _, v := range data {
			self.airwire_y += v
			self.dialog_index = 3
			time.BlockDelay(40)
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}
		UI.System("sudo rfkill block all")
		self.Screen.TitleBar.InAirPlaneMode = true
	}
}

func (self *AirplanePage) ToggleMode() {

}

func (self *AirplanePage) UnBlockAll() {

	self.Screen.MsgBox.SetText("TurningOn")
	self.Screen.MsgBox.Draw()
	UI.System("sudo rfkill unblock all")
	self.Screen.TitleBar.InAirPlaneMode = false
}

func (self *AirplanePage) OnLoadCb() {
	self.Scrolled = 0
	self.PosY = 0
	//self.DrawOnce = false

	out := UI.System("sudo rfkill list | grep yes | cut -d \" \" -f3")

	if strings.Contains(out, "yes") {
		self.Screen.TitleBar.InAirPlaneMode = true
		self.airwire_y = 50 + 43
		self.dialog_index = 1
	} else {
		self.dialog_index = 0
		self.airwire_y = 50
		self.Screen.TitleBar.InAirPlaneMode = false
	}
}

func (self *AirplanePage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["B"] {
		self.ToggleModeAni()
	}

	if ev.Data["Key"] == UI.CurKeys["X"] {
		self.UnBlockAll()
		self.Screen.SwapAndShow()
		time.BlockDelay(1000)
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
}

func (self *AirplanePage) Draw() {

	self.ClearCanvas()
	self.Icons["DialogBoxs"].NewCoord(145, 23)
	self.Icons["airwire"].NewCoord(80, self.airwire_y)

	self.Icons["DialogBoxs"].SetIconIndex(self.dialog_index)

	self.Icons["DialogBoxs"].DrawTopLeft()

	self.Icons["airwire"].Draw()

	self.Icons["GS"].NewCoord(98, 118)
	self.Icons["GS"].Draw()

	if self.HWND != nil {
		surface.Fill(self.HWND, &color.Color{255, 255, 255, 255})
		rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
		surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
	}
}
