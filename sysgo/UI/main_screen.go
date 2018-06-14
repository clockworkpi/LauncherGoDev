package UI

import (

	"strings"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/display"	
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/color"
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
	
}
