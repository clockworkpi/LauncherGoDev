package Pico8

import (
	/*
		"github.com/veandco/go-sdl2/ttf"

		"github.com/cuu/gogame/surface"
		"github.com/cuu/gogame/event"
		"github.com/cuu/gogame/rect"
		"github.com/cuu/gogame/color"
	*/
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	//"github.com/clockworkpi/LauncherGoDev/sysgo/DBUS"
)

/******************************************************************************/
type GamePlugin struct {
	UI.Plugin
	MainPage *GamePage
}

func (self *GamePlugin) Init(main_screen *UI.MainScreen) {
	self.MainPage = NewGamePage()
	self.MainPage.SetScreen(main_screen)
	self.MainPage.SetName("Pico8")
	self.MainPage.Init()
}

func (self *GamePlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.MainPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ GamePlugin
