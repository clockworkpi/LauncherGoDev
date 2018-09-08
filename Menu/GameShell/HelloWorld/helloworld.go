package main

import (
	"github.com/veandco/go-sdl2/ttf"
	
	"github.com/cuu/LauncherGo/sysgo/UI"
)

type InfoPageListItem struct{
	PosX   int
	PosY   int
	Width  int
	Height int
	Labels map[string]LabelInterface
	Icons  map[string]IconItemInterface
	Fonts  map[string]*ttf.Font
}

func (self *InfoPageListItem) SetSmallText( text string) {
	
}

type HelloWorldPage struct {
	UI.Page
}

func NewHelloWorldPage() *HelloWorldPage {
	p := &HelloWorldPage{}
	
	p.FootMsg = [5]string{"Nav.","","","Back",""}

	return p
}



type HelloWorldPlugin struct {
	UI.Plugin
	Page UI.PageInterface
}


func (self *HelloWorldPlugin) Init( main_screen *UI.MainScreen ) {
	self.Page = HelloWorldPage{}
	self.Page.SetScreen( main_screen)
	self.Page.SetName("HelloWorld")
	self.Page.Init()
}

func (self *HelloWorldPlugin) Run( main_screen *UI.MainScreen ) {
	if main_screen != nil {
		main_screen.PushPage(self.Page)
		main_screen.Draw()
		main_screen.SwapAndShow()
	}
}

var APIOBJ HelloWorldPlugin





