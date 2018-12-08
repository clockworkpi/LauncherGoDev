package Brightness

import(
  "fmt"
  
  "io/ioutil"
  "strconv"
  
  "github.com/veandco/go-sdl2/sdl"

  "github.com/cuu/gogame/event"
  "github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/surface"
  "github.com/cuu/gogame/rect"

  "github.com/cuu/LauncherGoDev/sysgo"
  "github.com/cuu/LauncherGoDev/sysgo/UI"
)

type OnChangeCB_T func(int)

type SliderIcon struct {
  UI.IconItem
  Parent *BSlider
  
}
func NewSliderIcon() *SliderIcon {
 	p := &SliderIcon{}
	p.MyType = UI.ICON_TYPES["EXE"]
	p.Align = UI.ALIGN["VCenter"]
  
  return p
}

func (self *SliderIcon) Draw() {
  if self.Parent == nil {
    fmt.Println("Error: SliderIcon Draw Parent nil")
    return
  }
	parent_x,parent_y := self.Parent.Coord()
	
	if self.Label != nil {
//		lab_x,lab_y := self.Label.Coord()
		lab_w,lab_h:= self.Label.Size()
		
		if self.Align == UI.ALIGN["VCenter"] {
//			fmt.Println("IconItem Draw VCenter:",lab_w,lab_h,self.Label.GetText())
			
			self.Label.NewCoord( self.PosX - lab_w/2 + parent_x, self.PosY + self.Height/2+6+parent_y)
			
		}else if self.Align == UI.ALIGN["HLeft"] {
			self.Label.NewCoord( self.PosX + self.Width/2+3+parent_x, self.PosY - lab_h/2 + parent_y)
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		surface.Blit(self.Parent.GetCanvasHWND(), self.ImgSurf,draw.MidRect(self.PosX + parent_x, self.PosY + parent_y,
			self.Width,self.Height, UI.Width, UI.Height),nil)
	}
}

type SliderMultiIcon struct {
  UI.MultiIconItem
  Parent *BSlider
}

func NewSliderMultiIcon() *SliderMultiIcon {
  p := &SliderMultiIcon{}
	p.MyType = UI.ICON_TYPES["EXE"]
	p.Align = UI.ALIGN["VCenter"]
    
  p.IconIndex = 0
	p.IconWidth = 18
	p.IconHeight = 18
  
  return p
}

func (self *SliderMultiIcon) Draw() {
  if self.Parent == nil {
    fmt.Println("Error: SliderMultiIcon Draw Parent nil")
    return
  }  
	parent_x,parent_y := self.Parent.Coord()
	
	if self.Label != nil {
//		lab_x,lab_y := self.Label.Coord()
		lab_w,lab_h:= self.Label.Size()
		if self.Align == UI.ALIGN["VCenter"] {
			self.Label.NewCoord( self.PosX - lab_w/2 + parent_x,        self.PosY + self.Height/2+6 + parent_y)
		}else if self.Align == UI.ALIGN["HLeft"] {
			self.Label.NewCoord( self.PosX + self.Width/2+3 + parent_x, self.PosY - lab_h/2 + parent_y )
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		
		portion := rect.Rect(0,self.IconIndex*self.IconHeight,self.IconWidth,self.IconHeight)
		
		surface.Blit(self.Parent.GetCanvasHWND(),
			self.ImgSurf,draw.MidRect(self.PosX + parent_x, self.PosY + parent_y,
			self.Width,self.Height, UI.Width, UI.Height),&portion)
	}
}

type BSlider struct {
    UI.Slider
    
    BGpng *SliderIcon
    BGwidth   int
    BGheight  int
    //NeedleSurf 
    Scale  *SliderMultiIcon
    Parent *BrightnessPage
    
    OnChangeCB OnChangeCB_T
    
}

func NewBSlider() *BSlider {
  p := &BSlider{}
  p.Range = [2]int{0,255}
  p.Value = 0
  
  p.BGwidth = 179
  p.BGheight = 153
  
  return p
}

func (self *BSlider) GetCanvasHWND() *sdl.Surface {
  return self.CanvasHWND
}

func (self *BSlider) Init() {

  self.Width = self.Parent.Width
  self.Height = self.Parent.Height
  
  self.BGpng = NewSliderIcon()
  self.BGpng.ImgSurf = UI.MyIconPool.GetImgSurf("light")
  self.BGpng.MyType = UI.ICON_TYPES["STAT"]
  self.BGpng.Parent = self
  self.BGpng.Adjust(0,0,self.BGwidth,self.BGheight,0)
  
  self.Scale = NewSliderMultiIcon()
  self.Scale.MyType = UI.ICON_TYPES["STAT"]
  self.Scale.Parent = self
  self.Scale.ImgSurf = UI.MyIconPool.GetImgSurf("scale")
  self.Scale.IconWidth = 82
  self.Scale.IconHeight = 63
  self.Scale.Adjust(0,0,82,63,0)

}

func (self *BSlider) SetValue( brt int) {
  self.Value = brt
}

func (self *BSlider) Further() {
  self.Value += 1
  
  if self.Value > 9 {
    self.Value = 9
  }

  if self.OnChangeCB != nil {
    self.OnChangeCB(self.Value)
  }
  
}

func (self *BSlider) StepBack() {
  self.Value -= 1
  
  if self.Value < 0 {
    self.Value = 0
  }
  
  if self.OnChangeCB != nil {
    self.OnChangeCB(self.Value)
  }
}

func (self *BSlider) Draw() {
  self.BGpng.NewCoord(self.Width/2,self.Height/2+11)
  self.BGpng.Draw()
  
  self.Scale.NewCoord(self.Width/2,self.Height/2)

  icon_idx := self.Value-1
  if icon_idx <0 {
    icon_idx = 0
  }
  
  self.Scale.IconIndex = icon_idx
  self.Scale.Draw()
  
}

type BrightnessPage struct {
  UI.Page
  MySlider *BSlider
}

func NewBrightnessPage() *BrightnessPage {
  p:= &BrightnessPage{}
  
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10
	p.Align = UI.ALIGN["SLeft"]
  
  p.FootMsg = [5]string{"Nav","","","Back","Enter"}  

  return p
}

func (self *BrightnessPage) Init() {
  self.CanvasHWND = self.Screen.CanvasHWND
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height
  
  self.MySlider = NewBSlider()
  
  self.MySlider.Parent = self
  
  self.MySlider.SetCanvasHWND(self.CanvasHWND)
  self.MySlider.OnChangeCB = self.WhenSliderDrag
  
  self.MySlider.Init()
  
  brt := self.ReadBackLight()
  
  self.MySlider.SetValue(brt)
  

}

func (self *BrightnessPage) ReadBackLight() int {
  
  if UI.FileExists(sysgo.BackLight) == false {
    return 0
  }
  
  lines,err := UI.ReadLines(sysgo.BackLight)
  
  if err != nil {
    fmt.Println(err)
    return 0
  }
  
  for _,v := range lines {
    n,e := strconv.Atoi(v)
    if e == nil {
      return n
    }else {
      fmt.Println(e)
      return 0
    }
    break
  }
  
  return 0
}

func (self *BrightnessPage) OnLoadCb() {
  brt := self.ReadBackLight()
  
  self.MySlider.SetValue(brt)

}

func (self *BrightnessPage) SetBackLight( newbrt int){
  
  newbrt_str := fmt.Sprintf("%d",newbrt)
  
  if UI.FileExists(sysgo.BackLight) {
    err:= ioutil.WriteFile(sysgo.BackLight,[]byte(newbrt_str),0644)
    if err != nil {
      fmt.Println(err)
    }
  }else{
    fmt.Println(sysgo.BackLight, " file not existed")
  }
}

func (self *BrightnessPage) WhenSliderDrag( val int) {
  self.SetBackLight(val)
}

func (self *BrightnessPage) KeyDown(ev *event.Event) {
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  
  if ev.Data["Key"] == UI.CurKeys["Right"] {
    self.MySlider.Further()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }

  if ev.Data["Key"] == UI.CurKeys["Left"] {
    self.MySlider.StepBack()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }  
}

func (self *BrightnessPage) Draw() {

  self.ClearCanvas()
  self.MySlider.Draw()

}

