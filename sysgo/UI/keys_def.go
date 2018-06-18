package UI

import (
	"../../sysgo"
)


var CurKeys map[string]string

var GameShell map[string]string
var PC        map[string]string


func DefinePC() {
	PC["UP"] = "Up"
	PC["Down"] = "Down"
	PC["Left"] = "Left"
	PC["Right"] = "Right"
	PC["Menu"] = "Escape"
	PC["X"]    = "X"
	PC["Y"]    = "Y"
	PC["A"]    = "A"
	PC["B"]    = "B"

	PC["Vol-"]  = "Space"
	PC["Vol+"]  = "Return"
	PC["Space"] = "Space"
	PC["Enter"] = "Return"
	PC["Start"] = "S"
}

func DefineGameShell() {
	GameShell["UP"] = "Up"
	GameShell["Down"] = "Down"
	GameShell["Left"] = "Left"
	GameShell["Right"] = "Right"
	GameShell["Menu"] = "Escape"
	GameShell["X"]    = "U"
	GameShell["Y"]    = "I"
	GameShell["A"]    = "J"
	GameShell["B"]    = "K"

	GameShell["Vol-"]  = "Space"
	GameShell["Vol+"]  = "Return"
	GameShell["Space"] = "Space"
	GameShell["Enter"] = "K"
	GameShell["Start"] = "Return"
}

func init(){
	GameShell = make(map[string]string)
	PC        = make(map[string]string)

	DefineGameShell()
	DefinePC()
	
	if sysgo.CurKeySet == "GameShell" {
		CurKeys = GameShell
	}else {
		CurKeys = PC
	}
}
