package Pico8

import (
	"fmt"

	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/color"

	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

type GamePage struct {
	UI.Page
	ListFontObj *ttf.Font
	URLColor    *color.Color
	TextColor   *color.Color
	Labels      map[string]UI.LabelInterface
	Icons       map[string]UI.IconItemInterface

	
	MsgBox *UI.MessageBox
	NotFound bool

	Pkg *UI.CommercialSoftwarePackage
}

func NewGamePage() *GamePage {
	p := &GamePage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav.", "", "", "Back", ""}


	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")
	p.ListFontObj = UI.MyLangManager.TrFont("varela18")

	//p.Labels = make(map[string]UI.LabelInterface)

	//p.Icons = make(map[string]UI.IconItemInterface)
	p.NotFound = true
	
	p.Pkg = UI.NewCommercialSoftwarePackage("/home/cpi/games/PICO-8/pico-8/pico8_dyn","/home/cpi/launchergo/Menu/GameShell/50_Pico8/")

	return p
}

func (self *GamePage) OnLoadCb() {
	self.PosY = 0
	if self.Pkg.IsInstalled() {
		self.Pkg.RunSetup()	
		fmt.Println("Run pico8")
		self.MsgBox.SetText("Running Pico8")
		self.Screen.RunEXE(self.Pkg.GetRunScript())

	}else{
		self.MsgBox.SetText("Please purchase the PICO-8 and copy it to the \"~/games/PICO-8\"")
	}

	fmt.Println("GamePage OnLoadCb")
}


func (self *GamePage) Init() {
	if self.Screen == nil {
		panic("No Screen")
	}

	if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
		self.HWND = self.Screen.CanvasHWND
		self.CanvasHWND = surface.Surface(self.Screen.Width, self.Screen.Height)
	}

	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height
	
	self.MsgBox = UI.NewMessageBox()
	self.MsgBox.Parent = self.CanvasHWND
	self.MsgBox.Init("Please purchase the PICO-8 and copy it to the \"~/games/PICO-8\"",self.ListFontObj,nil,self.Width,self.Height)

	self.Pkg.Init()	
}

func (self *GamePage) KeyDown(ev *event.Event) {
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Refresh()
	}
	return
}

func (self *GamePage) Draw() {
	self.ClearCanvas()
	
	self.MsgBox.Draw()

	if self.HWND != nil {
		surface.Fill(self.HWND, UI.MySkinManager.GiveColor("white"))
		rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
		surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
	}
}
