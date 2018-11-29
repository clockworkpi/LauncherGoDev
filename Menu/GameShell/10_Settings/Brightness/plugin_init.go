package main

/*
 * sysgo.BackLight  
 */
import (
/*
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
*/	
	"github.com/cuu/LauncherGo/sysgo/UI"
	//"github.com/cuu/LauncherGo/sysgo/DBUS"
)

/******************************************************************************/
type BrightnessPlugin struct {
	UI.Plugin
  BrightnessPage *BrightnessPage
}


func (self *BrightnessPlugin) Init( main_screen *UI.MainScreen ) {
  
}

func (self *BrightnessPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.BrightnessPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ BrightnessPlugin
