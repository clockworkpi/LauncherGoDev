package Warehouse

import (
	"fmt"
	"os"
	gotime "time"
	"encoding/json"

	"github.com/veandco/go-sdl2/ttf"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/cuu/gogame/image"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	"github.com/cuu/grab"
)

type WareHouseIndex struct {
    List []map[string]string `json:"list"`
}


type ImageDownloadProcessPage struct {
	UI.Page
	
	ListFontObj *ttf.Font
	
	URLColor  *color.Color
	TextColor *color.Color
	
	Downloader *grab.Client
	resp       *grab.Response
	req        *grab.Request
	URL        string
	Value      int
	LoadingLabel *UI.LabelInterface

	Img        *sdl.Surface
	Downloading chan bool
	Parent *WareHouse
	
}


func NewImageDownloadProcessPage() *ImageDownloadProcessPage {
	p := &ImageDownloadProcessPage{}
	p.ListFontObj = UI.MyLangManager.TrFont("varela13")
	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")
	p.FootMsg = [5]string{"Nav.","","","Back",""}

	return p
}

func (self *ImageDownloadProcessPage) Init() {
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.CanvasHWND = self.Screen.CanvasHWND
	self.LoadingLabel = UI.NewLabel()
	self.LoadingLabel.SetCanvasHWND(self.CanvasHWND)
	self.LoadingLabel.Init("Loading",self.ListFontObj,nil)
	self.LoadingLabel.SetColor(self.TextColor)

	self.Downloader = grab.NewClient()
	self.Downloading = make(chan bool)
}


func (self *ImageDownloadProcessPage) OnLoadCb() {

	if len(self.URL) < 10 {
		return
	}

	self.ClearCanvas()
	self.Screen.Draw()
	self.Screen.SwapAndShow()

	parts := strings.Split(self.URL,"/")
	filename := strings.TrimSpace(parts[len(parts)-1])
	local_dir := strings.Split(self.URL,"raw.githubusercontent.com")
	home_path, _ := os.UserHomeDir()
	
	if len(local_dir) > 1 {
		menu_file := local_dir[1]
		local_menu_file := fmt.Sprintf("%s/aria2downloads%s",
			home_path,menu_file)

		if UI.FileExists(local_menu_file) {
			self.Img = image.Load(local_menu_file)
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}else {
			
			self.req,_ = grab.NewRequest("/tmp",self.URL)
			self.resp = self.Downloader.Do(self.req)
			for len(self.Downloading) > 0 {
				<-self.Downloading
			}
			self.Downloading <- true
			
			go self.UpdateProcessInterval(400)
			
		}
	}
}

func (self *ImageDownloadProcessPage) UpdateProcessInterval(ms int) {
	
	t := gotime.NewTicker(ms * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				self.resp.BytesComplete(),
				self.resp.Size,
				100*self.resp.Progress())

		case <-self.resp.Done:
			// download is complete
			break
		case v:= <-self.Downloading
			if v == false {
				t.Stop()
				break
			}
		}
	}

	dst_filename := self.resp.Filename
	
	if err := self.resp.Err(); err == nil {//download successfully
		home_path, _ := os.UserHomeDir()
		parts := strings.Split(self.URL,"/")
		filename := strings.TrimSpace(parts[len(parts)-1])
		local_dir := strings.Split(self.URL,"raw.githubusercontent.com")

		local_menu_file := ""
		menu_file := ""
		
		if len(local_dir) > 1 {
			menu_file = local_dir[1]
			local_menu_file = fmt.Sprintf("%s/aria2downloads%s",
				home_path,menu_file)
		}

		dl_file := path.Join("/tmp",filename)
		if UI.IsDirectory( Path.Base(local_menu_file) ) == false {
			merr := os.MkdirAll( Path.Base(local_menu_file), os.ModePerm)
			if merr != nil {
				panic(merr)
			}
		}

		UI.CopyFile(dl_file,local_menu_file)
			
	}

	if UI.FileExists(dst_filename) {
		if self.Screen.CurPage() == self {
			self.Img = image.Load(dst_filename)
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}
	}
	
	
}


func (self *ImageDownloadProcessPage) KeyDown(ev  *event.Event) {

	if IsKeyMenuOrB(ev.Data["Key")) {

		self.Downloading <- false

		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
		self.URL = ""
	}
}

func (self *ImageDownloadProcessPage) Draw() {
	self.ClearCanvas()
	self.LoadingLabel.NewCoord( (UI.Width - self.LoadingLabel.Width)/2,(UI.Height-44)/2);
	self.LoadingLabel.Draw()
	if self.Img != nil {
		self.CanvasHWND.Blit(self.Img,draw.MidRect(UI.Width/2,(UI.Height-44)/2,
			self.Img.Width,self.Img.Height,
			UI.Width,UI.Height-44))
	}
}

