package Warehouse

import (
	"fmt"
	"os"
	gotime "time"
	"strings"
	"path"
	"path/filepath"
	//"encoding/json"

	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/sdl"
	
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/cuu/gogame/image"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/surface"

	"github.com/cuu/grab"
)

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
	LoadingLabel UI.LabelInterface

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
	LoadingLabel := UI.NewLabel()
	LoadingLabel.SetCanvasHWND(self.CanvasHWND)
	LoadingLabel.Init("Loading",self.ListFontObj,nil)
	LoadingLabel.SetColor(self.TextColor)
	self.LoadingLabel = LoadingLabel
	
	self.Downloader = grab.NewClient()
	self.Downloading = make(chan bool,1)
}

func (self *ImageDownloadProcessPage) SetDownloading(v bool) {
	for len(self.Downloading) > 0 {
		<- self.Downloading
	}

	self.Downloading <- v
}

func (self *ImageDownloadProcessPage) OnLoadCb() {

	if len(self.URL) < 10 {
		return
	}

	self.ClearCanvas()
	self.Screen.Draw()
	self.Screen.SwapAndShow()

	//parts := strings.Split(self.URL,"/")
	//filename := strings.TrimSpace(parts[len(parts)-1])
	local_dir := strings.Split(self.URL,"raw.githubusercontent.com")
	home_path, _ := os.UserHomeDir()
	
	if len(local_dir) > 1 {
		menu_file := local_dir[1]
		local_menu_file := fmt.Sprintf(aria2dl_folder,
			home_path,menu_file)

		if UI.FileExists(local_menu_file) {
			self.Img = image.Load(local_menu_file)
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}else {
			
			self.req,_ = grab.NewRequest("/tmp",self.URL)
			self.resp = self.Downloader.Do(self.req)
			
			self.SetDownloading(true)
			
			go self.UpdateProcessInterval(400)
			
		}
	}
}

func (self *ImageDownloadProcessPage) UpdateProcessInterval(ms int) {
	ms_total := 0
	t := gotime.NewTicker(gotime.Duration(ms) * gotime.Millisecond)
	defer t.Stop()
L:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				self.resp.BytesComplete(),
				self.resp.Size,
				100*self.resp.Progress())
			ms_total += ms
			if(ms_total > 10000) {
				fmt.Println("Get preview image timeout")
				break L
			}
		case <-self.resp.Done:
			// download is complete
			break L
		case v:= <-self.Downloading:
			if v == false {
				break L
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
			local_menu_file = fmt.Sprintf(aria2dl_folder,
				home_path,menu_file)
		}

		dl_file := path.Join("/tmp",filename)
		if UI.IsDirectory( filepath.Dir(local_menu_file) ) == false {
			merr := os.MkdirAll( filepath.Dir(local_menu_file), os.ModePerm)
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

	if UI.IsKeyMenuOrB(ev.Data["Key"]) {

		self.SetDownloading(false)

		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
		self.URL = ""
	}
}

func (self *ImageDownloadProcessPage) Draw() {
	self.ClearCanvas()
	w,_ := self.LoadingLabel.Size()
	self.LoadingLabel.NewCoord( (UI.Width - w)/2,(UI.Height-44)/2);
	self.LoadingLabel.Draw()
	if self.Img != nil {
		surface.Blit(self.CanvasHWND,
			self.Img,
			draw.MidRect(UI.Width/2,(UI.Height-44)/2,int(self.Img.W),int(self.Img.H),UI.Width,UI.Height-44),
			nil)
	}
}

