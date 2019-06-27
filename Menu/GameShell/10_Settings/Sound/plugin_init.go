package Sound

/*
 * need amixer
 * `sudo apt-get install alsa-utils`
 */
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
type SoundPlugin struct {
	UI.Plugin
  SoundPage *SoundPage
}


func (self *SoundPlugin) Init( main_screen *UI.MainScreen ) {
	self.SoundPage = NewSoundPage()
	self.SoundPage.SetScreen( main_screen)
	self.SoundPage.SetName("Sound Volume")
	self.SoundPage.Init()  
}

func (self *SoundPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.SoundPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ SoundPlugin
