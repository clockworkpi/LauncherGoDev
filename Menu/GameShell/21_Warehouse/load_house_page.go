package Warehouse

import (
	"fmt"
	"os"
	"io/ioutil"
	gotime "time"
	"strings"
	"encoding/json"
	"path"
	"path/filepath"
	
	"github.com/veandco/go-sdl2/ttf"
	
	//"github.com/cuu/gogame/image"
	//"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/time"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/cuu/grab"
)

type WareHouseIndex struct {
    List []map[string]string `json:"list"`
}

type LoadHousePage struct {
	UI.Page
	ListFontObj *ttf.Font
	URLColor  *color.Color
	TextColor *color.Color

	Downloader *grab.Client
	resp       *grab.Response
	req        *grab.Request
	
	URL  string
	Downloading chan bool
	LoadingLabel UI.LabelInterface

	Parent *WareHouse
}

func NewLoadHousePage() *LoadHousePage {
	p := &LoadHousePage{}

	p.ListFontObj = UI.MyLangManager.TrFont("varela18")
	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")
	p.FootMsg = [5]string{"Nav.","","","Back","Cancel"}

	return p
}

func (self *LoadHousePage) Init() {
	
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.CanvasHWND = self.Screen.CanvasHWND
	self.LoadingLabel = UI.NewLabel()
	self.LoadingLabel.SetCanvasHWND(self.CanvasHWND)
	self.LoadingLabel.Init("Loading",self.ListFontObj,nil)
	self.LoadingLabel.SetColor(self.TextColor)

	self.Downloader = grab.NewClient()
	self.Downloading = make(chan bool,1)
	
}

func (self *LoadHousePage) SetDownloading(v bool) {
	for len(self.Downloading) > 0 {
		<- self.Downloading
	}

	self.Downloading <- v
}

func (self *LoadHousePage) OnLoadCb() {
	
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
	fmt.Println("LoadHouse OnLoadCb")
	if len(local_dir) > 1 {
		menu_file := local_dir[1]
		local_menu_file := fmt.Sprintf("%s/aria2downloads%s",
			home_path,menu_file)

		if UI.FileExists(local_menu_file) {
			var result WareHouseIndex
			jsonFile, err := os.Open(local_menu_file)
			if err != nil {
        fmt.Println(err)
				return
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &result)
			self.Parent.MyStack.Push(result.List)
			
			self.Leave()
		} else {
			
			self.req,_ = grab.NewRequest("/tmp",self.URL)
			self.resp = self.Downloader.Do(self.req)

			self.SetDownloading(true)
			fmt.Println("Start Download index json to /tmp,grab")
			go self.UpdateProcessInterval(400)
		}
		
	}
}

func (self *LoadHousePage) UpdateProcessInterval(ms int) {
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
			if(ms_total > 5000) {
				fmt.Println("LoadHouse Timeout")
				break L
			}
		case <-self.resp.Done:
			// download is complete
			fmt.Println("Grab Download House done")
			break L
		case v:= <- self.Downloading:
			if v == false {
				break L
			}
		}		
	}
	
	dst_filename := self.resp.Filename
	fmt.Println("dst_filename ",dst_filename)
	
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
		fmt.Println("dl_file: ",dl_file)
		fmt.Println(local_menu_file)

		
		if UI.IsDirectory( filepath.Dir(local_menu_file) ) == false {
			merr := os.MkdirAll( filepath.Dir(local_menu_file), os.ModePerm)
			if merr != nil {
				panic(merr)
			}
		}

		UI.CopyFile(dl_file,local_menu_file)
		
		os.Remove(dl_file)
		
		var result WareHouseIndex
		jsonFile, err := os.Open(local_menu_file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &result)
		self.Parent.MyStack.Push(result.List)		
		self.Leave()
		
	} else {
		fmt.Println(err)
		self.Screen.MsgBox.SetText("Fetch house failed")
		self.Screen.MsgBox.Draw()
		self.Screen.SwapAndShow()
		time.BlockDelay(500)
	}
	
}

func (self *LoadHousePage) Leave() {

	self.SetDownloading(false)
	
	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()
	self.URL = ""
	
}

func (self *LoadHousePage) KeyDown(ev *event.Event) {
	if UI.IsKeyMenuOrB(ev.Data["Key"]) {
		self.Leave()
	}
	
}

func (self *LoadHousePage) Draw() {
	self.ClearCanvas()
	w,_ := self.LoadingLabel.Size()
	self.LoadingLabel.NewCoord( (UI.Width - w)/2,(UI.Height-44)/2);
	self.LoadingLabel.Draw()
	
}
