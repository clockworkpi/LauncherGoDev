package UI

import(
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
  
	"github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/rect"
  "github.com/cuu/gogame/color"
)

type AboveAllPatch struct {

  Widget
  
  Text string
  
  FontObj *ttf.Font
  
  Color *color.Color
  ValColor *color.Color
  
  CanvasHWND *sdl.Surface
  
  Icons map[string]IconItemInterface
  
  Value int 
}

func NewAboveAllPatch() *AboveAllPatch {
  
  p := &AboveAllPatch{}
  p.PosX = Width /2
  p.PosY = Height /2
  p.Width = 50
  p.Height = 120
  
  p.FontObj = Fonts["veramono20"]
  p.Color   = MySkinManager.GiveColor("Text")
  p.ValColor = MySkinManager.GiveColor("URL")
  
  p.Icons = make( map[string]IconItemInterface )
  
  p.Value = 0

  return p
}

func (self *AboveAllPatch) SetCanvasHWND( _canvashwnd *sdl.Surface) {
  
  self.CanvasHWND = _canvashwnd

}

func (self *AboveAllPatch) Draw() {
  start_rect := draw.MidRect(self.PosX,self.PosY,self.Width,self.Height,Width,Height)
  draw.AARoundRect(self.CanvasHWND,start_rect,self.Color,3,0,self.Color)
  
  if self.Value  > 10 {
    vol_height := int(float64(self.Height) * (float64(self.Value)/100.0))
    dheight    := self.Height - vol_height
    
    vol_rect   := rect.Rect(self.PosX - self.Width/2,self.PosY - self.Height/2+dheight,self.Width,vol_height)
    
    draw.AARoundRect(self.CanvasHWND,&vol_rect,self.ValColor,3,0,self.ValColor)
    
  }else {
      vol_height := 10
      dheight    := self.Height - vol_height
      vol_rect   := rect.Rect(self.PosX - self.Width/2,self.PosY - self.Height/2+dheight,self.Width,vol_height)
      
      draw.AARoundRect(self.CanvasHWND,&vol_rect,self.ValColor,3,0,self.ValColor)
  }
}


