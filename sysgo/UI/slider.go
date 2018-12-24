package UI

import (
	//"fmt"
  //"math"
	//"sync"
	
	"github.com/veandco/go-sdl2/sdl"

	//"github.com/cuu/gogame/surface"
	//"github.com/cuu/gogame/draw"
  //"github.com/cuu/gogame/rect"
  //"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/event"

	//"github.com/cuu/gogame/transform"
	//"github.com/clockworkpi/LauncherGoDev/sysgo/easings"
	
)

type SliderInterface interface {
  WidgetInterface
  
  Init()
  SetValue()
  SetRange(m1,m2 int)
  SetCanvasHWND( canvas *sdl.Surface)
  KeyDown(ev *event.Event)
  Draw()
}

type Slider struct {
  Widget
  
  Value int
  
  CanvasHWND *sdl.Surface
  
  Range [2]int
}

func NewSlider() *Slider {
  p := &Slider{}
  p.Range = [2]int{0,255}
  p.Value = 0
  return p
}

func (self *Slider) Init() {
  self.Value = 0
}

func (self *Slider) SetValue(v int) {
  self.Value = v
}

func (self *Slider) SetRange(m1 ,m2 int) {
  if m1 >= m2 {
    return
  }
  self.Range[0] = m1
  self.Range[1] = m2
}

func (self *Slider)	SetCanvasHWND( canvas *sdl.Surface) {
	self.CanvasHWND = canvas
}

func (self *Slider) KeyDown(ev *event.Event) {
}

func (self *Slider) Draw() {
  
}
