package UI

import (

	"strings"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/display"	
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/time"
	"github.com/cuu/gogame/event"
)

var (
	emulator_flag = "action.config"
	plugin_flag   = "plugin.config"
)

type MessageBox struct {
	Label
	Parent *MainScreen

}

func NewMessageBox() *MessageBox {
	m := &MessageBox{}
	m.Color = &color.Color{83,83,83,255}
	
	return m
}

func (self *MessageBox) Init( text string, font_obj *ttf.Font, col *color.Color) {
	if col != nil {
		self.Color = col
	}

	self.Text = text
	self.FontObj = font_obj

	self.Width = 0
	self.Height = 0

	self.CanvasHWND = surface.Surface(self.Parent.Width, self.Parent.Height)
	self.HWND       = self.Parent.CanvasHWND
	
}

func (self *MessageBox) SetText( text string) {
	self.Text = text
}

func (self *MessageBox) Draw() {
	self.Width = 0
	self.Height = 0
	surface.Fill(self.CanvasHWND, color.Color{255,255,255,255} )

	words := strings.Split(self.Text," ")
	space,_ := font.Size(self.FontObj," ")

	max_width := self.Parent.Width - 40
	x := 0
	y := 0

	row_total_width := 0
	lines := 0

	for _,word := range words {
		word_surface := font.Render( self.FontObj, word, true, self.Color,nil)
		word_width := word_surface.W
		word_height := word_surface.H
		row_total_width += word_width
		if lines == 0 {
			lines += word_height
		}

		if (row_total_width + space ) >= max_width {
			x = 0
			y += word_height
			row_total_width = word_width
			lines+=word_height
		}

		dest_rect := rect.Rect(x,y, word_width,word_height)
		surface.Blit(self.CanvasHWND, word_surface, &dest_rect,nil)
		if len(words) == 1 {
			x+=word_width
		} else {
			x += word_width+space
		}
		
		if x > self.Width {
			self.Width = x
		}

		if lines >= self.Parent.Height - 40 {
			break
		}
	}

	self.Height = lines

	padding := 5
	x = (self.Parent.Width - self.Width) / 2
	y = (self.Parent.Height - self.Height) /2

	rect_ := rect.Rect(x-padding,y-padding, self.Width+padding*2, self.Height+padding*2)
	
	draw.Rect(self.HWND , &color.Color{255,255,255,255},&rect_,0)

	if self.HWND != nil {
		rect__ := draw.MidRect(self.Parent.Width/2, self.Parent.Height/2,self.Width,self.Height,Width,Height)
		dest_rect := rect.Rect(0,0,self.Width,self,Height)
		surface.Blit(self.HWND, rect__, &dest_rect,nil)
	}

	draw.Rect(self.HWND , &color.Color{0,0,0,255},&rect_,1)
	
}

type MainScreen struct {
	Pages []PageInterface
	PageMax int
	PageIndex int
	PosX  int
	PosY  int
	Width int
	Height int
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
}


func NewMainScreen() *MainScreen {
	m := &MainScreen{}
  
	m.PosY = TitleBar_BarHeight+1
	m.Width = Width
	m.Height = Height - FootBar_BarHeight - TitleBar_BarHeight - 1
	m.MyPageStack = NewPageStack()

	m.MsgBoxFont = Fonts["veramono20"]
	m.IconFont   = Fonts["varela15"]
}

func (self *MainScreen) Init() {
	self.CanvasHWND = surface.Surface(self.Width,self.Height)
	
	self.MsgBox     = NewMessageBox()
	self.MsgBox.Parent = self
	self.MsgBox.Init(" ", self.MsgBoxFont, &color.Color{83,83,83})

	self.SkinManager = NewSkinManager()
	self.SkinManager.Init()

	
}

func (self *MainScreen) FartherPages() { // right after ReadTheDirIntoPages
	self.PageMax = len(self.Pages)

	for i:=0;i< self.PageMax; i++ {
		self.Pages[i].SetIndex(i)
		self.Pages[i].SetCanvasHWND(self.CanvasHWND)
		self.Pages[i].UpdateIconNumbers() // IconNumbers always == len(Pages[i].Icons)
		self.Pages[i].SetScreen(self)

		if self.Pages[i].GetIconNumbers() > 1 {
			self.Pages[i].SetPsIndex(1)
			self.Pages[i].SetIconIndex( 1 )
		}
	}

	self.CurrentPage = self.Pages[ self.PageIndex ]
	self.CurrentPage.SetOnShow(true)
}



func (self *MainScreen) CurPage() PageInterface {
	return self.CurrentPage
}

func (self *MainScreen) PushCurPage() {
	self.MyPageStack.Push(self.CurrentPage)
}

func (self *MainScreen) SetCurPage( pg PageInterface) {
	self.CurrentPage = pg
	pg.OnLoadCb()
}

func (self *MainScreen) PushPage( pg PageInterface) {
	self.PushCurPage()
	self.SetCurPage(pg)
}

func (self *MainScreen) AppendPage( pg PageInterface ) {
	self.Pages = append(self.Pages,pg)
}

func (self *MainScreen) ClearCanvas() {
	surface.Fill(self.CanvasHWND, color.Color{255,255,255,255} ) 
}

func (self *MainScreen) SwapAndShow() {
	if self.HWND != nil {
		rect_ := rect.Rect( self.PosX,self.PosY,self.Width,self.Height)
		surface.Blit(self.HWND,self.CanvasHWND,*rect_, nil)
	}
}

func (self *MainScreen) ExtraName(name string) string {

	parts := strings.Split(name,"_")
	if len(parts) >  1 {
		return parts[1]
	}else if len(parts) == 1 {
		return parts[0]
	}else {
		return name
	}	
}

func (self *MainScreen) IsPluginPackage(dirname string ) bool {
	ret := false
	files,err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
		return false
	}

	for _,f := range files {
		if f.IsDir() {
			//pass
		}else {
			if strings.HasSuffix(f.Name(),plugin_flag) == true {
				ret = true
				break
			}
		}
	}

	return ret
}

func (self *MainScreen) IsEmulatorPackage(dirname string ) bool {
	ret := false
	files,err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
		return false
	}

	for _,f := range files {
		if f.IsDir() {
			//pass
		}else {
			if strings.HasSuffix(f.Name(),emulator_flag) == true {
				ret = true
				break
			}
		}
	}

	return ret	
}

func (self *MainScreen) ReadTheDirIntoPages(_dir string, pglevel int, cur_page PageInterface) {
	
	if FileExists(_dir) == false && IsDirectory(_dir) == false {
		return
	}

	files,err := ioutil.ReadDir(_dir)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _,f := range files { // already sorted
		if IsDirectory( _dir +"/"+f.Name()) {
			if pglevel == 0 {
				page := NewPage()
				page.Name = self.ExtraName(f.Name())
				self.Pages = append(self.Pages, page)
				self.ReadTheDirIntoPages(_dir+"/"+f.Name(),pglevel+1, self.Pages[ len(self.Pages) - 1] )
			}else{ // on cur_page now
				i2:= self.ExtraName(f.Name())
				iconitem := NewIconItem()
				iconitem.AddLabel(i2,self.IconFont)
				if FileExists( SkinMap(_dir+"/"+i2+".png")) {
					iconitem.ImageName = SkinMap(_dir+"/"+i2+".png")
				}else {
					untitled := NewUntitledIcon()
					untitled.Init()
					if len(i2) > 1 {
						untitled.SetWords(i2[0],i2[1])
					}else if len(i2) == 1 {
						untitled.SetWords(i2[0],i2[0])
					}else {
						untitled.SetWords("G","s")
					}
					iconitem.ImgSurf = untitled.Surface()
					iconitem.ImageName = ""
				}

				if self.IsPluginPackage(_dir+"/"+f.Name()) {
					iconitem.MyType = ICON_TYPES["FUNC"]
					iconitem.CmdPath = f.Name()
					cur_page.AppendIcon(iconitem)
					//Init it 
				}else {
					iconitem.MyType = ICON_TYPES["DIR"]
					linkpage := NewPage()
					linkpage.Name = i2					
					iconitem.LinkPage = linkpage
					cur_page.AppendIcon(iconitem)
					self.ReadTheDirIntoPages(_dir+"/"+f.Name(),pglevel+1, iconitem.LinkPage)
				}
				
			}
		} else if IsAFile(_dir+"/"+f.Name()) && (pglevel > 0) {
			if strings.HasSuffix(strings.ToLower(f.Name()),IconExt) {
				i2 := self.ExtraName(f.Name())
				iconitem = NewIconItem()
				iconitem.CmdPath = _dir+"/"+f.Name()
				MakeExecutable( iconitem.CmdPath )
				iconitem.MyType = ICON_TYPES["EXE"]
				if FileExists( SkinMap( _dir+"/"+ ReplaceSuffix(i2,"png"))) {
					iconitem.ImageName = SkinMap( _dir+"/"+ ReplaceSuffix(i2,"png"))
				}else {
					untitled:= NewUntitledIcon()
					untitled.Init()
					if len(i2) > 1 {
						untitled.SetWords(i2[0],i2[1])
					}else if len(i2) == 1 {
						untitled.SetWords(i2[0],i2[0])
					}else {
						untitled.SetWords("G","s")
					}
					iconitem.ImgSurf = untitled.Surface()
					iconitem.ImageName = ""
				}

				iconitem.AddLabel(strings.Split(i2,".")[0], self.IconFont)
				iconfont.LinkPage = nil
				cur_page.AppendIcon(iconitem)
			}
		}
	}
}


func (self *MainScreen) RunEXE( cmdpath string) {
	self.DrawRun()
	self.SwapAndShow()

	time.Delay(1000)

	cmdpath = strings.Trim(cmdpath," ")

	cmdpath = CmdClean(cmdpath)
	
	event.Post(event.RUNEVT,cmdpath)
	
}

func (self *MainScreen) OnExitCb() {
	self.CurrentPage.OnExitCb()
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
}


func (self *MainScreen) DrawRun() {
	self.MsgBox.SetText("Launching....")
	self.MsgBox.Draw()
}

func (self *MainScreen) Draw() {
	self.CurrentPage.Draw()
	if self.TitleBar != nil {
		self.TitleBar.Draw( self.CurrentPage.GetName())
	}

	if self.FootBar != nil {
		self.FootBar.SetLabelTexts( self.CurrentPage.GetFootMsg())
		self.FootBar.Draw()
	}
}
