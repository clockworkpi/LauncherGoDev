package UI

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
)

type Textarea struct {
	Widget
	BackgroundColor *color.Color
	CanvasHWND      *sdl.Surface
	MyWords         []string
	BlitWords       []string
	FontObj         *ttf.Font
	LineNumber      int
	TextLimit       int
	TextFull        bool
	TextIndex       int
	BlitIndex       int
}

func NewTextarea() *Textarea {
	p := &Textarea{}

	p.TextLimit = 63
	p.TextFull = false

	p.MyWords = make([]string, 0)
	p.BlitWords = make([]string, 0)

	p.BackgroundColor = &color.Color{228, 228, 228, 255}

	return p
}

func (self *Textarea) Init() {
	self.FontObj = Fonts["veramono24"]
}

func (self *Textarea) SubTextIndex() {
	self.TextIndex -= 1
	if self.TextIndex < 0 {
		self.TextIndex = 0
	}
}

func (self *Textarea) AddTextIndex() {

	self.TextIndex += 1
	if self.TextIndex > len(self.MyWords) {
		self.TextIndex = len(self.MyWords)
	}
}

func (self *Textarea) ResetMyWords() {
	self.MyWords = nil
	self.TextIndex = 0
}

func (self *Textarea) RemoveFromLastText() []string {
	if len(self.MyWords) > 0 {
		self.SubTextIndex()
		if self.TextIndex < len(self.MyWords) {
			self.MyWords = append(self.MyWords[:self.TextIndex], self.MyWords[(self.TextIndex+1):]...)
		}
	}

	return self.MyWords
}

func (self *Textarea) AppendText(alphabet string) {
	self.AppendAndBlitText(alphabet)
}

func (self *Textarea) AppendAndBlitText(alphabet string) {
	if self.TextFull == false {

		if self.TextIndex <= len(self.MyWords) {
			self.MyWords = append(self.MyWords[:self.TextIndex],
				append([]string{alphabet}, self.MyWords[self.TextIndex:]...)...)

			self.BlitText()
			self.AddTextIndex()
		}

	} else {
		fmt.Printf("is Full %s\n", strings.Join(self.MyWords, ""))
	}

}

func (self *Textarea) BuildBlitText() {
	blit_rows := make([][]string, 0)
	blit_rows = append(blit_rows, []string{})

	w := 0
	//	xmargin   := 5
	endmargin := 15
	linenumber := 0
	cursor_row := 0

	for i, v := range self.MyWords {
		t := font.Render(self.FontObj, v, true, &color.Color{8, 135, 174, 255}, nil)
		t_width := surface.GetWidth(t)
		w += t_width

		blit_rows[linenumber] = append(blit_rows[linenumber], v)

		if i == self.TextIndex-1 {
			cursor_row = linenumber
		}

		if w+t_width >= self.Width-endmargin {
			w = 0
			linenumber += 1
			blit_rows = append(blit_rows, []string{})
		}
		t.Free()
	}

	if len(blit_rows) == 1 {
		self.BlitWords = blit_rows[0]
		self.BlitIndex = self.TextIndex
	} else if len(blit_rows) == 2 || cursor_row < 2 {
		self.BlitWords = append(blit_rows[0], blit_rows[1]...)
		self.BlitIndex = self.TextIndex

	} else {
		self.BlitWords = append(blit_rows[cursor_row-1], blit_rows[cursor_row]...)
		self.BlitIndex = self.TextIndex

		for i, v := range blit_rows {
			if i == cursor_row-1 {
				break
			}
			self.BlitIndex -= len(v)
		}
	}

}

func (self *Textarea) BlitText() {
	//blit every single word into surface and calc the width ,check multi line
	self.BuildBlitText()

	w := 0
	xmargin := 5
	endmargin := 15

	x := self.PosX + xmargin
	y := self.PosY

	linenumber := 0

	if len(self.MyWords) > self.TextLimit {
		self.TextFull = true
	} else {
		self.TextFull = false
	}

	for _, v := range self.BlitWords {
		t := font.Render(self.FontObj, v, true, &color.Color{8, 135, 174, 255}, nil)
		w += surface.GetWidth(t)

		if w >= self.Width-endmargin && linenumber == 0 {
			linenumber += 1
			x = self.PosX + xmargin
			y = self.PosY + surface.GetHeight(t)*linenumber
			w = 0
		}

		rect_ := rect.Rect(x, y, 0, 0)
		surface.Blit(self.CanvasHWND, t, &rect_, nil)
		x += surface.GetWidth(t)
		t.Free()
	}
}

func (self *Textarea) Cursor() {
	w := 0
	xmargin := 5
	endmargin := 15
	x := self.PosX + xmargin
	y := self.PosY
	linenumber := 0

	for _, v := range self.BlitWords[:self.BlitIndex] {
		t := font.Render(self.FontObj, v, true, &color.Color{8, 135, 174, 255}, nil)
		w += surface.GetWidth(t)

		if w >= self.Width-endmargin && linenumber == 0 {
			x = self.PosX + xmargin
			y = self.PosY + surface.GetHeight(t)
			w = 0
			linenumber += 1
		}

		if w >= self.Width-endmargin*3 && linenumber > 0 {
			x += surface.GetWidth(t)
			break
		}
		x += surface.GetWidth(t)
		t.Free()
	}

	c_t := font.Render(self.FontObj, "_", true, &color.Color{0, 0, 0, 255}, nil)
	rect_ := rect.Rect(x+1, y-2, 0, 0)
	surface.Blit(self.CanvasHWND, c_t, &rect_, nil)
	c_t.Free()
}

func (self *Textarea) Draw() {

	rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)

	draw.AARoundRect(self.CanvasHWND, &rect_, self.BackgroundColor, 4, 0, self.BackgroundColor)

	self.BlitText()
	self.Cursor()

}
