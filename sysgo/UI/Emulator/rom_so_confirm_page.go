package Emulator

import (
  "strconv"
  "strings"
  
  "path/filepath"
  
	"github.com/cuu/gogame/surface"
  "github.com/cuu/LauncherGoDev/sysgo"
  "github.com/cuu/LauncherGoDev/sysgo/UI"

)

type RomSoConfirmPage struct {
  UI.ConfirmPage
  
  DownloadPage *UI.DownloadProcessPage
  
}

func NewRomSoConfirmPage() *RomSoConfirmPage {
  p := &RomSoConfirmPage{}
  p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10
	p.Align = ALIGN["SLeft"]
  
  p.ListFont = UI.Fonts["veramono18"]
  p.FootMsg = [5]string{"Nav","","","Cancel","Yes"}
  p.ConfirmText ="Do you want to setup this game engine automatically?"
  
  return p
  
}

func (self *RomSoConfirmPage) Init() {
  self.PosX = self.Index * self.Screen.Width
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height

  self.CanvasHWND = self.Screen.CanvasHWND
  
  li := UI.NewMultiLabel()
  li.SetCanvasHWND(self.CanvasHWND)
  li.Width = 160
  li.Init(self.ConfirmText,self.ListFont)
  
  li.PosX = (self.Width - li.Width)/2
  li.PosY = (self.Height - li.Height)/2

  self.BGPosX = li.PosX-20
  self.BGPosY = li.PosY-20
  self.BGWidth = li.Width+40
  self.BGHeight = li.Height+40  
  
  self.MyList = append(self.MyList ,li )
  
}

func (self *RomSoConfirmPage) SnapMsg(msg string) {
  self.MyList[0].SetText(msg)
  self.Screen.Draw()
  self.Screen.SwapAndShow()
  self.MyList[0].SetText(self.ConfirmText)
}

func (self *RomSoConfirmPage) OnReturnBackCb() {
  self.ReturnToUpLevelPage()
  self.Screen.Draw()
  self.Screen.SwapAndShow()
}

func (self *RomSoConfirmPage) KeyDown(ev *event.Event) {
  
  if ev.Data["Key"] == UI.CurKeys["Menu"] || ev.Data["Key"] == UI.CurKeys["A"] {
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  if ev.Data["Key"] == UI.CurKeys["B"] {
    if UI.CheckBattery < 5 {
      self.SnapMsg("Battery must over 5%")
    }else {
      if self.DownloadPage == nil {
        self.DownloadPage = UI.NewDownloadProcessPage()
        self.DownloadPage.Screen = self.Screen
        self.DownloadPage.Name   = "Downloading"
        self.DownloadPage.Init()
      }
      
      self.Screen.PushPage(self.DownloadPage)
      self.Screen.Draw()
      self.Screen.SwapAndShow()
      
      if sysgo.CurKeySet == "PC" {
        so_url := self.Parent.EmulatorConfig.SO_URL
        so_url = strings.Replace(so_url,"armhf","x86_64",-1)
        fmt.Println(so_url)
        self.DownloadPage.StartDownload(so_url,filepath.Dir(self.Parent.EmulatorConfig.ROM_SO))
        
      }else{
        so_url := self.Parent.EmulatorConfig.SO_URL
        self.DownloadPage.StartDownload(so_url,filepath.Dir(self.Parent.EmulatorConfig.ROM_SO))
      }
    }
  }
}

func (self *RomSoConfirmPage) Draw() {
  self.ClearCanvas()
  self.DrawBG()
  for _,v := range self.MyList{
    v.Draw()
  }  
    

}
