package Bluetooth

import (
  "log"
  "os"
  
	"github.com/muka/go-bluetooth/api"
	//"github.com/muka/go-bluetooth/emitter"
	"github.com/muka/go-bluetooth/linux"


/*
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
*/	

	"github.com/cuu/LauncherGoDev/sysgo/UI"
	//"github.com/cuu/LauncherGoDev/sysgo/DBUS"
)

/******************************************************************************/
type BluetoothPlugin struct {
	UI.Plugin
  BluetoothPage *BluetoothPage
}

const (
  adapterID = "hci0"
)


func (self *BluetoothPlugin) Init( main_screen *UI.MainScreen ) {
  
 	log.Println("Reset bluetooth device")
  
	a := linux.NewBtMgmt(adapterID)
	err := a.Reset()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	
  err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device
		showDeviceInfo(dev)
	}))
  
  if err != nil {
    fmt.Println(err)
  }
    

	self.BluetoothPage = NewBluetoothPage()
	self.BluetoothPage.SetScreen( main_screen)
	self.BluetoothPage.SetName("Bluetooth")
	self.BluetoothPage.Init()  
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
