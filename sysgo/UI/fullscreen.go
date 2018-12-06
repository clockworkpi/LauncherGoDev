package UI

import (

  "github.com/cuu/gogame/surface"

)

type FullScreen struct {
  Widget
  CanvasHWND *sdl.Surface
  HWND       *sdl.Surface
  
}

func NewFullScreen() *FullScreen {
  p := &FullScreen{}
  
  return p

}

func (self *FullScreen) SwapAndShow() {
  if self.HWND !=nil {
    rect_ := rect.Rect(self.PosX,self.PosY,self.Width,self.Height)
    surface.Blit(self.HWND,self.CanvasHWND,&rect_,nil)
    SwapAndShow()
  }

}

func (self *FullScreen) Draw() {


}
