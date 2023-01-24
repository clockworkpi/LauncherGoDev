package ButtonsLayout

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
type ButtonsLayoutPlugin struct {
	UI.Plugin
	ButtonsLayoutPage *ButtonsLayoutPage
}

func (self *ButtonsLayoutPlugin) Init(main_screen *UI.MainScreen) {
	self.ButtonsLayoutPage = NewButtonsLayoutPage()
	self.ButtonsLayoutPage.SetScreen(main_screen)
	self.ButtonsLayoutPage.SetName("Buttons Layout")
	self.ButtonsLayoutPage.Init()
}

func (self *ButtonsLayoutPlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.ButtonsLayoutPage)
		main_screen.Refresh()
	}
}

var APIOBJ ButtonsLayoutPlugin
