package LauncherPy

import (
  "log"
  "os/exec"
  "os/user"
  "github.com/cuu/gogame/time"
  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
  
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
    time.BlockDelay(550)
    usr, _ := user.Current()
    dir := usr.HomeDir
    if usr.Username == "root" {
      dir = "/home/cpi"
    }

    cmd := exec.Command("sed","-i","s/launchergo/launcher/g",dir+"/.bashrc")
    err := cmd.Run()
    if err != nil {
	log.Println("sed failed", err)
    }else {
    
    cmd = exec.Command("sudo","reboot")
    cmd.Run()
    }
	}
}

var APIOBJ LauncherPyPlugin
