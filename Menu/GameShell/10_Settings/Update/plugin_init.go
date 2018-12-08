package Update
import (
  "github.com/cuu/LauncherGoDev/sysgo/UI"
)
/******************************************************************************/
type UpdatePlugin struct {
	UI.Plugin
	Page UI.PageInterface
}


func (self *UpdatePlugin) Init( main_screen *UI.MainScreen ) {
	self.Page = NewUpdatePage()
	self.Page.SetScreen( main_screen)
	self.Page.SetName("Update")
	self.Page.Init()
}

func (self *UpdatePlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
		main_screen.PushPage(self.Page)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ UpdatePlugin
