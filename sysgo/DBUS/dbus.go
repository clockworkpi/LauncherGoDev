package DBUS

import (
	"fmt"
	//"strconv"
	"github.com/godbus/dbus"
)

type DbusInterface struct {
	Dest string
	Path dbus.ObjectPath
	Iface string
	Obj *dbus.Object
  SigFuncs  map[string]interface{}
}

func NewDbusInterface(conn *dbus.Conn,dest string, path dbus.ObjectPath ,iface string) *DbusInterface {
  m := &DbusInterface{}
	o := conn.Object(dest,path)

	m.Obj = o.(*dbus.Object)
	m.Dest = dest
	m.Path = path
  
	m.SigFuncs = make(map[string]interface{})

	if len(iface) > 2 {
		m.Iface = iface
	}
	return m
}

func (self *DbusInterface) Method(name string, args...interface{} ) *dbus.Call {
	var method string
	if self.Iface != "" {
		method = fmt.Sprintf("%s.%s.%s", self.Dest, self.Iface,name)
	}else {
		method = fmt.Sprintf("%s.%s", self.Dest,name)
	}

	if args != nil {
		return self.Obj.Call( method , 0,args...)
	}else {
		return self.Obj.Call( method, 0)
	}
}

func (self *DbusInterface) Get( thecall *dbus.Call, retvalues ...interface{}) error {
	if len(thecall.Body) == 0 {
		return nil
	}
	err:=  thecall.Store(retvalues...)
  
	if err != nil {
		panic(fmt.Sprintf("Failed: %s,%s",err,thecall.Method))
	}
  
  return err
}

func (self *DbusInterface) EnableSignal(signame string) {
  iface := self.Dest
  if self.Iface != "" {
    iface = iface+ "."+self.Iface
  }
  self.Obj.AddMatchSignal(iface,signame)
}


func (self *DbusInterface) HandleSignal( sig *dbus.Signal) {
	
	iface := self.Dest
	if self.Iface != "" {
		iface = iface+ "."+self.Iface
	}

	if strings.HasPrefix(sig.Name,iface) {
		func_name := strings.Replace( sig.Name, iface, "",-1)[1:]
		for k,v := range self.SigFuncs {
			if k == func_name {
				v.(func([]interface{}))(sig.Body)
				break
			}
		}
	}

}

type DBusInterface interface {
	WifiStrength() int
	IsWifiConnectedNow() bool 	
}

type DBus struct {
	Conn *dbus.Conn
	Daemon *DbusInterface
  Wifi    *DbusInterface
}

func NewDBus() *DBus {
	d := &DBus{}
	return d
}

func (self *DBus) Init() {
  //conn_option := dbus.WithSignalHandler(self)
  
	conn, err := dbus.SystemBus()
  //conn,err := dbus.SystemBusPrivate(conn_option)
  
  if err != nil {
    panic(fmt.Sprintf("Failed to connect to system bus:", err))
  }

	self.Conn = conn

	self.Daemon = NewDbusInterface(conn, "org.wicd.daemon","/org/wicd/daemon" ,"",)
	self.Wifi   = NewDbusInterface(conn, "org.wicd.daemon","/org/wicd/daemon/wireless","wireless")
}

func (self *DBus) WifiStrength() int {
	var fast bool
	var iwconfig string
  var sig_display_type int
	var strength int

	self.Daemon.Get( self.Daemon.Method("NeedsExternalCalls"), &fast)

	if fast == false {
			self.Wifi.Get( self.Wifi.Method("GetIwconfig"), &iwconfig  )
	}else{
		iwconfig = ""
	}
	self.Daemon.Get( self.Daemon.Method("GetSignalDisplayType"), &sig_display_type )
	
	if sig_display_type == 0 {
		self.Wifi.Get( self.Wifi.Method("GetCurrentSignalStrength",iwconfig), &strength)
	} else{
		self.Wifi.Get( self.Wifi.Method("GetCurrentDBMStrength",iwconfig), &strength)
	}

	return strength	
}

func (self *DBus) check_for_wireless(iwconfig string, wireless_ip string)  bool {
  var network string
  var sig_display_type int
	var strength int
	if wireless_ip == "" {
		return false
	}

	self.Wifi.Get( self.Wifi.Method("GetCurrentNetwork",iwconfig), &network)
	self.Daemon.Get( self.Daemon.Method("GetSignalDisplayType"), &sig_display_type )
	
	if sig_display_type == 0 {
		self.Wifi.Get( self.Wifi.Method("GetCurrentSignalStrength",iwconfig), &strength)
	}else {
		self.Wifi.Get( self.Wifi.Method("GetCurrentDBMStrength",iwconfig), &strength)
	}

	if strength == 0 {
		return false
	}
	strength_str := ""
	self.Daemon.Get( self.Daemon.Method("FormatSignalForPrinting",strength), &strength_str)

	return true
}

func (self *DBus) IsWifiConnectedNow() bool {
  var fast bool
  var iwconfig string
  var wireless_connecting bool
	var wireless_ip string

	self.Wifi.Get( self.Wifi.Method("CheckIfWirelessConnecting"), &wireless_connecting  )
	self.Daemon.Get( self.Daemon.Method("NeedsExternalCalls"), &fast)
	if wireless_connecting == true {
		return false
	}else {
		if fast == false {
			self.Wifi.Get( self.Wifi.Method("GetIwconfig"), &iwconfig  )
		}else {
			iwconfig = ""
		}

		self.Wifi.Get( self.Wifi.Method("GetWirelessIP", iwconfig), &wireless_ip)
		
		if self.check_for_wireless(iwconfig,wireless_ip) == true {
			return true
		}else {
			return false
		}
		
	}	
}

func (self *DBus) ListenSignal() {
	c := make(chan *dbus.Signal, 10)
	self.Conn.Signal(c)
  
  for v := range c {
    fmt.Printf("%+v %#v\n",v,v)
    fmt.Printf("body len:%d \n\n",len(v.Body)) 
    
    self.Wifi.HandleSignal(v)
    self.Daemon.HandleSignal(v)
    
  }  
}

var DBusHandler *DBus //global 

func init() {
  DBusHandler = NewDBus()
  DBusHandler.Init()
  
  go DBusHandler.ListenSignal()
  
  
}
