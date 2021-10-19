package Warehouse

import (
	//"fmt"
	//"os"
	//gotime "time"
	
	//"github.com/cuu/gogame/image"
	//"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

type Aria2DownloadProcessPage struct {
	UI.Page
	URLColor  *color.Color
	TextColor *color.Color

	Icons          map[string]UI.IconItemInterface
	FileNameLabel UI.LabelInterface
	SizeLabel     UI.LabelInterface
	
}

func NewAria2DownloadProcessPage() *Aria2DownloadProcessPage {
	p := &Aria2DownloadProcessPage{}
	p.Icons = make(map[string]UI.IconItemInterfac)
	
	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")

	p.FootMsg = [5]string{"Nav.","","Pause","Back","Cancel"}

	return p
}

func (self *Aria2DownloadProcessPage) Init() {
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.CanvasHWND = self.Screen.CanvasHWND

	bgpng := UI.NewIconItem()
	bgpng.ImgSurf = UI.MyIconPool.GiveIconSurface("rom_download")
	bgpng.MyType = UI.ICON_TYPES["STAT"]
	bgpng.Parent = self
	bgpng.Adjust(0,0,UI.MyIconPool.Width("rom_download"),UI.MyIconPool.Height("rom_download"),0)
	self.Icons["bg"] = bgpng
	
}
