package Storage


import (
/*
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
*/	
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"

)

/******************************************************************************/
type StoragePlugin struct {
	UI.Plugin
  StoragePage *StoragePage
}


func (self *StoragePlugin) Init( main_screen *UI.MainScreen ) {
	self.StoragePage = NewStoragePage()
	self.StoragePage.SetScreen( main_screen)
	self.StoragePage.SetName("Storage")
	self.StoragePage.Init()  
}

func (self *StoragePlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.StoragePage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ StoragePlugin
