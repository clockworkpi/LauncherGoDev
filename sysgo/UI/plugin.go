package UI

import (
	"plugin"
)
var PluginPool = NewPluginPool()

type PluginInterface {
	Name() string
	Init(screen *MainScreen)
	Run(screen *MainScreen)
}

type Plugin struct {
	Name string // only ID for plugin
}

func (self *Plugin) Name() string {
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
	name := pi.Name()
	
	if _,ok := PluginPool[name]; ok {
		return false
	}
	
	PluginPool[name] = pi
}

func LoadPlugin( pname string) (*plugin.Plugin,error) {
	return plugin.Load(pname)
}

func InitPlugin(p *plugin.Plugin, main_screen *MainScreen) {
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

	pi.Init(main_screen)
}

func RunPlugin(p *plugin.Plugin, main_screen *MainScreen) {
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

