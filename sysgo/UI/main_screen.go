package UI

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	//"encoding/json"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"path/filepath"
	gotime "time"

	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/display"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/time"

	"github.com/cuu/gogame/event"

	"github.com/clockworkpi/LauncherGoDev/sysgo"
)

//eg: MainScreen
type ScreenInterface interface {
	AppendPage(pg PageInterface)
	ClearCanvas()
	CurPage() PageInterface
	Draw()
	ExtraName(name string) string
	FartherPages()
	Init()
	IsEmulatorPackage(dirname string) bool
	IsExecPackage(dirname string) bool
	IsPluginPackage(dirname string) bool
	KeyDown(ev *event.Event)
	OnExitCb()
	HookExitCb()
	PushCurPage()
	PushPage(pg PageInterface)
	RunEXE(cmdpath string)
	SetCurPage(pg PageInterface)
	SwapAndShow()
	IsWifiConnectedNow()
	Refresh()
}

type PluginConfig struct {
	NAME    string `json:"NAME"` // plugin name,default could be the same as Plugin Folder's name
	SO_FILE string `json:"SO_FILE"`
}

type MessageBox struct {
	Label
	Parent *sdl.Surface
	ParentWidth int
	ParentHeight int
	HWND   *sdl.Surface
}

func NewMessageBox() *MessageBox {
	m := &MessageBox{}
	m.Color = &color.Color{83, 83, 83, 255}

	return m
}

func (self *MessageBox) Init(text string, font_obj *ttf.Font, col *color.Color,w int, h int) {
	if col != nil {
		self.Color = col
	}

	self.Text = text
	self.FontObj = font_obj

	self.Width = 0
	self.Height = 0
	
	if w == 0 || h == 0 {
		self.CanvasHWND = surface.Surface(Width,Height)
		self.ParentWidth = Width
		self.ParentHeight = Height
	}else{
		self.CanvasHWND = surface.Surface(w,h)
		self.ParentWidth = w
		self.ParentHeight = h
	}
	self.HWND = self.Parent

}

func (self *MessageBox) SetText(text string) {
	self.Text = text
}

func (self *MessageBox) Draw() {
	self.Width = 0
	self.Height = 0
	surface.Fill(self.CanvasHWND, &color.Color{255, 255, 255, 255})

	words := strings.Split(self.Text, " ")
	space, _ := font.Size(self.FontObj, " ")

	max_width := self.ParentWidth - 40
	x := 0
	y := 0

	row_total_width := 0
	lines := 0

	for _, word := range words {
		word_surface := font.Render(self.FontObj, word, true, self.Color, nil)
		word_width := int(word_surface.W)
		word_height := int(word_surface.H)
		row_total_width += word_width
		if lines == 0 {
			lines += word_height
		}

		if (row_total_width + space) >= max_width {
			x = 0
			y += word_height
			row_total_width = word_width
			lines += word_height
		}

		dest_rect := rect.Rect(x, y, word_width, word_height)
		surface.Blit(self.CanvasHWND, word_surface, &dest_rect, nil)
		word_surface.Free()
		if len(words) == 1 {
			x += word_width
		} else {
			x += word_width + space
		}

		if x > self.Width {
			self.Width = x
		}

		if lines >= self.ParentHeight-40 {
			break
		}
	}

	self.Height = lines

	padding := 5
	x = (self.ParentWidth - self.Width) / 2
	y = (self.ParentHeight - self.Height) / 2

	rect_ := rect.Rect(x-padding, y-padding, self.Width+padding*2, self.Height+padding*2)

	if self.HWND != nil {

		draw.Rect(self.HWND, &color.Color{255, 255, 255, 255}, &rect_, 0)

		rect__ := draw.MidRect(self.ParentWidth/2, self.ParentHeight/2, self.Width, self.Height, Width, Height)

		dest_rect := rect.Rect(0, 0, self.Width, self.Height)

		surface.Blit(self.HWND, self.CanvasHWND, rect__, &dest_rect)

		draw.Rect(self.HWND, &color.Color{0, 0, 0, 255}, &rect_, 1)

	}

}

type MainScreen struct {
	Widget
	Pages     []PageInterface
	Child     []PageInterface

	PageMax   int
	PageIndex int

	MyPageStack *PageStack

	CurrentPage PageInterface
	CanvasHWND  *sdl.Surface
	HWND        *sdl.Surface
	TitleBar    *TitleBar
	FootBar     *FootBar
	MsgBox      *MessageBox
	MsgBoxFont  *ttf.Font
	IconFont    *ttf.Font
	SkinManager *SkinManager

	CounterScreen *CounterScreen
	Closed        bool

	UIPluginList []*UIPlugin

	LastKey     string
	LastKeyDown gotime.Time
	
	updateScreen chan bool
}

func NewMainScreen() *MainScreen {
	m := &MainScreen{}

	m.PosY = TitleBar_BarHeight + 1
	m.Width = Width
	m.Height = Height - FootBar_BarHeight - TitleBar_BarHeight - 1
	m.MyPageStack = NewPageStack()

	m.MsgBoxFont = Fonts["veramono20"]
	m.IconFont = Fonts["varela15"]
	m.Closed = false
	
	m.updateScreen = make(chan bool,1)

	return m
}

func (self *MainScreen) Init() {
	self.CanvasHWND = surface.Surface(self.Width, self.Height)

	self.MsgBox = NewMessageBox()
	self.MsgBox.Parent = self.CanvasHWND
	self.MsgBox.Init(" ", self.MsgBoxFont, &color.Color{83, 83, 83, 255},self.Width,self.Height)

	self.SkinManager = NewSkinManager()
	self.SkinManager.Init()

	self.CounterScreen = NewCounterScreen()
	self.CounterScreen.HWND = self.HWND
	self.CounterScreen.Init()

	//self.GenList() // load predefined plugin list,ready to be injected ,or ,as a .so for dynamic loading
	go func() {
		self.RefreshLoop() 
	}()
}

func (self *MainScreen) FartherPages() { // right after ReadTheDirIntoPages
	self.PageMax = len(self.Pages)

	for i := 0; i < self.PageMax; i++ {
		self.Pages[i].SetIndex(i)
		self.Pages[i].SetCanvasHWND(self.CanvasHWND)
		self.Pages[i].UpdateIconNumbers() // IconNumbers always == len(Pages[i].Icons)
		self.Pages[i].SetScreen(self)
		self.Pages[i].Adjust()

		if self.Pages[i].GetIconNumbers() > 1 {
			self.Pages[i].SetPsIndex(1)
			self.Pages[i].SetIconIndex(1)
		}
	}

	self.CurrentPage = self.Pages[self.PageIndex]
	self.CurrentPage.SetOnShow(true)
}

func (self *MainScreen) CurPage() PageInterface {
	return self.CurrentPage
}

func (self *MainScreen) PushCurPage() {
	self.MyPageStack.Push(self.CurrentPage)
}

func (self *MainScreen) SetCurPage(pg PageInterface) {
	self.CurrentPage = pg
	pg.OnLoadCb()
}

func (self *MainScreen) PushPage(pg PageInterface) {
	self.PushCurPage()
	self.SetCurPage(pg)
}

func (self *MainScreen) AppendPage(pg PageInterface) {
	self.Pages = append(self.Pages, pg)
}

func (self *MainScreen) ClearCanvas() {
	surface.Fill(self.CanvasHWND, &color.Color{255, 255, 255, 255})
}

func (self *MainScreen) SwapAndShow() {
	if self.HWND != nil {
		rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
		surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
	}


	display.Flip()
}

func (self *MainScreen) ExtraName(name string) string {

	parts := strings.Split(name, "_")
	if len(parts) > 1 {
		return parts[1]
	} else if len(parts) == 1 {
		return parts[0]
	} else {
		return name
	}
}

//ExecPackage is all-in-one folder ,Name.sh,Name.png,etc
func (self *MainScreen) IsExecPackage(dirname string) bool {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println(err)
		return false
	}

	bname := filepath.Base(dirname)
	bname = self.ExtraName(bname)

	for _, v := range files {
		if v.Name() == bname+".sh" {
			return true
		}
	}

	return false
}

func (self *MainScreen) IsPluginPackage(dirname string) bool {
	ret := false
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, f := range files {
		if f.IsDir() {
			//pass
		} else {
			if strings.HasSuffix(f.Name(), Plugin_flag) == true {
				ret = true
				break
			}
		}
	}

	return ret
}

func (self *MainScreen) IsEmulatorPackage(dirname string) bool {
	ret := false
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, f := range files {
		if f.IsDir() {
			//pass
		} else {
			if strings.HasSuffix(f.Name(), Emulator_flag) == true {
				ret = true
				break
			}
		}
	}

	return ret
}

func (self *MainScreen) IsWifiConnectedNow() bool {

	cli := fmt.Sprintf("ip -4 addr show %s | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}'", sysgo.WifiDev)
	out := System(cli)
	if len(out) < 6 {
		return false
	}
	return true
}

func (self *MainScreen) GetWirelessIP() string {

	cli := fmt.Sprintf("ip -4 addr show %s | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}'", sysgo.WifiDev)
	out := SystemTrim(cli)

	return out
}

func (self *MainScreen) RunEXE(cmdpath string) {
	self.DrawRun()
	self.SwapAndShow()

	time.BlockDelay(1000)

	cmdpath = strings.Trim(cmdpath, " ")
	cmdpath = CmdClean(cmdpath)

	event.Post(RUNEVT, cmdpath)

}

func (self *MainScreen) HookExitCb( page PageInterface) {
	self.Child = append(self.Child,page)
}

func (self *MainScreen) OnExitCb() {
	PageLen := len(self.Child)

        for i := 0; i < PageLen; i++ {
		self.Child[i].OnExitCb()
	}
}

func (self *MainScreen) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == "T" {
		self.DrawRun()
		self.SwapAndShow()
		return
	}

	if ev.Data["Key"] == "Space" {
		self.Draw()
		self.SwapAndShow()
	}

	self.CurrentPage.KeyDown(ev)
	self.LastKey = ev.Data["Key"]

}
func (self *MainScreen) ShowMsg(content string, blocktime int) {
	self.MsgBox.SetText(content)
        self.MsgBox.Draw()
        self.SwapAndShow()

	if blocktime > 0 {
		time.BlockDelay(blocktime)
		self.CurrentPage.Draw()
		self.SwapAndShow()
	}

}
func (self *MainScreen) DrawRun() {
	self.MsgBox.SetText("Launching....")
	self.MsgBox.Draw()
}

func (self *MainScreen) Draw() {
	if self.CurrentPage != nil {
		self.CurrentPage.Draw()
	}

	if self.TitleBar != nil {
		//every plugin_init should not do any Draw actions since CurrentPage might be nil at that time
		self.TitleBar.Draw(self.CurrentPage.GetName())
	}

	if self.FootBar != nil {
		self.FootBar.SetLabelTexts(self.CurrentPage.GetFootMsg())
		self.FootBar.Draw()
	}
}

func (self *MainScreen) Refresh() {
	self.updateScreen <- true	
}

func (self *MainScreen) RefreshLoop() {
L:
	for {
		select {
	                case v:= <- self.updateScreen:
			if v == true {
				self.Draw()
				self.SwapAndShow()
			}
                        if v== false {
                                break L
                        }
		}
	}
}
