package UI

import (
	"log"
	goplugin "plugin"
)

var PluginPool = NewPluginPool()

type PluginInterface interface {
	GetName() string
	Init(screen *MainScreen)
	Run(screen *MainScreen)
}

type Plugin struct {
	Name string // only ID for plugin
}

func (self *Plugin) GetName() string {
	return self.Name
}

func (self *Plugin) Init( screen *MainScreen) {
	
}


func (self *Plugin) Run( screen *MainScreen) {
	
}


func NewPluginPool() map[string]PluginInterface {
	pp :=make( map[string]PluginInterface )
	return pp
}

func PluginPoolRegister( pi PluginInterface ) bool {
	name := pi.GetName()
	
	if _,ok := PluginPool[name]; ok {
		return false
	}
	
	PluginPool[name] = pi
	return true
	
}

func LoadPlugin( pname string) (*goplugin.Plugin,error) {
	return goplugin.Open(pname)
}

func InitPlugin(p *goplugin.Plugin, main_screen *MainScreen) PluginInterface {
	symAPI,err := p.Lookup("APIOBJ")

	if err!= nil {
		log.Fatal( "init plugin failed")
		return nil
	}

	var pi PluginInterface
	pi,ok := symAPI.(PluginInterface)
	if !ok {
		log.Fatal("unexpected type from module symbol")
		return nil
	}

	//PluginPoolRegister(pi)
	
	pi.Init(main_screen)

	return pi
}

func RunPlugin(p *goplugin.Plugin, main_screen *MainScreen) {
	symAPI,err := p.Lookup("APIOBJ")

	if err!= nil {
		log.Fatal( "init plugin failed")
		return
	}

	var pi PluginInterface
	pi,ok := symAPI.(PluginInterface)
	if !ok {
		log.Fatal("unexpected type from module symbol")
		return
	}
	pi.Run(main_screen)
}

const (
  PluginPackage = iota
  PluginSo
)

type UIPlugin struct{ //Loadable and injectable
  Type int // 0 == loadable package, 1  == .so 
  SoFile string
  FolderName string
  LabelText  string
  EmbInterface  PluginInterface
}
