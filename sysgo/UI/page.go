package UI

import (
	"sync"
	
	"github.com/veandco/go-sdl2/sdl"
	
	"github.com/cuu/gogame/font"
	
)

type element struct {
    data interface{}
    next *element
}

type PageStack struct {
    lock *sync.Mutex
    head *element
    Size int
}

func (stk *PageStack) Push(data interface{}) {
    stk.lock.Lock()

    element := new(element)
    element.data = data
    temp := stk.head
    element.next = temp
    stk.head = element
    stk.Size++

    stk.lock.Unlock()
}

func (stk *PageStack) Pop() interface{} {
    if stk.head == nil {
        return nil
    }
    stk.lock.Lock()
    r := stk.head.data
    stk.head = stk.head.next
    stk.Size--

    stk.lock.Unlock()

    return r
}

func (stk *PageStack) Length() int {
	return stk.Size
}

func NewPageStack() *PageStack {
    stk := new(PageStack)
    stk.lock = &sync.Mutex{}
    return stk
}


type PageSelectorInterface interface {
	Adjust(x,y,w,h,alpha int)
	Draw()
}

type PageSelector struct {
	
	PosX int
	PosY int
	Width int
	Height int
	Parent PageInterface
	Alpha int
	OnShow bool
	IconSurf  *sdl.Surface
	
}

func NewPageSelector() *PageSelector {
	p := &PageSelector{}
	return p
}

func (self *PageSelector) Init(x,y,w,h,alpha int) {
	self.Adjust(x,y,w,h,alpha)
}

func (self *PageSelector) Adjust(x,y,w,h,alpha int) {
	self.PosX = x
	self.PosY = y
	self.Width = w
	self.Height = h
	self.Alpha  = alpha
}

func (self *PageSelector) Draw() {
	canvas  := self.Parent.GetCanvasHWND()
	idx     := self.Parent.GetPsIndex()
	iconidx := self.Parent.GetIconIndex()
	icons   := self.Parent.GetIcons()
	
	if idx < len(icons) {
		icon_x ,_ := icons[idx].Coord()
		_,icon_y  := icons[iconidx].Coord()
		
		parent_x,parent_y := self.Parent.Coord()
		parent_w,parent_h := self.Parent.Size()
		
		x := icon_x + parent_x
		y := icon_y // only use current icon's PosY
		
		rect_ = draw.MidRect(x,y, self.Width, self.Height, parent_w,parent_h)
		if rect_.W <=0 || rect_.H <= 0 {
			return
		}
		
		if self.IconSurf != nil {
			surface.Blit(canvas,self.IconSurf, rect_,nil)
		}
		
	}
}


type PageInterface interface {
	// shared functions
	// GetScreen
	// GetIcons
	// SetScreen
	// SetFootMsg
	// SetCanvasHWND
	// GetCanvasHWND
	// GetHWND
	// SetHWND
	// AdjustHLeftAlign
	// AdjustSAutoLeftAlign
	// SetPsIndex
	// SetIconIndex
	// GetPsIndex
	// GetIconIndex
	// Coord
	// Size

	
}

type Page struct {
	PosX int
	PosY int
	Width int
	Height int
	Icons []IconItemInterface // slice ,use append
	IconNumbers int
	IconIndex int
	PrevIconIndex int
	
	Ps PageSelectorInterface
	PsIndex int

	Index int

	Align string
	
	CanvasHWND *sdl.Surface
	HWND       *sdl.Surface

	OnShow bool
	Name  string
	Screen *MainScreen
	
	PageIconMargin int // default 20
	FootMsg  [5]string

	SelectedIconTopOffset int
	EasingDur int
}

func NewPage() *Page {
	p := &Page{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 30

	p.Align = ALIGN["SLeft"]
	
	p.FootMsg = [5]string{"Nav.","","","","Enter"}
	
	return p
}

func (self *Page) AdjustHLeftAlign() {
	self.PosX = self.Index*self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	cols := int(Width/IconWidth)
	rows := int( self.IconNumbers * IconWidth) / self.Width + 1
	cnt := 0
	
	if rows < 1 {
		rows = 1
	}
	
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			start_x := IconWidth/2 + j*IconWidth
			start_y := IconHeight/2 + i*IconHeight
			icon    := self.Icons[cnt]
			icon.Adjust(start_x,start_y,IconWidth-4,IconHeight-4,0)
			icon.SetIndex(cnt)
			icon.SetParent(self)
			if cnt >= self.IconNumbers -1 {
				break
			}
			cnt += 1
		}
	}

	ps := NewPageSelector()
	ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
	ps.Parent = self

	ps.Init(IconWidth/2,TitleBar_BarHeight+IconHeight/2, 92,92,128) //hard coded of the blueselector png size
	
	self.Ps = ps
	self.PsIndex = 0
	self.OnShow = false
}

func (self *Page) AdjustSLeftAlign() { // ## adjust coordinator and append the PageSelector
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	start_x := (self.PageIconMargin + IconWidth + self.PageIconMargin ) / 2
	start_y := self.Height/2

	for i := 0; i < self.IconNumbers; i++ {
		it := self.Icons[i]
		it.SetParent(self)
		it.SetIndex(i)
		it.Adjust(start_x+i*self.PageIconMargin+i*IconWidth, start_y, IconWidth-6,IconHeight-6,0)

		old_surf := it.GetImageSurf()
		
		it_w,it_h := it.Size() //width height changed by Adjust above
		it.SetImageSurf( transform.SmoothScale(old_surf,it_w,it_h) )
	}

	ps := NewPageSelector()
	ps.IconSurf = MyIconPool.GetImageSurf("blueselector")
	ps.Parent = self
	ps.Init(start_x,start_y,92,92,128)

	self.Ps = ps
	self.PsIndex = 0
	self.OnShow = false

	if self.IconNumbers > 1 {
		self.PsIndex = 1
		self.IconIndex = self.PsIndex
		self.PrevIconIndex = self.IconIndex
		cur_icon_x,cur_icon_y := self.Icons[self.IconIndex].Coord()
		self.Icons[self.IconIndex].NewCoord(cur_icon_x, cur_icon_y - self.SelectedIconTopOffset )
	}
}


func (self *Page) AdjustSAutoLeftAlign() { //  ## adjust coordinator and append the PageSelector
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	start_x := (self.PageIconMargin + IconWidth + self.PageIconMargin ) / 2
	start_y := self.Height/2

	
}


