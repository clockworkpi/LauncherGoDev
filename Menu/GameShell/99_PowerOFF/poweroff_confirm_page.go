package PowerOFF

import (
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/cuu/gogame/event"
)

type PowerOFFConfirmPage struct {
	UI.ConfirmPage
}

func NewPowerOFFConfirmPage() *PowerOFFConfirmPage {

	p := &PowerOFFConfirmPage{}
	p.ListFont = UI.Fonts["veramono20"]
	p.ConfirmText = "Awaiting Input"
	p.FootMsg = [5]string{"Nav", "Reboot", "", "Cancel", "Shutdown"}

	p.ConfirmPage.ConfirmText = p.ConfirmText
	p.ConfirmPage.FootMsg = p.FootMsg
	p.ConfirmPage.ListFont = p.ListFont

	return p

}

func (self *PowerOFFConfirmPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["Menu"] || ev.Data["Key"] == UI.CurKeys["A"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()

	}

	if ev.Data["Key"] == UI.CurKeys["B"] {
		cmdpath := ""

		if UI.CheckBattery() < 20 {
			cmdpath = "feh --bg-center sysgo/gameshell/wallpaper/gameover.png;"
		} else {
			cmdpath = "feh --bg-center sysgo/gameshell/wallpaper/seeyou.png;"
		}

		cmdpath = cmdpath + "sleep 3;"

		cmdpath = cmdpath + "sudo halt -p"

		event.Post(UI.RUNSYS, cmdpath)

	}

	if ev.Data["Key"] == UI.CurKeys["X"] {
		cmdpath := "feh --bg-center sysgo/gameshell/wallpaper/seeyou.png;"
		cmdpath += "sleep 3;"
		cmdpath += "sudo reboot"

		event.Post(UI.RUNSYS, cmdpath)
	}

}
