package Bluetooth

import (
  "fmt"
  //"os"
  //"log"
  "strings"
  
  //"github.com/fatih/structs"
  /*
  "github.com/veandco/go-sdl2/ttf"
  "github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/surface"
  "github.com/cuu/gogame/rect"
  
  "github.com/cuu/gogame/color"
  "github.com/cuu/gogame/font"
  */
  "github.com/cuu/gogame/time"
  "github.com/cuu/gogame/event"
  "github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
  //"github.com/muka/go-bluetooth/bluez"
  "github.com/muka/go-bluetooth/bluez/profile"
  
  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

func set_trusted(path string) {
	devices ,err := api.GetDevices()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	for i,v := range devices {
		fmt.Println(i, v.Path)
		if strings.Contains(v.Path,path) {
			fmt.Println("Found device")
			dev1,_ := v.GetClient()
			err:=dev1.SetProperty("Trusted",true)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}


type Agent struct{
	BusName string
	AgentInterface string
	AgentPath string
  Leader *BluetoothPlugin
}

func (self *Agent) Release()  *dbus.Error {
	return nil
}

func (self *Agent) RequestPinCode(device dbus.ObjectPath) (pincode string, err *dbus.Error) {
	fmt.Println("RequestPinCode",device)
	set_trusted(string(device))
	return "0000",nil
}

func (self *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	fmt.Println( fmt.Sprintf("DisplayPinCode (%s, %s)" ,device, pincode))
  self.Leader.PairPage.ShowPinCode(string(device),pincode)
	return nil
}

func (self *Agent) RequestPasskey(device dbus.ObjectPath)  (passkey uint32, err *dbus.Error) {
	set_trusted(string(device))
	return 0,nil
}

func (self *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	fmt.Println(fmt.Sprintf("DisplayPasskey %s, %06u entered %u" ,device, passkey, entered))
  self.Leader.PairPage.ShowPassKey(string(device),passkey,entered)
	return nil
}


func (self *Agent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	fmt.Println(fmt.Sprintf("RequestConfirmation (%s, %06d)", device, passkey))
	set_trusted(string(device))
	return nil
}


func (self *Agent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	fmt.Printf("RequestAuthorization (%s)\n" ,device)
	return nil
}

func (self *Agent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	fmt.Printf("AuthorizeService (%s, %s)",device, uuid) //directly authrized
	return nil
}

func (self *Agent) Cancel() *dbus.Error {
	fmt.Println("Cancel")	
	return nil
}

func (self *Agent) RegistrationPath() string {
	return self.AgentPath
}

func (self *Agent) InterfacePath() string {
	return self.AgentInterface
}

func RegisterAgent(agent profile.Agent1Interface, caps string) (err error) {
	//agent_path := AgentDefaultRegisterPath // we use the default path
	agent_path := agent.RegistrationPath() // we use the default path
	fmt.Println("The Agent Path: ", agent_path)
	// Register agent
	am := profile.NewAgentManager1(agent_path)

	// Export the Go interface to DBus
	err = am.ExportGoAgentToDBus(agent)
	if err != nil { return err }

	// Register the exported interface as application agent via AgenManager API
	err = am.RegisterAgent(agent_path, caps)
	if err != nil { return err }

	// Set the new application agent as Default Agent
	err = am.RequestDefaultAgent(agent_path)
	if err != nil { return err }

	return
}

type BleAgentPairPage struct {
  UI.Page
  
  Pin string
  Pass  string
  DevObj *api.Device
  Leader *BluetoothPlugin
}

func NewBleAgentPairPage() *BleAgentPairPage {
  p := &BleAgentPairPage{}
  p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]
	
	p.FootMsg = [5]string{"Nav.","","","Back",""}

  return p
}

func (self *BleAgentPairPage) Init() {
  self.PosX = self.Index * self.Screen.Width
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height
  
  self.CanvasHWND = self.Screen.CanvasHWND

}

func (self *BleAgentPairPage) ShowPinCode(device string,pincode string) {
  fmt.Println( fmt.Sprintf("ShowPinCode %s %d" ,device,pincode))
  if self.Screen.CurPage() != self {
    self.Screen.PushPage(self)
    self.ClearCanvas()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  self.Pin = pincode
  txt := self.Pin
  if len(self.Pin) > 0 {
    txt = fmt.Sprintf("Pin code: %s",self.Pin)
  }
  
  self.Screen.MsgBox.SetText(txt)
  self.Screen.MsgBox.Draw()
  self.Screen.SwapAndShow()
}

func (self *BleAgentPairPage) ShowPassKey(device string,passkey uint32,entered uint16) {
  fmt.Println(fmt.Sprintf("ShowPassKey %06d %d",passkey,entered) )
  if self.Screen.CurPage() != self {
    self.Screen.PushPage(self)
    self.ClearCanvas()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  self.Pass = fmt.Sprintf("%06d",passkey)  
  txt := self.Pass
  if len(self.Pass) > 0 {
    txt = fmt.Sprintf("Pair code: %s",self.Pass)
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

func (self *BleAgentPairPage) PairErrorCb( err_msg string) {
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
      c, err := self.DevObj.GetClient()
      if err != nil {
        fmt.Println(err)
        return
      }
      c.CancelPairing()
    }
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()          
  }
}

func (self *BleAgentPairPage) Draw() {
// DoNothing
}

