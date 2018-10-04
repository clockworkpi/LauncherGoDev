package UI

import (
	"github.com/cuu/gogame/color"
)
//sysgo/UI/keyboard_keys.layout
type KeyboardIcon struct {
	IconItem
	Color *color.Color

	Str  string
}

func NewKeyboardIcon() *KeyboardIcon {
	p := &KeyboardIcon{}

	p.MyType = ICON_TYPES["NAV"]

	return p
}

type KeyboardSelector struct {
	PageSelector
	Parent *Keyboard
}


func NewKeyboardSelector() * KeyboardSelector {
	p := &KeyboardSelector{}

	return p
}

func (self *KeyboardSelector) Draw() {
	
}

type Keyboard {
	Page

	SectionNumbers int
	SectionIndex int
	Icons  map[string]UI.IconItemInterface

	KeyboardLayoutFile string ///sysgo/UI/keyboard_keys.layout

	
}
