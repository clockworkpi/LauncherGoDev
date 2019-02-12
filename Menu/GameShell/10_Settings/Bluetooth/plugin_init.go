package Bluetooth

import (
  "fmt"
  "log"
  "os"
  
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/linux/btmgmt"
  "github.com/muka/go-bluetooth/bluez/profile"

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
type BluetoothPlugin struct {
	UI.Plugin
  BluetoothPage *BluetoothPage
  PairPage      *BleAgentPairPage
}

const (
  adapterID = "hci0"
	BUS_NAME = "org.bluez"
	AGENT_INTERFACE = "org.bluez.Agent1"
	AGENT_PATH = "/gameshell/bleagentgo"  
  
)


func (self *BluetoothPlugin) InitAgent() {
	agent := &Agent{}
	agent.BusName = BUS_NAME
	agent.AgentInterface = AGENT_INTERFACE
	agent.AgentPath = AGENT_PATH
  agent.Leader = self
	RegisterAgent(agent, profile.AGENT_CAP_KEYBOARD_DISPLAY)
}

func (self *BluetoothPlugin) Init( main_screen *UI.MainScreen ) {
  
 	log.Println("Reset bluetooth device")
  
	a := btmgmt.NewBtMgmt(adapterID)
	err := a.Reset()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	
	self.BluetoothPage = NewBluetoothPage()
	self.BluetoothPage.SetScreen( main_screen)
	self.BluetoothPage.SetName("Bluetooth")
  self.BluetoothPage.Leader  = self
	self.BluetoothPage.Init()  
  
  self.PairPage = NewBleAgentPairPage()
  self.PairPage.SetScreen( main_screen) 
  self.PairPage.SetName("Bluetooth pair")
  self.PairPage.Leader = self
  self.PairPage.Init()
  
  self.InitAgent()
  
  err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
		//discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		//dev := discoveryEvent.Device
		//showDeviceInfo(dev)
    self.BluetoothPage.RefreshDevices()
    self.BluetoothPage.GenNetworkList()
    main_screen.Draw()
    main_screen.SwapAndShow()
    
	}))
  
  if err != nil {
    fmt.Println(err)
  }  
  
}

func (self *BluetoothPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
    main_screen.PushCurPage()
    main_screen.SetCurPage(self.BluetoothPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ BluetoothPlugin
