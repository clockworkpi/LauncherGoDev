package Emulator

import (

  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

type ActionConfig struct {
	ROM string `ini:"ROM"`
	ROM_SO string `ini:"ROM_SO"`
	EXT []string   `ini:"EXT,omitempty"`
	EXCLUDE []string `ini:"EXCLUDE,omitempty"`
	FILETYPE string  `ini:"FILETYPE"`   // defalut is file
	LAUNCHER string  `ini:"LAUNCHER"`
	TITLE   string   `ini:"TITLE"` // defaut is Game
	SO_URL string    `ini:"SO_URL"`
	RETRO_CONFIG string `ini:"RETRO_CONFIG"`
}


var (
  FavGID = 31415
  FavGname = "cpifav"
  
)


type MyEmulator struct { // as leader of RomListPage and FavListPage, it's a PluginInterface
  Name string
  RomPage *RomListPage
  FavPage *FavListPage
  DeleteConfirmPage *UI.DeleteConfirmPage
  EmulatorConfig *ActionConfig
  
  SpeedMax int
  SpeedTimeInter int
  
}

func NewMyEmulator() *MyEmulator{
  p := &MyEmulator{}
  
  p.SpeedMax = 5
  p.SpeedTimeInter = 300
  
  return p
}

func (self *MyEmulator) GetName() string {
  return "MyEmulator"
}

func (self *MyEmulator) Init(main_screen  *UI.MainScreen) {
  self.DeleteConfirmPage = UI.NewDeleteConfirmPage()
  self.DeleteConfirmPage.Screen = main_screen
  self.DeleteConfirmPage.Name  = "Delete Confirm"
  self.DeleteConfirmPage.Init()

  self.RomPage = NewRomListPage()
  self.RomPage.Screen = main_screen
  self.RomPage.Name  = self.EmulatorConfig.TITLE
  self.RomPage.EmulatorConfig = self.EmulatorConfig
  self.RomPage.Leader = self
  self.RomPage.Init()
  
  self.FavPage = NewFavListPage()
  self.FavPage.Screen = main_screen
  self.FavPage.Name = "FavouriteGames"
  self.FavPage.EmulatorConfig = self.EmulatorConfig
  self.FavPage.Leader = self
  self.FavPage.Init()
  
  
  
}

func (self *MyEmulator) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.RomPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}



