package Airplane

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
type AirplanePlugin struct {
	UI.Plugin
  AirplanePage *AirplanePage
}


func (self *AirplanePlugin) Init( main_screen *UI.MainScreen ) {
	self.AirplanePage = NewAirplanePage()
	self.AirplanePage.SetScreen( main_screen)
	self.AirplanePage.SetName("Airplane")
	self.AirplanePage.Init()  
}

func (self *AirplanePlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.AirplanePage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ AirplanePlugin
