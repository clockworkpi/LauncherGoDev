package Sound

import (
	"fmt"
	"github.com/cuu/gogame/event"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"

	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

type OnChangeCB_T func(int)

type SliderIcon struct {
	UI.IconItem
	Parent *SoundSlider
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
	parent_x, parent_y := self.Parent.Coord()

	if self.Label != nil {
		//		lab_x,lab_y := self.Label.Coord()
		lab_w, lab_h := self.Label.Size()

		if self.Align == UI.ALIGN["VCenter"] {
			//			fmt.Println("IconItem Draw VCenter:",lab_w,lab_h,self.Label.GetText())

			self.Label.NewCoord(self.PosX-lab_w/2+parent_x, self.PosY+self.Height/2+6+parent_y)

		} else if self.Align == UI.ALIGN["HLeft"] {
			self.Label.NewCoord(self.PosX+self.Width/2+3+parent_x, self.PosY-lab_h/2+parent_y)
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		surface.Blit(self.Parent.GetCanvasHWND(), self.ImgSurf, draw.MidRect(self.PosX+parent_x, self.PosY+parent_y,
			self.Width, self.Height, UI.Width, UI.Height), nil)
	}
}

type SliderMultiIcon struct {
	UI.MultiIconItem
	Parent *SoundSlider
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
	parent_x, parent_y := self.Parent.Coord()

	if self.Label != nil {
		//		lab_x,lab_y := self.Label.Coord()
		lab_w, lab_h := self.Label.Size()
		if self.Align == UI.ALIGN["VCenter"] {
			self.Label.NewCoord(self.PosX-lab_w/2+parent_x, self.PosY+self.Height/2+6+parent_y)
		} else if self.Align == UI.ALIGN["HLeft"] {
			self.Label.NewCoord(self.PosX+self.Width/2+3+parent_x, self.PosY-lab_h/2+parent_y)
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {

		portion := rect.Rect(0, self.IconIndex*self.IconHeight, self.IconWidth, self.IconHeight)

		surface.Blit(self.Parent.GetCanvasHWND(),
			self.ImgSurf, draw.MidRect(self.PosX+parent_x, self.PosY+parent_y,
				self.Width, self.Height, UI.Width, UI.Height), &portion)
	}
}

type SoundSlider struct {
	UI.Slider

	BGpng    *SliderIcon
	BGwidth  int
	BGheight int
	//NeedleSurf
	Scale  *SliderMultiIcon
	Parent *SoundPage

	OnChangeCB OnChangeCB_T

	snd_segs [][2]int
}

func NewSoundSlider() *SoundSlider {
	p := &SoundSlider{}
	p.Range = [2]int{0, 255}
	p.Value = 0

	p.BGwidth = 192
	p.BGheight = 173

	p.snd_segs = [][2]int{[2]int{0, 20}, [2]int{21, 40}, [2]int{41, 50},
		[2]int{51, 60}, [2]int{61, 70}, [2]int{71, 85},
		[2]int{86, 90}, [2]int{91, 95}, [2]int{96, 100}}

	return p
}

func (self *SoundSlider) GetCanvasHWND() *sdl.Surface {
	return self.CanvasHWND
}

func (self *SoundSlider) Init() {
	self.Width = self.Parent.Width
	self.Height = self.Parent.Height

	self.BGpng = NewSliderIcon()
	self.BGpng.ImgSurf = UI.MyIconPool.GetImgSurf("vol")
	self.BGpng.MyType = UI.ICON_TYPES["STAT"]
	self.BGpng.Parent = self
	self.BGpng.Adjust(0, 0, self.BGwidth, self.BGheight, 0)

	self.Scale = NewSliderMultiIcon()
	self.Scale.MyType = UI.ICON_TYPES["STAT"]
	self.Scale.Parent = self
	self.Scale.ImgSurf = UI.MyIconPool.GetImgSurf("scale")
	self.Scale.IconWidth = 82
	self.Scale.IconHeight = 63
	self.Scale.Adjust(0, 0, 82, 63, 0)
}

func (self *SoundSlider) SetValue(vol int) { // pct 0 - 100
	for i, v := range self.snd_segs {
		if vol >= v[0] && vol <= v[1] {
			self.Value = i
			break
		}
	}
}

func (self *SoundSlider) Further() {
	self.Value += 1

	if self.Value >= len(self.snd_segs)-1 {
		self.Value = len(self.snd_segs) - 1
	}

	vol := self.snd_segs[self.Value][0] + (self.snd_segs[self.Value][1]-self.snd_segs[self.Value][0])/2

	if self.OnChangeCB != nil {
		self.OnChangeCB(vol)
	}
}

func (self *SoundSlider) StepBack() {
	self.Value -= 1

	if self.Value < 0 {
		self.Value = 0
	}

	vol := self.snd_segs[self.Value][0] + (self.snd_segs[self.Value][1]-self.snd_segs[self.Value][0])/2

	if self.OnChangeCB != nil {
		self.OnChangeCB(vol)
	}
}

func (self *SoundSlider) Draw() {
	self.BGpng.NewCoord(self.Width/2, self.Height/2)
	//fmt.Printf("%x\n",self.BGpng.Parent)
	self.BGpng.Draw()

	self.Scale.NewCoord(self.Width/2, self.Height/2)

	self.Scale.IconIndex = self.Value

	self.Scale.Draw()

}

type SoundPage struct {
	UI.Page

	MySlider *SoundSlider
}

func NewSoundPage() *SoundPage {
	p := &SoundPage{}

	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10
	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav", "", "", "Back", "Enter"}

	return p
}

func (self *SoundPage) Init() {
	self.CanvasHWND = self.Screen.CanvasHWND
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.MySlider = NewSoundSlider()

	self.MySlider.Parent = self
	self.MySlider.SetCanvasHWND(self.CanvasHWND)

	self.MySlider.OnChangeCB = self.WhenSliderDrag

	self.MySlider.Init()

	v, err := GetVolume()
	if err == nil {
		self.MySlider.SetValue(v)
	} else {
		fmt.Println(err)
	}
}

func (self *SoundPage) OnLoadCb() {
	v, err := GetVolume()
	if err == nil {
		self.MySlider.SetValue(v)
	} else {
		fmt.Println(err)
	}
}

func (self *SoundPage) WhenSliderDrag(val int) { //value 0 - 100
	if val < 0 || val > 100 {
		return
	}

	self.Screen.TitleBar.SetSoundVolume(val)

	SetVolume(val)
}

func (self *SoundPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Right"] {
		self.MySlider.Further()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Left"] {
		self.MySlider.StepBack()
		self.Screen.Refresh()
	}

}

func (self *SoundPage) Draw() {
	self.ClearCanvas()
	self.MySlider.Draw()
}
