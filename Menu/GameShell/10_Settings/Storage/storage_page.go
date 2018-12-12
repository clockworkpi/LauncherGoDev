package Storage

import (
	"fmt"
	"syscall"
  
  "github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/color"

  "github.com/cuu/LauncherGoDev/sysgo/UI"
  
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type StoragePage struct {
  UI.Page
  
  BGpng UI.IconItemInterface
  BGwidth int // 96
  BGheight int // 73
  BGlabel UI.LabelInterface
  FreeLabel UI.LabelInterface
  
  BGmsg  string // "%.1GB of %.1fGB Used"
  
  DskUsg [2]float64
  HighColor *color.Color //MySkinManager.GiveColor('High')
  
}

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

func NewStoragePage() *StoragePage {
  p:= &StoragePage{}
  
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]
	
	p.FootMsg = [5]string{"Nav","","","Back",""}

  p.HighColor = &color.Color{51, 166, 255,255}
  
  p.BGwidth = 96
  p.BGheight = 73
  
  p.BGmsg = "%.1fGB of %.1fGB Used"
  return p
}

func (self *StoragePage) DiskUsage() (float64,float64) {
  disk := DiskUsage("/")
  
  all := float64(disk.All)/float64(GB)
  
  free := float64(disk.Free)/float64(GB)
  
  return free,all

}

func (self *StoragePage) Init() {
  
  self.DskUsg[0],self.DskUsg[1] = self.DiskUsage()
  
  self.CanvasHWND = self.Screen.CanvasHWND
  self.Width      = self.Screen.Width
  self.Height     = self.Screen.Height
  
  bgpng := UI.NewIconItem()
  bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("icon_sd")
  bgpng.MyType = UI.ICON_TYPES["STAT"]
  bgpng.Parent = self
  
  bgpng.AddLabel( fmt.Sprintf(self.BGmsg,self.DskUsg[1]-self.DskUsg[0],self.DskUsg[1]),UI.Fonts["varela15"])
  bgpng.Adjust(0,0,self.BGwidth,self.BGheight,0)
  
  self.BGpng = bgpng
  self.BGlabel = UI.NewLabel()
  self.BGlabel.SetCanvasHWND(self.CanvasHWND)
  
  usage_percent := int((self.DskUsg[0]/self.DskUsg[1])*100.0)
  
  self.BGlabel.Init(fmt.Sprintf("%d%%",usage_percent ),UI.Fonts["varela25"],nil)
  self.BGlabel.SetColor(self.HighColor)
  
  self.FreeLabel = UI.NewLabel()
  self.FreeLabel.SetCanvasHWND(self.CanvasHWND)
  self.FreeLabel.Init("Free",UI.Fonts["varela13"],nil)
  self.FreeLabel.SetColor(self.BGlabel.(*UI.Label).Color)

}

func (self *StoragePage) Draw() {
  self.ClearCanvas()
  
  self.BGpng.NewCoord(self.Width/2,self.Height/2-10)
  self.BGpng.Draw()
  
  self.BGlabel.NewCoord(self.Width/2-28,self.Height/2-30)
  self.BGlabel.Draw()

  x,_ := self.BGlabel.Coord()
  self.FreeLabel.NewCoord(x+10   ,self.Height/2)
  self.FreeLabel.Draw()

  usage_percent := (self.DskUsg[0]/self.DskUsg[1] )
  if usage_percent < 0.1 {
    usage_percent = 0.1
  }
  
  rect_ := draw.MidRect(self.Width/2,self.Height-30,170,17, UI.Width,UI.Height)

  draw.AARoundRect(self.CanvasHWND,rect_,&color.Color{169,169,169,255},5,0,&color.Color{169,169,169,255})
  
  rect2_ := draw.MidRect(self.Width/2,self.Height-30,int(170.0*(1.0-usage_percent)),17, UI.Width,UI.Height)
  
  rect2_.X = rect_.X
  rect2_.Y = rect_.Y
  
  
  draw.AARoundRect(self.CanvasHWND,rect2_,&color.Color{131,199,219,255},5,0,&color.Color{131,199,219,255})
  

}
