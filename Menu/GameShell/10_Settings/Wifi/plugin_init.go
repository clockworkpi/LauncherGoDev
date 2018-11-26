package main 

import (
/*
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
*/	
	"github.com/cuu/LauncherGo/sysgo/UI"
	"github.com/cuu/LauncherGo/sysgo/DBUS"
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
  self.ScanPage.Wireless = DBUS.DBusHandler.Wifi
  self.ScanPage.Daemon   = DBUS.DBusHandler.Daemon
  
  self.ScanPage.Screen = main_screen
  
  self.ScanPage.Init()
  
}

func (self *WifiPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
		main_screen.PushPage(self.ScanPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ WifiPlugin
