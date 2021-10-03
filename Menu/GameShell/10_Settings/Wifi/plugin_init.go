package Wifi

import (
  //gotime "time"
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
type WifiPlugin struct {
	UI.Plugin
	ScanPage *WifiList
  PasswordPage *UI.Keyboard
}


func (self *WifiPlugin) Init( main_screen *UI.MainScreen ) {

  self.PasswordPage = UI.NewKeyboard()
  self.PasswordPage.Name = "Enter wifi password"
  self.PasswordPage.Screen= main_screen
  self.PasswordPage.Init()
  
  self.ScanPage = NewWifiList()
  self.ScanPage.Name = "Scan wifi"
  
  self.ScanPage.Screen = main_screen
  
  self.PasswordPage.Caller = self.ScanPage
  
  self.ScanPage.Init()
  
}

func (self *WifiPlugin) Run( main_screen *UI.MainScreen ) {
  if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.ScanPage)
    main_screen.Draw()
    main_screen.SwapAndShow()
  }
}

var APIOBJ WifiPlugin
