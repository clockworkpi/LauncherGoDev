package UI

import (
  "log"
  
	"github.com/itchyny/volume-go"
	"github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/rect"
  //"github.com/cuu/gogame/color"

)

type SoundPatch struct {
  AboveAllPatch
  
  snd_segs [][2]int
  Needle int
  Parent *MainScreen
}

func NewSoundPatch() *SoundPatch {
  p := &SoundPatch{}
  p.PosX = Width /2
  p.PosY = Height /2
  p.Width = 50
  p.Height = 120
  
  p.FontObj = Fonts["veramono20"]
  p.Color   = MySkinManager.GiveColor("Text")
  p.ValColor = MySkinManager.GiveColor("URL")
  
  p.Icons = make( map[string]IconItemInterface )
  
  p.Value = 0
   
  p.snd_segs = [][2]int{ [2]int{0,20},[2]int{21,40},[2]int{41,50},
                        [2]int{51,60},[2]int{61,70},[2]int{71,85},
                        [2]int{86,90},[2]int{91,95},[2]int{96,100}}
  
  
  return p
}


func (self *SoundPatch) Init() {
  self.SetCanvasHWND(self.Parent.CanvasHWND)
}

func (self *SoundPatch) VolumeUp() int {
  
  vol, err := volume.GetVolume()
  if err != nil {
    log.Printf("SoundPatch VolumeUp get volume failed: %+v", err)
		vol = 0
  }
  
	for i,v := range self.snd_segs {
		if vol >= v[0] && vol <= v[1] {
			self.Needle = i
			break
		}
	}
  
  self.Needle += 1
  
  if self.Needle > len(self.snd_segs) -1 {
    self.Needle = len(self.snd_segs) -1
  }
  
  val := self.snd_segs[self.Needle][0] +  (self.snd_segs[self.Needle][1] - self.snd_segs[self.Needle][0])/2
    
  volume.SetVolume(val)
  
  self.Value = self.snd_segs[self.Needle][1]
  
  self.Parent.TitleBar.SetSoundVolume(val)
  
  return self.Value
}

func (self *SoundPatch) VolumeDown() int {
  vol, err := volume.GetVolume()
  if err != nil {
    log.Printf("SoundPatch VolumeDown get volume failed: %+v\n", err)
		vol = 0
  }
  
	for i,v := range self.snd_segs {
		if vol >= v[0] && vol <= v[1] {
			self.Needle = i
			break
		}
	}
  
  self.Needle -= 1
  
  if self.Needle < 0 {
    self.Needle = 0
  }
  
  val := self.snd_segs[self.Needle][0]
  
  if val < 0 {
    val = 0
  }
  
  volume.SetVolume(val)
  
  self.Value = val
  
  self.Parent.TitleBar.SetSoundVolume(val)
  
  return self.Value

}

func (self *SoundPatch) Draw() {

  for i:=0;i< (self.Needle+1);i++ {
    vol_rect := rect.Rect(80+i*20, self.Height/2+20,10, 40)
    draw.AARoundRect(self.CanvasHWND,&vol_rect,MySkinManager.GiveColor("Front"),3,0,MySkinManager.GiveColor("Front"))
  }
}

