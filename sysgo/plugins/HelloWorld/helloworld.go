package main

import (
	"../../UI"
)

type HelloWorldPage struct {
	UI.Page
}




type HelloWorldPlugin struct {
	UI.Plugin
}


func (self *HelloWorldPlugin) Init( main_screen *UI.MainScreen ) {
	
}

func (self *HelloWorldPlugin) Run( main_screen *UI.MainScreen ) {
	
}

var APIOBJ HelloWorldPlugin





