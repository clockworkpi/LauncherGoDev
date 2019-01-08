package GateWay

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
type GatewayPlugin struct {
	UI.Plugin
  Page1st *GateWayPage
}


func (self *GatewayPlugin) Init( main_screen *UI.MainScreen ) {
	self.Page1st = NewGateWayPage()
	self.Page1st.SetScreen( main_screen)
	self.Page1st.SetName("Gateway switch")
	self.Page1st.Init()  
}

func (self *GatewayPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.Page1st)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ GatewayPlugin
