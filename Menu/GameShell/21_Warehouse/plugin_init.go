package Warehouse

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
type WareHousePlugin struct {
	UI.Plugin
	MainPage *WareHouse
}

var aria2dl_folder = "%s/aria2download%s"

func (self *WareHousePlugin) Init(main_screen *UI.MainScreen) {
	self.MainPage = NewWareHouse()
	self.MainPage.SetScreen(main_screen)
	self.MainPage.SetName("Warehouse")
	self.MainPage.Init()
}

func (self *WareHousePlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.MainPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ WareHousePlugin
