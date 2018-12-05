package LauncherPy

import (
  "os/exec"
  "github.com/cuu/gogame/time"
	"github.com/cuu/LauncherGoDev/sysgo/UI"
  
)
/******************************************************************************/
type LauncherPyPlugin struct {
	UI.Plugin
}

func (self *LauncherPyPlugin) Init( main_screen *UI.MainScreen ) {

}

func (self *LauncherPyPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.MsgBox.SetText("Rebooting to Launcher")
    main_screen.MsgBox.Draw()
		main_screen.SwapAndShow()
    time.BlockDelay(300)
    cmd := exec.Command("sed","-i","s/launchergo/launcher/g","~/.bashrc")
    cmd.Run()
    
    cmd = exec.Command("sudo","reboot")
    cmd.Run()
	}
}

var APIOBJ LauncherPyPlugin
