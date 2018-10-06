package UI

import (
	
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/surface"
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

func (self *KeyboardIcon) Draw() {

	rect_ := draw.MidRect(self.PosX,self.PosY,self.Width,self.Height,Width,Height)
	
	surface.Blit(self.Parent.GetCanvasHWND(),self.ImgSurf,rect_,nil)
	
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

type Keyboard struct {
	Page

	SectionNumbers int
	SectionIndex int
	Icons  map[string]IconItemInterface

	KeyboardLayoutFile string ///sysgo/UI/keyboard_keys.layout

	LeftOrRight int

	RowIndex int

	Textarea *Textarea
	Selector *KeyboardSelector

	
}

func NewKeyboard() *Keyboard {
	p := &Keyboard{}

	p.SectionNumbers = 3
	p.SectionIndex = 1

	p.Icons =  make( map[string]IconItemInterface )

	p.LeftOrRight = 1

	p.RowIndex = 0
	
	p.FootMsg = [5]string{"Nav.","ABC","Done","Backspace","Enter"}
	
	return p
	
}

func (self *Keyboard) ReadLayoutFile( fname string) {

	/*
	LayoutIndex := 0

	content ,err := ReadLines(fname)
  */
	
	
	
}
