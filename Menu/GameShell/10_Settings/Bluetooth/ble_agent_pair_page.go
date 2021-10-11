package Bluetooth

import (
	"fmt"
	//"os"
	//"log"
	//"strings"

	//"github.com/fatih/structs"
	/*
	  "github.com/veandco/go-sdl2/ttf"
	  "github.com/cuu/gogame/draw"
	  "github.com/cuu/gogame/surface"
	  "github.com/cuu/gogame/rect"

	  "github.com/cuu/gogame/color"
	  "github.com/cuu/gogame/font"
	*/
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/time"
	//"github.com/godbus/dbus"
	//"github.com/muka/go-bluetooth/api"
	//"github.com/muka/go-bluetooth/bluez"
	//"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

type BleAgentPairPage struct {
	UI.Page

	Pin    string
	Pass   string
	DevObj *device.Device1
	Leader *BluetoothPlugin
}

func NewBleAgentPairPage() *BleAgentPairPage {
	p := &BleAgentPairPage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav.", "", "", "Back", ""}

	return p
}

func (self *BleAgentPairPage) Init() {
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.CanvasHWND = self.Screen.CanvasHWND

}

func (self *BleAgentPairPage) ShowPinCode(device string, pincode string) {
	fmt.Println(fmt.Sprintf("ShowPinCode %s %d", device, pincode))
	if self.Screen.CurPage() != self {
		self.Screen.PushPage(self)
		self.ClearCanvas()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	self.Pin = pincode
	txt := self.Pin
	if len(self.Pin) > 0 {
		txt = fmt.Sprintf("Pin code: %s", self.Pin)
	}

	self.Screen.MsgBox.SetText(txt)
	self.Screen.MsgBox.Draw()
	self.Screen.SwapAndShow()
}

func (self *BleAgentPairPage) ShowPassKey(device string, passkey uint32, entered uint16) {
	fmt.Println(fmt.Sprintf("ShowPassKey %06d %d", passkey, entered))
	if self.Screen.CurPage() != self {
		self.Screen.PushPage(self)
		self.ClearCanvas()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	self.Pass = fmt.Sprintf("%06d", passkey)
	txt := self.Pass
	if len(self.Pass) > 0 {
		txt = fmt.Sprintf("Pair code: %s", self.Pass)
	}

	self.Screen.MsgBox.SetText(txt)
	self.Screen.MsgBox.Draw()
	self.Screen.SwapAndShow()

}

func (self *BleAgentPairPage) PairOKCb() {
	self.ClearCanvas()
	self.Screen.Draw()
	self.Screen.SwapAndShow()

	self.Screen.MsgBox.SetText("Device paired")
	self.Screen.MsgBox.Draw()
	self.Screen.SwapAndShow()

	time.BlockDelay(1500)

	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()
	self.Screen.FootBar.ResetNavText()

}

func (self *BleAgentPairPage) PairErrorCb(err_msg string) {
	self.ClearCanvas()
	self.Screen.Draw()
	self.Screen.SwapAndShow()

	self.Screen.MsgBox.SetText(err_msg)
	self.Screen.MsgBox.Draw()
	self.Screen.SwapAndShow()

	time.BlockDelay(1500)

	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()
	self.Screen.FootBar.ResetNavText()

}

func (self *BleAgentPairPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		if self.DevObj != nil {
			err := self.DevObj.CancelPairing()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
}

func (self *BleAgentPairPage) Draw() {
	// DoNothing
}
