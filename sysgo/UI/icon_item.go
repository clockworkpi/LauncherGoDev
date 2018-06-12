package UI

import (
	
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/transform"
	"github.com/cuu/gogame/utils"
)

type IconItemInterface interface {
	Init(x,y,w,h,at int)
	SetIndex(i int)
	SetParent( p interface{} )
	SetLabelColor(col *color.Color)
	Coord() (int,int)
	NewCoord(x,y int)
	Size() (int,int)
	AddLabel(text string, fontobj *ttf.Font)
	AdjustLinkPage()
	GetImageSurf() *sdl.Surface
	SetImageSurf(newsurf *sdl.Surface)
	CreateImageSurf()
	ChangeImgSurfColor(col *color.Color)
	Draw() 
}

type IconItem struct {
	PosX int
	PosY int
	Width int
	Height int
	ImageName string
	ImgSurf *sdl.Surface
	Parent  PageInterface
	Index   int
	MyType int
	CmdPath  interface{}
	LinkPage PageInterface
	Label  LabelInterface
	Align  int
	AnimationTime int
}


func NewIconItem() *IconItem {
	i := &IconItem{}
	i.MyType = ICON_TYPES["EXE"]

	i.Align = ALIGN["VCenter"]
	
	l := NewLabel()
	
	i.Label = l
	
	return i
}


func (self *IconItem) Init(x,y,w,h,at int) {
	self.PosX = x
	self.PosY = y
	self.Width = w
	self.Height = h
	self.AnimationTime = at
}

func (self *IconItem) SetIndex(i int) {
	self.Index = i
}

func (self *IconItem) SetParent(p interface{} ) {
	self.Parent = p
}

func (self *IconItem) SetLabelColor(col *color.Color) {
	self.Label.SetColor(col)
}

func (self *IconItem) Coord() (int,int) {
	return self.PosX,self.PosY
}

func (self *IconItem) NewCoord(x,y int) {
	self.PosX = x
	self.PosY = y
}

func (self *IconItem) NewCoord(x,y int) {
	return self.Width,self.Height
}

func (self *IconItem) AddLabel(text string, fontobj *ttf.Font) {
	if self.Label == nil {
		l:= NewLabel()
		self.Label = l
	}else {
		self.Label.Init(text,fontobj)
	}	
}

func (self *IconItem) AdjustLinkPage() {
	if self.MyType == ICON_TYPES["DIR"] && self.LinkPage != nil {
		self.LinkPage.SetIndex(0)
		self.LinkPage.SetAlign(ALIGN["SLeft"])
		self.LinkPage.SetIconNumbers( len(self.LinkPage.GetIcons()) )
		self.LinkPage.SetScreen(self.Parent.GetScreen())
		self.LinkPage.SetCanvasHWND( (self.Parent.GetScreen()).CanvasHWND )
		self.LinkPage.SetFootMsg([5]string{ "Nav.","","","Back","Enter" } )
		if self.LinkPage.GetAlign() == ALIGN["HLeft"] {
			self.LinkPage.AdjustHLeftAlign()
		}else if self.LinkPage.GetAlign() == ALIGN["SLeft"] {
			self.LinkPage.AdjustSAutoLeftAlign()
			if self.LinkPage.GetIconNumbers() > 1 {
				self.LinkPage.SetPsIndex(1)
				self.LinkPage.SetIconIndex ( 1 ) 
			}
		}
	}
}


func (self *IconItem) GetImageSurf() *sdl.Surface {
	return self.ImgSurf
}

func (self *IconItem) SetImageSurf(newsurf *sdl.Surface) {
	self.ImgSurf = newsurf
}


func (self *IconItem) CreateImageSurf() {
	if self.ImgSurf == nil && self.ImageName != "" {
		self.ImgSurf = image.Load(self.ImageName)
		if self.ImgSurf.W  > IconWidth  || self.ImgSurf.H > IconHeight {
			self.ImgSurf = transform.Scale(self.ImgSurf,IconWidth,IconHeight)
		}
	}
}

func (self *IconItem) ChangeImgSurfColor(col *color.Color) {
	utils.ColorSurface(self.ImgSurf,col)
}

func (self *IconItem) Draw() {
	
	parent_x,parent_y := self.Parent.Coord()
	
	if self.Label != nil {
		lab_x,lab_y := self.Label.Coord()
		lab_w,lab_h:= self.Label.Size()
		
		if self.Align == ALIGN["VCenter"] {
			self.Label.NewCoord( self.PosX - lab_w/2 + parent_x, self.PosY + lab_h/2+6+parent_y)
		}else if self.Align == ALIGN["HLeft"] {
			self.Label.NewCoord( self.PosX + self.Width/2+3+parent_x, self.PosY - lab_h/2 + parent_y)
		}

		self.Label.Draw()
	}

	if self.ImgSurf != nil {
		surface.Blit(self.Parent.GetCanvasHWND(), self.ImgSurf,draw.MidRect(self.PosX + parent_x, self.PosY + parent_y,
			self.Width,self.Height, Width, Height),nil)
	}
}

