package TinyCloud

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
type TinyCloudPlugin struct {
	UI.Plugin
	MainPage *TinyCloudPage
}

func (self *TinyCloudPlugin) Init(main_screen *UI.MainScreen) {
	self.MainPage = NewTinyCloudPage()
	self.MainPage.SetScreen(main_screen)
	self.MainPage.SetName("Tiny cloud")
	self.MainPage.Init()
}

func (self *TinyCloudPlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.MainPage)
		main_screen.Refresh()
	}
}

var APIOBJ TinyCloudPlugin
