package UI

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	
	"github.com/cuu/gogame/color"
)

type Textarea struct {
	Widget
	BackgroundColor *color.Color
	CanvasHWND *sdl.Surface
	MyWords []string
	BlitWords []string
	FontObj *ttf.Font
	LineNumber int
	TextLimit int
	TextFull bool
	TextIndex int
	BlitIndex int
}

func NewTextarea() *Textarea {
	p := &Textarea{}

	p.TextLimit = 63
	p.TextFull = false

	p.MyWords = make([]string,0)
	p.BlitWords = make([]string,0)

	p.BackgroundColor = &color.Color{228,228,228,255}
	
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
			self.MyWords = append(self.MyWords[:self.TextIndex],self.MyWords[(self.TextIndex+1):]...)
		}
	}

	return self.MyWords
}


func (self *Textarea) AppendText( alphabet string) {
	self.AppendAndBlitText(alphabet)
}

func (self *Textarea) AppendAndBlitText(alphabet string) {
	if self.TextFull == false {
		
		if self.TextIndex <= len(self.MyWords) {
			m = append(m[:idx], append([]string{"U"}, m[idx:]...)...)
			self.MyWords = append(self.MyWords[:self.TextIndex],
				append([]string{alphabet},self.MyWords[self.TextIndex:]...)...)
			
			self.BlitText()
			self.AddTextIndex()
		}
		
	}else {
		fmt.Printf("is Full %s",strings.Join(self.MyWords,""))
	}
	
}

func (self *Textarea) BuildBlitText() {
	blit_rows := make([][]string,0)

	w := 0
	xmargin := 5
	endmargin :=15
	x := self.PosX + xmargin
	linenumber := 0
	cursor_row := 0

	
}
