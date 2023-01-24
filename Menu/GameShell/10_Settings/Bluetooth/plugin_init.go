package Bluetooth

import (
	"fmt"
	//"log"
	//"os"
	//"time"
	"github.com/godbus/dbus/v5"
	//"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/hw"
	//"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	//"github.com/muka/go-bluetooth/bluez/profile/adapter"
	/*
	   "github.com/veandco/go-sdl2/ttf"

	   "github.com/cuu/gogame/surface"
	   "github.com/cuu/gogame/event"
	   "github.com/cuu/gogame/rect"
	   "github.com/cuu/gogame/color"
	*/

	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	log "github.com/sirupsen/logrus"
)

/******************************************************************************/
type BluetoothPlugin struct {
	UI.Plugin
	BluetoothPage *BluetoothPage
	PairPage      *BleAgentPairPage
}

const (
	adapterID       = "hci0"
	BUS_NAME        = "org.bluez"
	AGENT_INTERFACE = "org.bluez.Agent1"
)

func (self *BluetoothPlugin) InitAgent() {

	conn, err := dbus.SystemBus()
	if err != nil {
		return
	}

	ag := agent.NewSimpleAgent()
	err = agent.ExposeAgent(conn, ag, agent.CapKeyboardDisplay, true)
	if err != nil {
		fmt.Println(fmt.Errorf("SimpleAgent: %s", err))
		return
	}

}

func (self *BluetoothPlugin) Init(main_screen *UI.MainScreen) {

	log.Println("Reset bluetooth device")

	btmgmt := hw.NewBtMgmt(adapterID)
	btmgmt.SetPowered(true)

	self.BluetoothPage = NewBluetoothPage()
	self.BluetoothPage.SetScreen(main_screen)
	self.BluetoothPage.SetName("Bluetooth")
	self.BluetoothPage.Leader = self
	self.BluetoothPage.Init()

	self.PairPage = NewBleAgentPairPage()
	self.PairPage.SetScreen(main_screen)
	self.PairPage.SetName("Bluetooth pair")
	self.PairPage.Leader = self
	self.PairPage.Init()

	self.InitAgent()

	/*
	  a, err := adapter.GetAdapter(adapterID)
	  if err != nil {
	    fmt.Println(err)
	    return
	  }

	  discovery, cancel, err := api.Discover(a, nil)
	  if err != nil {
	    fmt.Println(err)
	  }

	  defer cancel()

	  wait := make(chan error)

	  go func() {
	    for dev := range discovery {
	      if dev == nil {
	        return
	      }
	      wait <- nil
	    }
	  }()

	  go func() {
	    sleep := 5
	    time.Sleep(time.Duration(sleep) * time.Second)
	    log.Debugf("Discovery timeout exceeded (%ds)", sleep)
	    wait <- nil
	  }()

	  err = <-wait
	  if err != nil {
	    fmt.Println(err)
	  }
	*/

	//self.BluetoothPage.RefreshDevices()
	//self.BluetoothPage.GenNetworkList()

}

func (self *BluetoothPlugin) Run(main_screen *UI.MainScreen) {
	if main_screen != nil {
		main_screen.PushCurPage()
		main_screen.SetCurPage(self.BluetoothPage)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ BluetoothPlugin
