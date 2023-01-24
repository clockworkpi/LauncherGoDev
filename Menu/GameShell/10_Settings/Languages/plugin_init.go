package Languages

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
type LanguagesPlugin struct {
	UI.Plugin
	LanguagesPage *LanguagesPage
}

func (self *LanguagesPlugin) Init(main_screen *UI.MainScreen) {
	self.LanguagesPage = NewLanguagesPage()
	self.LanguagesPage.SetScreen(main_screen)
	self.LanguagesPage.SetName("Languages")
	self.LanguagesPage.Init()
}

func (self *LanguagesPlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.LanguagesPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ LanguagesPlugin
