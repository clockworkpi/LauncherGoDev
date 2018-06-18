package UI

import (
	
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/image"
	"github.com/cuu/gogame/transform"
	"github.com/cuu/gogame/utils"
)

type IconItemInterface interface {
	Init(x,y,w,h,at int)
	Adjust(x,y,w,h,at int)
	GetCmdPath() string
	SetCmdPath( path string)
	
	SetMyType( thetype int )
	GetMyType() int

	GetIconIndex() int
	SetIconIndex(idx int)
	
	GetIndex() int
	SetIndex(i int)
	
	SetParent( p interface{} )
	
	SetLabelColor(col *color.Color)
	SetLabelText(text string)
	GetLabelText() string
	
	Coord() (int,int)
	NewCoord(x,y int)

	TotalWidth() int
	Size() (int,int)

	
	AddLabel(text string, fontobj *ttf.Font)
	GetLinkPage() PageInterface
	AdjustLinkPage()
	GetImgSurf() *sdl.Surface
	SetImgSurf(newsurf *sdl.Surface)
	CreateImgSurf()
	ChangeImgSurfColor(col *color.Color)
	
	Clear()

	GetCmdInvoke() PluginInterface

	
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
	IconIndex int
	MyType int
	CmdPath  string
	CmdInvoke PluginInterface
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

func (self *IconItem) Adjust(x,y,w,h,at int) {
	self.PosX = x
	self.PosY = y
	self.Width = w
	self.Height = h
	self.AnimationTime = at

	if self.Label != nil {
		self.Label.SetCanvasHWND(self.Parent.GetCanvasHWND())
	}

	self.CreateImgSurf()
	self.AdjustLinkPage()
	
}


func (self *IconItem) GetCmdPath() string {
	return self.CmdPath
}

func (self *IconItem) SetCmdPath( path string) {
	self.CmdPath = path
}

func (self *IconItem) SetMyType( thetype int ) {
	self.MyType = thetype
}

func (self *IconItem) GetMyType() int {
	return self.MyType
}

func (self *IconItem) GetIconIndex() int {
	return self.IconIndex
}

func (self *IconItem) SetIconIndex( idx int) {
	self.IconIndex = idx
}

func (self *IconItem) GetIndex() int {
	return self.Index
}

func (self *IconItem) SetIndex(i int) {
	self.Index = i
}

func (self *IconItem) SetParent(p interface{} ) {
	self.Parent = p.(PageInterface)
}

func (self *IconItem) SetLabelColor(col *color.Color) {
	self.Label.SetColor(col)
}

func (self *IconItem) GetLabelText() string {
	return self.Label.GetText()
}

func (self *IconItem) SetLabelText(text string) {
	self.Label.SetText(text)
}

func (self *IconItem) Coord() (int,int) {
	return self.PosX,self.PosY
}

func (self *IconItem) NewCoord(x,y int) {
	self.PosX = x
	self.PosY = y
}

func (self *IconItem) TotalWidth() int {
	return 0
}

func (self *IconItem) Size() (int,int) {
	return self.Width,self.Height
}

func (self *IconItem) AddLabel(text string, fontobj *ttf.Font) {
	if self.Label == nil {
		l:= NewLabel()
		self.Label = l
	}else {
		self.Label.Init(text,fontobj,nil)
	}	
}

func (self *IconItem) GetLinkPage() PageInterface {
	return self.LinkPage
}

func (self *IconItem) AdjustLinkPage() {
	if self.MyType == ICON_TYPES["DIR"] && self.LinkPage != nil {
		self.LinkPage.SetIndex(0)
		self.LinkPage.SetAlign(ALIGN["SLeft"])
		self.LinkPage.UpdateIconNumbers()
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


func (self *IconItem) GetImgSurf() *sdl.Surface {
	return self.ImgSurf
}

func (self *IconItem) SetImgSurf(newsurf *sdl.Surface) {
	self.ImgSurf = newsurf
}


func (self *IconItem) CreateImgSurf() {
	if self.ImgSurf == nil && self.ImageName != "" {
		self.ImgSurf = image.Load(self.ImageName)
		if int(self.ImgSurf.W)  > IconWidth  || int(self.ImgSurf.H) > IconHeight {
			self.ImgSurf = transform.Scale(self.ImgSurf,IconWidth,IconHeight)
		}
	}
}

func (self *IconItem) ChangeImgSurfColor(col *color.Color) {
	utils.ColorSurface(self.ImgSurf,col)
}

func (self *IconItem) Clear() {
	
}

func (self *IconItem) GetCmdInvoke() PluginInterface {
	return self.CmdInvoke
}

func (self *IconItem) Draw() {
	
	parent_x,parent_y := self.Parent.Coord()
	
	if self.Label != nil {
//		lab_x,lab_y := self.Label.Coord()
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

