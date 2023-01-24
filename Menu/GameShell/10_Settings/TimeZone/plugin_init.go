package TimeZone

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
type TimeZonePlugin struct {
	UI.Plugin
	TimeZonePage *TimeZoneListPage
}

func (self *TimeZonePlugin) Init(main_screen *UI.MainScreen) {
	self.TimeZonePage = NewTimeZoneListPage()
	self.TimeZonePage.SetScreen(main_screen)
	self.TimeZonePage.SetName("Timezone Selection")
	self.TimeZonePage.Init()
}

func (self *TimeZonePlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.TimeZonePage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ TimeZonePlugin
