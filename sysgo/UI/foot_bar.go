package UI

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	
)

var FootBar_BarHeight = 20

type FootBarIconItem struct {
	MultiIconItem
	Parent *FootBar
}

func NewFootBarIconItem() *FootBarIconItem {
	m := &FootBarIconItem{}
	return m
}

func (self *FootBarIconItem) TotalWidth() int {
	lab_w,_ := self.Label.Size()
	return self.Width + lab_w
}

func (self *FootBarIconItem) Draw() {
	
	if self.Label != nil {
		lab_w,lab_h:= self.Label.Size()
		if self.Align == ALIGN["VCenter"] {
			self.Label.NewCoord( self.PosX - lab_w/2, self.PosY+ self.Height/2+12)
		} else if self.Align == ALIGN["HLeft"] {
			self.Label.NewCoord( self.PosX + self.Width/2+3, self.PosY - lab_h/2)
		}
		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		portion := rect.Rect(0, self.IconIndex*self.IconHeight, self.IconWidth, self.IconHeight)
		surface.Blit(self.Parent.CanvasHWND, self.ImgSurf, draw.MidRect(self.PosX,self.PosY, self.Width,self.Height, Width,Height),&portion)
		
	}
	
}

type FootBar struct {

	PosX int
	PosY int
	Width int
	Height int
	BarHeight int
	BorderWidth int
	CanvasHWND *sdl.Surface
	HWND       *sdl.Surface
	Icons   map[string]IconItemInterface
	IconWidth int
	IconHeight int
	LabelFont *ttf.Font
	State   string
	SkinManager *SkinManager
	icon_base_path string
	
}


func NewFootBar() *FootBar {
	f := &FootBar{}
	f.Width = Width

	f.BorderWidth = 1
	f.BarHeight = FootBar_BarHeight
	f.Height = 20

	f.IconWidth = 18
	f.IconHeight = 18
	
	f.LabelFont = Fonts["veramono10"]
	f.State = "normal"
	f.icon_base_path = SkinMap("gameshell/footbar_icons/")

	f.Icons = make(map[string]IconItemInterface)
	
	return f
}

func (self *FootBar) ReadFootBarIcons( icondir string) {
	if FileExists(icondir) == false && IsDirectory(icondir) == false {
		return
	}
	keynames := [5]string{"nav","x","y","a","b"}

	share_surf := image.Load(icon_base_path+"footbar.png")

	files,err := ioutil.ReadDir(icondir)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i,v := range keynames { // share_surf contains same number of image pieces of keynames
		it := NewFootBarIconItem()
		it.MyType = ICON_TYPES["NAV"]
		it.Parent = self
		it.ImgSurf = share_surf
		it.Align = ALIGN["HLeft"] // (X)Text
		it.AddLabel("game", self.LabelFont)
		it.Adjust( self.IconWidth/2+i*self.IconWidth, self.IconHeight/2+2, self.IconWidth,self.IconHeight,0)
		it.IconIndex = i
		self.Icons[v] = it
	}
}


func (self *FootBar) Init(main_screen *MainScreen) {
	self.CanvasHWND = surface.Surface(self.Width,self.Height)
	self.HWND = main_screen.HWND
	self.SkinManager = main_screen.SkinManager
	self.DBusManager = main_screen.DBusManager


	round_corners := NewFootBarIconItem()
	round_corners.IconWidth = 10
	round_corners.IconHeight = 10
	
	round_corners.MyType = ICON_TYPES["STAT"]
	round_corners.Parent = self
	round_corners.ImgSurf = MyIconPool.GetImageSurf["roundcorners"]
	round_corners.Adjust(0,0,10,10,0)
	
	self.Icons["round_corners"] = round_corners
	
}

func (self *FootBar) ResetNavText() {
	self.Icons["nav"].Label.SetText("Nav.")
	self.State = "normal"
	self.Draw()
}

func (self *FootBar) UpdateNavText(texts string) {
	self.State = "tips"
	my_text := font.Render(self.LabelFont, texts, true,self.SkinManager.GiveColor("Text"))

	left_width := self.Width - 18

	final_piece := ""

	for i,_ := range texts {
		text_ := texts[:i+1]
		my_text := font.Render(self.LabelFont, text_, true, self.SkinManager.GiveColor("Text"))
		final_piece  = text_
		if my_text.W >= left_width {
			break
		}
	}
	
	fmt.Printf("finalpiece %s\n", final_piece)

	self.Icons["nav"].Label.SetText(final_piece)
	self.Draw()
	
}

func (self *FootBar) SetLabelTexts( texts []string) {
	keynames := [5]string{"nav","x","y","a","b"}
	if len(texts) < 5 {
		log.Fatal("SetLabelTexts texts length error")
		return
	}

	for idx,x := range keynames {
		self.Icons[x].Label.SetText(texts[idx])
	}
	
}

func (self *FootBar) ClearCanvas() {
	surface.Fill( self.CanvasHWND,  self.SkinManager.GiveColor("White"))

	self.Icons["round_corners"].NewCoord(5,self.Height-5)
	self.Icons["round_corners"].SetIconIndex(2)
	self.Icons["round_corners"].Draw()


	self.Icons["round_corners"].NewCoord(self.Width - 5,self.Height - 5)
	self.Icons["round_corners"].SetIconIndex(3)

	self.Icons["round_corners"].Draw()
	
}

func (self *FootBar) Draw() {
	self.ClearCanvas()
	self.Icons["nav"].NewCoord(self.IconWidth/2+3, self.IconHeight/2+2)
	self.Icons["nav"].Draw()

	if self.State == "normal" {
		_w := 0

		for i,x := range []string{"b","a","y","x"} {
			if self.Icons[x].Label.GetText() != "" {
				if i== 0 {
					_w += self.Icons[x].TotalWidth()
				}else {
					_w += self.Icons[x].TotalWidth()+5
				}

				start_x := self.Width - _w
				start_y := self.IconHeight/2+2
				self.Icons[x].NewCoord(start_x, start_y)
				self.Icons[x].Draw()
			}
		}

		draw.Line(self.CanvasHWND, self.SkinManager.GiveColor("Line"),0,0,Width,0,self.BorderWidth)

		if self.HWND != nil {
			rect_ := rect.Rect(self.PosX, Height - self.Height, Width, self.BarHeight)
			surface.Blit(self.HWND,self.CanvasHWND, &rect_,nil)
		}
}
