package MusicPlayer

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
	"github.com/fhs/gompd/v2/mpd"

)

/******************************************************************************/
type MusicPlayerPlugin struct {
	UI.Plugin
	MusicPlayerPage *MusicPlayerPage
	MpdClient *mpd.Client
}

func (self *MusicPlayerPlugin) Init(main_screen *UI.MainScreen) {
	self.MusicPlayerPage = NewMusicPlayerPage()
	self.MusicPlayerPage.SetScreen(main_screen)
	self.MusicPlayerPage.SetName("Music Player")
	self.MpdClient = nil
	self.MusicPlayerPage.MpdClient = self.MpdClient
	self.MusicPlayerPage.Init()
}

func (self *MusicPlayerPlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.MusicPlayerPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ MusicPlayerPlugin
