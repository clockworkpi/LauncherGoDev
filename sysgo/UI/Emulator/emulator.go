package Emulator

import (

  "github.com/cuu/LauncherGoDev/sysgo/UI"
)

type ActionConfig struct {
	ROM string `json:"ROM"`
	ROM_SO string `json:"ROM_SO"`
	EXT []string   `json:"EXT"`
	EXCLUDE []string `json:"EXCLUDE"`
	FILETYPE string  `json:"FILETYPE"`   // defalut is file
	LAUNCHER string  `json:"LAUNCHER"`
	TITLE   string   `json:"TITLE"` // defaut is Game
	SO_URL string    `json:"SO_URL"`
	RETRO_CONFIG string `json:"RETRO_CONFIG"`
}


var (
  FavGID = 31415
  FavGname = "cpifav"
  
)


type MyEmulator struct { // as leader of RomListPage and FavListPage, it's a PluginInterface
  Name string
  RomPage *RomListPage
  FavPage *FavListPage
}

func NewMyEmulator() *MyEmulator{
  p := &MyEmulator{}
  
  return p
}

func (self *MyEmulator) GetName() string {
  return "MyEmulator"
}

func (self *MyEmulator) Init(main_screen  *UI.MainScreen) {
  
}

func (self *MyEmulator) API(main_screen *UI.MainScreen) {
  
}



