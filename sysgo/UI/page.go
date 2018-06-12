package UI

import (
	"sync"
	
	"github.com/veandco/go-sdl2/sdl"
	
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

func NewPageStack() *PageStack {
    stk := new(PageStack)
    stk.lock = &sync.Mutex{}
    return stk
}


type PageSelector struct {
	
	PosX int
	PosY int
	Width int
	Height int
	Parent interface{} //
	Alpha int
	OnShow bool
	IconSurf  *sdl.Surface
	
}

func (p *PageSelector) Adjust(x,y,w,h,alpha int) {
	p.PosX = x
	p.PosY = y
	p.Width = w
	p.Height = h
	p.Alpha  = alpha
}

func (p *PageSelector) Draw() {
	
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
	Icons []interface{} // make first
	IconNumbers int
	IconIndex int
	PrevIconIndex int
	
	Ps interface{}
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
	
	p.FootMsg = [5]string{"Nav.","","","","Enter"}

	return p
}






