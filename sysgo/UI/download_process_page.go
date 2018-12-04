package UI

import (
  gotime "time"
  "net/url"
  
  "github.com/cuu/grab"
  "github.com/cuu/gogame/color"
)

type DownloadProcessPage struct {
  UI.Page
    
  URL string
  DST_DIR string
  Value   int
  PngSize map[string][2]int
  
  FileNameLabel LabelInterface
  SizeLabel     LabelInterface
  
  Icons  map[string]IconItemInterface

  URLColor  *color.Color
  TextColor  *color.Color
  TheTicker *gotime.Ticker
  
  Downloader *grab.Client
  resp *grab.Response
  req *grab.Request
  
}


func NewDownloadProcessPage() *DownloadProcessPage {
  
  p := &DownloadProcessPage{}
  
  p.FootMsg = [5]string{"Nav","","","Back",""}
  
  p.URLColor = &color.Color{51, 166, 255,255 } // URL
  p.TextColor = &color.Color{83,83,83,255 } // Text
  
  p.PngSize = make(map[string][2]int,0)
  
  p.Icons=make(map[string]IconItemInterface)
  
  return p
}

func (self *DownloadProcessPage) Init() {
  self.PosX = self.Index * self.Screen.Width
  self.Width  = self.Screen.Width
  self.Height = self.Screen.Height
  
  self.CanvasHWND = self.Screen.CanvasHWND
  self.PngSize["bg"] = [2]int{48,79}
  self.PngSize["needwifi_bg"] = [2]int{253,132}
  
  bgpng := NewIconItem()
  bgpng.ImgSurf = MyIconPool.GetImgSurf("rom_download")
  bgpng.MyType = ICON_TYPES["STAT"]
  bgpng.Parent = self
  bgpng.Adjust(0,0,self.PngSize["bg"][0],self.PngSize["bg"][1],0)
  self.Icons["bg"] = bgpng
  
  needwifi_bg = NewIconItem()
  needwifi_bg.ImgSurf = MyIconPool.GetImgSurf("needwifi_bg")
  needwifi_bg.MyType = ICON_TYPES["STAT"]
  needwifi_bg.Parent = self
  needwifi_bg.Adjust(0,0,self.PngSize["needwifi_bg"][0],self.PngSize["needwifi_bg"][1],0)
  
  self.Icons["needwifi_bg"] = needwifi_bg

  self.FileNameLabel = NewLabel()
  self.FileNameLabel.SetCanvasHWND(self.CanvasHWND)
  self.FileNameLabel.Init("", Fonts["varela12"])
  
  self.SizeLabel = NewLabel()
  self.SizeLabel.SetCanvasHWND(self.CanvasHWND)
  self.SizeLabel.Init("0/0Kb",Fonts["varela12"])
  self.SizeLabel.SetColor( self.URLColor )
  
  self.Downloader = grab.NewClient()
  
}

func (self *DownloadProcessPage) OnExitCb() {
  
  //Stop Ticker and the Grab
  if self.TheTicker != nil {
    self.TheTicker.Stop()
  }
  
}

// should be in a gorotine
func (self *DownloadProcessPage) UpdateProcessInterval() {
  if self.TheTicker == nil {
    return
  }
  
  if self.Screen.CurPage() != self {
    return
  }
  
  
  for {
		select {
		case <-self.TheTicker.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				self.resp.BytesComplete(),
				self.resp.Size,
				100*self.resp.Progress())
    self.Value = int(100*self.resp.Progress())
    
		case <-self.resp.Done:
			// download is complete
      fmt.Println("download is complete ",self.Value)
      self.Value = 0 
      
      
			break
		}
  }
  
	if err := self.resp.Err(); err != nil {
    self.DownloadErr()
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
	}

  fmt.Printf("Download saved to %s/%v \n",self.DST_DIR, self.resp.Filename)
  
  filename := filepath.Base(self.resp.Filename)
  
  if strings.HasSuffix(filename,".zip") {
    cmd := exec.Command("unzip",filename)
    cmd.Dir = self.DST_DIR
    cmd.Run()
  }else if strings.HasSuffix(filename,".zsync") {
    cmd := exec.Command("rm","-rf",filename)
    cmd.Dir = self.DST_DIR
    cmd.Run()
  }else if strings.HasSuffix(filename,".tar.gz") {
    cmd := exec.Command("tar", "xf", filename)
    cmd.Dir= self.DST_DIR
    cmd.Run()
  }
  
  cmd := exec.Command("rm","-rf",filename)
  cmd.Dir = self.DST_DIR
  cmd.Run()
    
  self.TheTicker.Stop()
  
  self.DoneAndReturnUpLevel()
  
}

func (self *DownloadProcessPage) DownloadErr()  {
  self.Screen.MsgBox.SetText("Download Failed")
  self.Screen.MsgBox.Draw()
  self.Screen.SwapAndShow()  
}

func (self *DownloadProcessPage) DoneAndReturnUpLevel() {
  self.ReturnToUpLevelPage()
  self.Screen.Draw()
  self.Screen.SwapAndShow()
}



func (self *DownloadProcessPage) StartDownload(_url,dst_dir string) {
  
  if self.Screen.DBusManager.IsWifiConnectedNow() == false {
    return
  }
  
  _, err := url.ParseRequestURI(_url)
  if err == nil && UI.IsDirectory(dst_dir) {
    self.URL = _url
    self.DST_DIR = dst_dir
  }else{
  
    self.Screen.MsgBox.SetText("Invaid")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow() 
    fmt.Println("DownloadProcessPage StartDownload Invalid ",err)
    return
  }

  self.req, _ := grab.NewRequest(self.DST_DIR, _url)
  
  fmt.Printf("Downloading %v...\n", self.req.URL())
  
  self.resp = self.Downloader.Do(self.req)
  
  fmt.Printf("  %v\n", self.resp.HTTPResponse.Status)
  
  self.TheTicker = gotime.NewTicker(100 * gotime.Millisecond)
  
  go self.UpdateProcessInterval()
  
}


func (self *DownloadProcessPage) Draw() {

  self.ClearCanvas()
  
  if self.Screen.DBusManager.IsWifiConnectedNow() == false {
    self.Icons["needwifi_bg"].NewCoord(self.Width/2,self.Height/2)
    self.Icons["needwifi_bg"].Draw()
    return
    
  }
  
  self.Icons["bg"].NewCoord(self.Width/2,self.Height/2)
  self.Icons["bg"].Draw()
  
  percent := self.Value
  if percent < 10 {
    percent = 10
  }
  
  rect_ := draw.MidRect(self.Width/2,self.Height/2+33,170,17, UI.Width,UI.Height)
  
  draw.AARoundRect(self.CanvasHWND,rect_,
                  &color.Color{228,228,228,255},5,0,&color.Color{228,228,228,255})
  
  
  rect2_ := draw.MidRect( self.Width/2,self.Height/2+33,int(170.0*((float64)percent/100.0)),17, UI.Width,UI.Height )
  
  rect2_.X = rect_.X
  rect2_.Y = rect_.Y
  
  draw.AARoundRect(self.CanvasHWND,rect2_,
                  &color.Color{131, 199, 219,255},5,0,&color.Color{131, 199, 219,255})
  
  w,h: = self.FileNameLabel.Size()
  
  rect3_ := draw.MidRect(self.Width/2,self.Height/2+53,w, h,UI.Width,UI.Height)

  w, h = self.SizeLabel.Size()
  
  rect4 := draw.MidRect(self.Width/2,self.Height/2+70,w, h,UI.Width,UI.Height)
  
  self.FileNameLabel.NewCoord(int(rect3_.X),int(rect3_.Y))
  self.FileNameLabel.Draw()
  
  self.SizeLabel.NewCoord(int(rect4_.X),int(rect4_.Y))
  self.SizeLabel.Draw()
  
}
