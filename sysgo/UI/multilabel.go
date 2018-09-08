package UI

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/font"
)

//MultiLabel is also a LabelInterface
type MultiLabel struct {
	PosX       int
	PosY       int
	Width      int
	Height     int
	Text       string
	FontObj    *ttf.Font
	Color      *color.Color
	CanvasHWND *sdl.Surface
	//TextSurf *sdl.Surface
	MaxWidth   int
}

func NewMultiLabel() *MultiLabel{
	l := &MultiLabel{}
	l.Color = &color.Color{83,83,83,255}
	l.Width = 135
	l.Height = 100

	return l
}

func (self *MultiLabel) Init(text string, font_obj *ttf.Font,col *color.Color) {
	if col != nil {
		self.Color = col
	}
	
	self.Text = text
	self.FontObj = font_obj

	if self.CanvasHWND != nil {
		self.Draw()
	}
}

func (self *MultiLabel) SetCanvasHWND( canvas *sdl.Surface) {
	self.CanvasHWND = canvas
}

func (self *MultiLabel) Coord() (int,int) {
	return self.PosX,self.PosY
}

func (self *MultiLabel) Size() (int,int) {
	return self.Width,self.Height
}

func (self *MultiLabel) NewCoord(x,y int) {
	self.PosX = x
	self.PosY = y
	
}

func (self *MultiLabel) SetColor(col *color.Color){
	if col != nil {
		self.Color = col
	}
}

func (self *MultiLabel) GetText() string {
	return self.Text
}


func (self *MultiLabel) SetText(text string) {
	self.Text = text
	
}

func (self *MultiLabel) Draw() {
	font.SetBold(self.FontObj,false) // avoing same font tangling set_bold to others
	self.blit_text(self.CanvasHWND, self.Text,self.PosX,self.PosY,self.FontObj)	
}

// difference to Label
func (self *MultiLabel) blit_text(surf *sdl.Surface,text string, pos_x,pos_y int, fnt *ttf.Font) {
	words := make([][]string, 0)
	temp := strings.Split(text,"\n")
	for _,v := range temp {
		t := strings.Split(v," ")
		words = append(words,t)
	}

	space,_   := font.Size(fnt," ")
	max_width := self.Width
	x,y       := pos_x,pos_y

	row_total_width := 0
	lines := 0

	for i,line := range words[:4] {
		for _,word := range line[:12] {
			word_surface := font.Render(fnt,word,true,self.Color,nil)
			word_width   := surface.GetWidth(word_surface)
			word_height  := surface.GetHeight(word_surface)
			row_total_width += word_width
			if row_total_width+space >= max_width {
				x = pos_x
				y = y+word_height
				row_total_width = word_width
				if lines == 0 {
					lines = lines + word_height
				}else {
					lines = lines + word_height
				}	
			}
			rect_ := rect.Rect(x,y,self.Width,self.Height)
			surface.Blit(surf,word_surface,&rect_,nil)
			x += (word_width+space)	
		}
		x = pos_x
		y += word_height
		lines += word_height
	}

	self.Height = lines
	
}

