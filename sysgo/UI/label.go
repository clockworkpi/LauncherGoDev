package UI

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/font"
)

type LabelInterface interface {
	Init( text string, font_obj *ttf.Font,col *color.Color )
	SetCanvasHWND( canvas *sdl.Surface)
	Coord() (int,int)
	Size() (int,int)
	NewCoord(x,y int)
	NewSize(w,h int)
	SetColor(col *color.Color )
	GetText() string
	SetText(text string)
	Draw()
}

type Label struct {
	Widget
	Text string
	FontObj *ttf.Font
	Color  *color.Color
	CanvasHWND *sdl.Surface
//	TextSurf *sdl.Surface
}

func NewLabel() *Label {
	l := &Label{}
	l.Color = &color.Color{83,83,83,255}
	return l
}

func (self *Label) Init(text string, font_obj *ttf.Font,col *color.Color ) {
	if col != nil {
		self.Color = col
	}
		

	self.Text = text

	self.FontObj = font_obj

	self.Width,self.Height = font.Size(self.FontObj, self.Text)
	
}

func (self *Label) SetCanvasHWND( canvas *sdl.Surface) {
	self.CanvasHWND = canvas
}

func (self *Label) Coord() (int,int) {
	return self.PosX,self.PosY
}

func (self *Label) Size() (int,int) {
	return self.Width,self.Height
}

func (self *Label) NewCoord(x,y int) {
	self.PosX = x
	self.PosY = y
	
}

func (self *Label) SetColor(col *color.Color){
	if col != nil {
		self.Color = col
	}
}

func (self *Label) GetText() string {
	return self.Text
}


func (self *Label) SetText(text string) {
	self.Text = text
	self.Width,self.Height = font.Size(self.FontObj, self.Text)
}

func (self *Label) Draw() {
	font.SetBold(self.FontObj,false) // avoing same font tangling set_bold to others
	my_text := font.Render(self.FontObj,self.Text, true, self.Color, nil)

	rect_ := rect.Rect(self.PosX,self.PosY,self.Width,self.Height)
	
	surface.Blit(self.CanvasHWND,my_text,&rect_,nil)
	
}
