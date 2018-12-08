package Wifi

import (
  gotime "time"
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
  self.ScanPage.Wireless = main_screen.DBusManager.Wifi
  self.ScanPage.Daemon   = main_screen.DBusManager.Daemon
  
  self.ScanPage.Screen = main_screen
  
  self.ScanPage.Init()
  
  go func() {
    gotime.Sleep(2000 * gotime.Millisecond)
    
    if self.ScanPage.Daemon != nil {
    
      self.ScanPage.Daemon.EnableSignal("StatusChanged")
      self.ScanPage.Daemon.EnableSignal("ConnectResultsSent")
      self.ScanPage.Wireless.EnableSignal("SendStartScanSignal")
      self.ScanPage.Wireless.EnableSignal("SendEndScanSignal")
    
      self.ScanPage.Daemon.SigFuncs["StatusChanged"] = self.ScanPage.DbusDaemonStatusChangedSig
      self.ScanPage.Daemon.SigFuncs["ConnectResultSent"] = self.ScanPage.DbusConnectResultsSent
    
      self.ScanPage.Wireless.SigFuncs["SendStartScanSignal"] = self.ScanPage.WifiDbusScanStarted
      self.ScanPage.Wireless.SigFuncs["SendEndScanSignal"]   = self.ScanPage.WifiDbusScanFinishedSig
    }
  }()
  
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
