package Lima

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
type LimaPlugin struct {
	UI.Plugin
  GPUDriverPage *GPUDriverPage
}


func (self *LimaPlugin) Init( main_screen *UI.MainScreen ) {
	self.GPUDriverPage = NewGPUDriverPage()
	self.GPUDriverPage.SetScreen( main_screen)
	self.GPUDriverPage.SetName("GPU driver switch")
	self.GPUDriverPage.Init()  
}

func (self *LimaPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.GPUDriverPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ LimaPlugin
