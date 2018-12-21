package PowerOptions

import (
/*
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
*/	
	"github.com/cuu/LauncherGoDev/sysgo/UI"
	//"github.com/cuu/LauncherGoDev/sysgo/DBUS"
)

/******************************************************************************/
type PowerOptionsPlugin struct {
	UI.Plugin
  PowerOptionsPage *PowerOptionsPage
}


func (self *PowerOptionsPlugin) Init( main_screen *UI.MainScreen ) {
	self.PowerOptionsPage = NewPowerOptionsPage()
	self.PowerOptionsPage.SetScreen( main_screen)
	self.PowerOptionsPage.SetName("PowerOptions")
	self.PowerOptionsPage.Init()  
}

func (self *PowerOptionsPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.PowerOptionsPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ PowerOptionsPlugin
