package UI

import (
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/draw"
)

type ListItemInterface interface {
	Init(text string)

	Size() (int, int)
	NewSize(w, h int)
	Coord() (int, int)
	NewCoord(x, y int)

	GetLinkObj() PluginInterface
	Draw()
}

type ListItem struct {
	Widget
	Labels map[string]LabelInterface
	Icons  map[string]IconItemInterface
	Fonts  map[string]*ttf.Font

	LinkObj PluginInterface

	Parent PageInterface
}

func NewListItem() *ListItem {
	i := &ListItem{}
	i.Labels = make(map[string]LabelInterface)
	i.Icons = make(map[string]IconItemInterface)
	i.Fonts = make(map[string]*ttf.Font)

	i.Height = 30
	i.Width = 0

	return i
}

func (self *ListItem) Init(text string) {
	l := NewLabel()
	l.PosX = 16
	l.SetCanvasHWND(self.Parent.GetCanvasHWND())
	l.Init(text, self.Fonts["normal"], nil)
	self.Labels["Text"] = l

}

func (self *ListItem) Coord() (int, int) {
	return self.PosX, self.PosY
}

func (self *ListItem) Size() (int, int) {
	return self.Width, self.Height
}

func (self *ListItem) GetLinkObj() PluginInterface {
	return self.LinkObj
}

func (self *ListItem) Draw() {
	x_, _ := self.Labels["Text"].Coord()
	_, h_ := self.Labels["Text"].Size()

	self.Labels["Text"].NewCoord(x_, self.PosY+(self.Height-h_)/2)
	self.Labels["Text"].Draw()

	draw.Line(self.Parent.GetCanvasHWND(), &color.Color{169, 169, 169, 255},
		self.PosX, (self.PosY + self.Height - 1),
		(self.PosX + self.Width), (self.PosY + self.Height - 1), 1)

}
