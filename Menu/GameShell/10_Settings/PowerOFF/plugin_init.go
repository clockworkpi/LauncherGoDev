package PowerOFF

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
type PowerOFFPlugin struct {
	UI.Plugin
	PowerOFFPage *PowerOFFConfirmPage
}

func (self *PowerOFFPlugin) Init(main_screen *UI.MainScreen) {
	self.PowerOFFPage = NewPowerOFFConfirmPage()
	self.PowerOFFPage.SetScreen(main_screen)
	self.PowerOFFPage.SetName("PowerOFF")
	self.PowerOFFPage.Init()
}

func (self *PowerOFFPlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.PowerOFFPage)
		main_screen.Refresh()
	}
}

var APIOBJ PowerOFFPlugin
