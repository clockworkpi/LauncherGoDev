package UI

import (
	"fmt"
	
//	"math"
	"sync"
	
	"github.com/veandco/go-sdl2/sdl"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/draw"
//	"github.com/cuu/gogame/rect"
//	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/event"

	"github.com/cuu/gogame/transform"
	"github.com/cuu/LauncherGo/sysgo/easings"
	
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
	Init(x,y,w,h,alpha int)
	Adjust(x,y,w,h,alpha int)
	GetOnShow() bool
	SetOnShow(onshow bool)

	Coord() (int,int)
	NewCoord(x,y int)
	Size() (int,int)
	NewSize(w,h int)
	
	Draw()
}

type PageSelector struct {
	Widget
	
	Alpha int
	OnShow bool
	IconSurf  *sdl.Surface
	
	Parent PageInterface
}

func NewPageSelector() *PageSelector {
	p := &PageSelector{}
	p.OnShow = true
	
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

func (self *PageSelector) GetOnShow() bool {
	return self.OnShow
}

func (self *PageSelector) SetOnShow( onshow bool ) {
	self.OnShow = onshow
}

func (self *PageSelector) Draw() {
	canvas  := self.Parent.GetCanvasHWND()
	idx     := self.Parent.GetPsIndex()
	iconidx := self.Parent.GetIconIndex()
	icons   := self.Parent.GetIcons()
	
	if idx < len(icons) {
		icon_x ,_ := icons[idx].Coord()
		_,icon_y  := icons[iconidx].Coord()
		
		parent_x,_ := self.Parent.Coord()
		parent_w,parent_h := self.Parent.Size()
		
		x := icon_x + parent_x
		y := icon_y // only use current icon's PosY
		
		rect_ := draw.MidRect(x,y, self.Width, self.Height, parent_w,parent_h)
		if rect_.W <=0 || rect_.H <= 0 {
			return
		}
		
		if self.IconSurf != nil {
			surface.Blit(canvas,self.IconSurf, rect_,nil)
		}
		
	}
}


type PageInterface interface {
	// ## shared functions ##
	Adjust()
	Init()
	
	GetScreen() *MainScreen
	GetIcons() []IconItemInterface
	SetScreen( main_screen *MainScreen)
	SetFootMsg(footmsg [5]string)
	GetCanvasHWND() *sdl.Surface
	SetCanvasHWND( canvas *sdl.Surface)

	GetHWND() *sdl.Surface
	SetHWND(h *sdl.Surface)
	
	AdjustHLeftAlign()
	AdjustSAutoLeftAlign()

	SetPsIndex( idx int)
	GetPsIndex() int

	SetIndex(idx int)

	GetAlign() int
	SetAlign(al int)
	
	
	SetIconIndex(idx int)
	GetIconIndex() int

	Coord() (int, int)
	NewCoord(x,y int)
	Size() (int,int)
	NewSize(w,h int)
	
	
	UpdateIconNumbers()
	GetIconNumbers() int

	
	SetOnShow(on_show bool)
	GetOnShow() bool
	
	AppendIcon( it interface{} )
	ClearIcons()
	DrawIcons()

	GetMyList() []ListItemInterface
	
	GetName() string
	SetName(n string)
	GetFootMsg() [5]string

	KeyDown( ev *event.Event)

	ReturnToUpLevelPage()
	
	OnLoadCb()
	OnReturnBackCb()
	OnExitCb()

//	IconClick()
	ResetPageSelector()
	DrawPageSelector()

	ClearCanvas()
	Draw()

	
}

type Page struct {
	Widget
	Icons []IconItemInterface // slice ,use append
	IconNumbers int
	IconIndex int
	PrevIconIndex int
	
	Ps PageSelectorInterface
	PsIndex int

	Index int

	Align int
	
	CanvasHWND *sdl.Surface
	HWND       *sdl.Surface

	MyList []ListItemInterface
	
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
	p.EasingDur = 10

	p.Align = ALIGN["SLeft"]
	
	p.FootMsg = [5]string{"Nav.","","","","Enter"}
	
	return p
}

func (self *Page) GetScreen() *MainScreen {
	return self.Screen
}

func (self *Page) SetScreen(main_screen *MainScreen) {
	self.Screen = main_screen
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

		old_surf := it.GetImgSurf()
		
		it_w,it_h := it.Size() //width height changed by Adjust above
		it.SetImgSurf( transform.SmoothScale(old_surf,it_w,it_h) )
	}

	ps := NewPageSelector()
	ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
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

	if self.IconNumbers == 1 {
		start_x = self.Width/2
		start_y = self.Height/2
		it := self.Icons[0]
		it.SetParent(self)
		it.SetIndex(0)
		it.Adjust(start_x,start_y, IconWidth,IconHeight,0)
		/*
		old_surf := it.GetImgSurf()
		it_w,it_h := it.Size()
		it.SetImgSurf( transform.SmoothScale(old_surf, it_w,it_h))
		*/
	}else if self.IconNumbers == 2 {
		start_x = (self.Width - self.PageIconMargin - self.IconNumbers*IconWidth) / 2 + IconWidth/2
		start_y = self.Height /2

		for i:=0; i < self.IconNumbers; i++ {
			it := self.Icons[i]
			it.SetParent(self)
			it.SetIndex(i)
			it.Adjust( start_x+ i*self.PageIconMargin+i*IconWidth, start_y, IconWidth, IconHeight,0)
			/*
			old_surf := it.GetImgSurf()
			it_w,it_h := it.Size()
			it.SetImgSurf( transform.SmoothScale( old_surf, it_w,it_h))
			*/
		}
		
	}else if self.IconNumbers > 2 {
		for i:=0; i < self.IconNumbers; i++ {
			it := self.Icons[i]
			it.SetParent(self)
			it.SetIndex(i)
			it.Adjust(start_x+i*self.PageIconMargin + i*IconWidth, start_y, IconWidth, IconHeight, 0)
			/*
			old_surf := it.GetImgSurf()
			it_w,it_h := it.Size()
			it.SetImgSurf( transform.SmoothScale( old_surf, it_w,it_h))			
      */
		}
	}

	ps := NewPageSelector()
	ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
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



func (self *Page) InitLeftAlign() {
	self.PosX   = self.Index * self.Screen.Width
	self.Width  = self.Screen.Width
	self.Height = self.Screen.Height

	cols := int(self.Width/IconWidth)
	rows := int((self.IconNumbers * IconWidth) / self.Width +1)

	if rows < 1{
		rows = 1
	}
	cnt := 0
	start_x := 0
	start_y := 0
	
	for i:=0; i< rows; i++ {
		for j:=0; j< cols; j++ {
			start_x = IconWidth/2 + j*IconWidth
			start_y = TitleBar_BarHeight + IconHeight /2 + i*IconHeight

			icon := NewIconItem()
			icon.Init(start_x,start_y,IconWidth-4,IconHeight-4,0)
			icon.SetIndex(cnt)
			icon.SetParent(self)
			self.Icons = append(self.Icons, icon)
			if cnt >= (self.IconNumbers -1 ){
				break
			}
			cnt+=1
		}
	}

	ps := NewPageSelector()
	ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
	ps.Parent = self
	ps.Init(IconWidth/2,IconHeight/2,92,92,128)

	self.Ps = ps
	self.PsIndex = 0
	self.OnShow = false	
	
}

func (self *Page) Adjust() { // default init way,
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	start_x := 0
	start_y := 0
	
	if self.Align == ALIGN["HLeft"] {
		start_x = (self.Width - self.IconNumbers*IconWidth) / 2 + IconWidth/2
		start_y = self.Height/2

		for i:=0;i< self.IconNumbers; i++ {
			self.Icons[i].SetParent(self)
			self.Icons[i].SetIndex(i)
			self.Icons[i].Adjust(start_x + i*IconWidth, start_y, IconWidth, IconHeight,0)
		}

		ps := NewPageSelector()
		ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
		ps.Parent = self
		ps.Init(start_x,start_y, 92,92,128)
		self.Ps = ps
		self.PsIndex = 0
		self.OnShow = false
		
	}else if self.Align == ALIGN["SLeft"] {
		start_x = (self.PageIconMargin + IconWidth + self.PageIconMargin) / 2
		start_y = self.Height/2
		for i:=0;i< self.IconNumbers; i++ {
			it:=self.Icons[i]
			it.SetParent(self)
			it.SetIndex(i)
			it.Adjust(start_x + i*self.PageIconMargin+i*IconWidth, start_y, IconWidth, IconHeight,0)
		}
		ps := NewPageSelector()
		ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
		ps.Parent = self
		ps.Init(start_x,start_y-self.SelectedIconTopOffset, 92,92,128)
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
	
}

func (self *Page) GetOnShow() bool {
	return self.OnShow
}

func (self *Page) SetOnShow( on_show bool) {
	self.OnShow = on_show
}

func (self *Page) UpdateIconNumbers() {

	self.IconNumbers = len(self.Icons)
	
}

func (self *Page) GetIconNumbers() int {
	return self.IconNumbers
}

func (self *Page) Init() {

	if self.Screen != nil {
		if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
			self.CanvasHWND = self.Screen.CanvasHWND
		}
	}

	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	start_x := (self.Width - self.IconNumbers *IconWidth) /2 + IconWidth /2
	start_y := self.Height/2
	
	for i:=0; i< self.IconNumbers; i++ {
		it := NewIconItem()
		it.SetParent(self)
		it.SetIndex(i)
		it.Init(start_x + i * IconWidth, start_y, IconWidth,IconHeight, 0)
		self.Icons = append(self.Icons, it)
	}

	if self.IconNumbers > 0 {
		ps := NewPageSelector()
		ps.IconSurf = MyIconPool.GetImgSurf("blueselector")
		ps.Parent = self
		ps.Init(start_x,start_y, IconWidth+4, IconHeight+4, 128)
		self.Ps = ps
		self.PsIndex = 0
		self.OnShow = false
	}
}


func (self *Page) IconStepMoveData(icon_eh ,cuts int)  []int {  //  no Sine,No curve,plain movement steps data
	var all_pieces []int
	
	piece := float32( icon_eh / cuts )
	c := float32(0.0)
	prev := float32(0.0)
	for i:=0;i<cuts;i++ {
		c+= piece
		dx:= c-prev
		if dx < 0.5 {
			dx = float32(1.0)
		}
		dx += 0.9
		all_pieces = append(all_pieces, int(dx))
		if c >= float32(icon_eh) {
			break
		}
	}

	c = 0.0
	bidx := 0

	for _,v := range all_pieces {
		c += float32(v)
		bidx+=1
		if c >= float32(icon_eh) {
			break
		}
	}

	all_pieces = all_pieces[0:bidx]

	if len(all_pieces) < cuts {
		dff := cuts - len(all_pieces)
		var diffa []int
		for i:=0;i<dff;i++ {
			diffa= append(diffa,0)
		}
		
		all_pieces = append(all_pieces, diffa...)
	}

	return all_pieces		
}

func (self *Page) EasingData(start,distance int) []int {
	current_time := float32(0.0)
	start_posx   := float32(0.0)
	current_posx := start_posx
	final_posx   := float32(distance)
//	posx_init    := start
	dur          := self.EasingDur
	last_posx    := float32(0.0)

	var all_last_posx []int

	for i:=0;i<distance*dur;i++ {
		current_posx = float32(easings.SineIn(float32(current_time), float32(start_posx), float32(final_posx-start_posx),float32(dur)))
		if current_posx >= final_posx {
			current_posx = final_posx
		}
		dx := current_posx - last_posx
		all_last_posx = append(all_last_posx,int(dx))
		current_time+=1.0
		last_posx = current_posx
		if current_posx >= final_posx {
			break
		}
	}

	c := 0
	for _,v := range all_last_posx {
		c+=v
	}
	if c < int(final_posx - start_posx) {
		all_last_posx = append(all_last_posx, int( int(final_posx) - c ))
	}

	return all_last_posx	
}


func (self *Page) IconSmoothUp(icon_ew int) {
	data := self.EasingData(self.PosX,icon_ew)
	data2 := self.IconStepMoveData(self.SelectedIconTopOffset, len(data))

	for i,_ := range data {
		self.ClearCanvas()
		cur_icon_x,cur_icon_y := self.Icons[self.IconIndex].Coord()
		self.Icons[self.IconIndex].NewCoord(cur_icon_x, cur_icon_y - data2[i] )
		
		prev_icon_x,prev_icon_y := self.Icons[self.PrevIconIndex].Coord()
		
		if prev_icon_y < self.Height/2 {
			self.Icons[self.PrevIconIndex].NewCoord(prev_icon_x, prev_icon_y + data2[i])

			self.DrawIcons()
			self.Screen.SwapAndShow()
		}
	}
}

func (self *Page) IconsEasingLeft(icon_ew int) {
	data := self.EasingData(self.PosX, icon_ew)
	data2 := self.IconStepMoveData(self.SelectedIconTopOffset, len(data))

	for i,v := range data {
		self.ClearCanvas()
		
		self.PosX -= v
		
		cur_icon_x,cur_icon_y := self.Icons[self.IconIndex].Coord()
		self.Icons[self.IconIndex].NewCoord(cur_icon_x, cur_icon_y - data2[i] )

		prev_icon_x,prev_icon_y := self.Icons[self.PrevIconIndex].Coord()
		if prev_icon_y < self.Height/2 {
			self.Icons[self.PrevIconIndex].NewCoord(prev_icon_x, prev_icon_y + data2[i])
		}
		self.DrawIcons()
		self.Screen.SwapAndShow()
	}
}


func (self *Page) IconsEasingRight(icon_ew int) {
	data := self.EasingData(self.PosX, icon_ew)
	data2 := self.IconStepMoveData(self.SelectedIconTopOffset, len(data))

	for i,v := range data {
		self.ClearCanvas()
		
		self.PosX += v
		
		cur_icon_x,cur_icon_y := self.Icons[self.IconIndex].Coord()
		self.Icons[self.IconIndex].NewCoord(cur_icon_x, cur_icon_y - data2[i] )

		prev_icon_x,prev_icon_y := self.Icons[self.PrevIconIndex].Coord()
		if prev_icon_y < self.Height/2 {
			self.Icons[self.PrevIconIndex].NewCoord(prev_icon_x, prev_icon_y + data2[i])
		}
		self.DrawIcons()
		self.Screen.SwapAndShow()
	}
}

func (self *Page) EasingLeft(ew int) {
	data := self.EasingData(self.PosX,ew)

	for _, i := range data {
		self.PosX -= i
		self.Draw()
		self.Screen.SwapAndShow()
	}
}


func (self *Page) EasingRight(ew int) {
	data := self.EasingData(self.PosX,ew)

	for _, i := range data {
		self.PosX += i
		self.Draw()
		self.Screen.SwapAndShow()
	}
}

func (self *Page) MoveLeft(ew int) {
	self.PosX -= ew
}

func (self *Page) MoveRight(ew int) {
	self.PosX += ew
}

func (self *Page) ResetPageSelector() {
	self.PsIndex = 0
	self.IconIndex = 0
	self.Ps.SetOnShow(true)
}

func (self *Page) DrawPageSelector() {
	if self.Ps.GetOnShow() == true {
//		fmt.Println("DrawPageSelector")
		self.Ps.Draw()
	}
}

func (self *Page) MoveIconIndexPrev() bool {
	self.IconIndex -= 1
	if self.IconIndex < 0 {
		self.IconIndex = 0
		self.PrevIconIndex = self.IconIndex
		return false
	}

	self.PrevIconIndex = self.IconIndex + 1
	return true
}

func (self *Page) MoveIconIndexNext() bool {
	self.IconIndex+=1
	if self.IconIndex > (self.IconNumbers - 1) {
		self.IconIndex = self.IconNumbers -1
		self.PrevIconIndex = self.IconIndex
		return false
	}
	self.PrevIconIndex = self.IconIndex - 1
	return true
}


func (self *Page) IconClick() {
	if self.IconIndex > ( len(self.Icons) - 1) {
		return
	}

	cur_icon := self.Icons[self.IconIndex]

	if self.Ps.GetOnShow() == false {
		return
	}

	if cur_icon.GetMyType() == ICON_TYPES["EXE"] {
		fmt.Printf("IconClick: %s %d", cur_icon.GetCmdPath(), cur_icon.GetIndex() )
		self.Screen.RunEXE(cur_icon.GetCmdPath())
		return
	}

	if cur_icon.GetMyType() == ICON_TYPES["DIR"] {
		child_page := cur_icon.GetLinkPage()
		if child_page != nil {
			self.Screen.PushPage(child_page)
			child_page.Draw()
		}
		return
	}

	if cur_icon.GetMyType() == ICON_TYPES["FUNC"] {
		invoker := cur_icon.GetCmdInvoke()
		if invoker != nil {
			invoker.Run(self.Screen)
		}
		return
	}
}

func (self *Page) ReturnToUpLevelPage() {
	pop_page := self.Screen.MyPageStack.Pop()
	if pop_page != nil {
		page_ := pop_page.(PageInterface)
		page_.Draw()
		self.Screen.SetCurPage(page_)
	}else {
		if self.Screen.MyPageStack.Length() == 0 {
			if len(self.Screen.Pages) > 0 {
				if self.Screen.PageIndex < len(self.Screen.Pages) {
					self.Screen.CurrentPage = self.Screen.Pages[ self.Screen.PageIndex ]
					self.Screen.CurrentPage.Draw()
					fmt.Println( "OnTopLevel", self.Screen.PageIndex)
				}
			}
		}
	}
}

func (self *Page) ClearCanvas() {
	surface.Fill(self.CanvasHWND, self.Screen.SkinManager.GiveColor("White"))
}

func (self *Page) AppendIcon( it interface{} ) {
	self.Icons = append(self.Icons, it.(IconItemInterface))
}

func (self *Page) GetIcons() []IconItemInterface {
	return self.Icons
}

func (self *Page) ClearIcons() {
	for i:=0;i<self.IconNumbers; i++ {
		self.Icons[i].Clear()
	}
}

func (self *Page) DrawIcons() {
	for i:=0;i<self.IconNumbers; i++ {
		self.Icons[i].Draw()
	}	
}

func (self *Page) GetMyList() []ListItemInterface {
	return self.MyList
}

func (self *Page) KeyDown( ev *event.Event) {
	if ev.Data["Key"] == CurKeys["A"] {

		if self.FootMsg[3] == "Back" {
			self.ReturnToUpLevelPage()
			self.Screen.Draw()
			self.Screen.SwapAndShow()
			return
		}
	}

	if ev.Data["Key"] == CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == CurKeys["Right"] {
		if self.MoveIconIndexNext() == true {
			if self.IconIndex == (self.IconNumbers -1) || self.PrevIconIndex == 0 {
				self.IconSmoothUp(IconWidth + self.PageIconMargin)
			}else {
				self.IconsEasingLeft(IconWidth + self.PageIconMargin)
			}
			
			self.PsIndex = self.IconIndex
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}
	}

	if ev.Data["Key"] == CurKeys["Left"] {
		if self.MoveIconIndexPrev() == true {
			if self.IconIndex == 0 || self.PrevIconIndex == (self.IconNumbers -1) {
				self.IconSmoothUp(IconWidth + self.PageIconMargin)
			}else {
				self.IconsEasingRight(IconWidth + self.PageIconMargin)
			}
			self.PsIndex = self.IconIndex
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}
	}

	if ev.Data["Key"] == CurKeys["Enter"] {
		self.IconClick()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
	
}


func (self *Page) OnLoadCb() {
	
}

func (self *Page) OnReturnBackCb() {
	
}

func (self *Page) OnExitCb() {
	
}

func (self *Page) Draw() {
	self.ClearCanvas()
	self.DrawIcons()
	self.DrawPageSelector()
}

func (self *Page) GetFootMsg() [5]string {
	return self.FootMsg
}

func (self *Page) SetFootMsg(footmsg [5]string) {
	self.FootMsg = footmsg
}


func (self *Page) GetCanvasHWND() *sdl.Surface {
	return self.CanvasHWND
}

func (self *Page)	SetCanvasHWND( canvas *sdl.Surface) {
	self.CanvasHWND = canvas
}


func (self *Page)	GetHWND() *sdl.Surface {
	return self.HWND
}

func (self *Page)	SetHWND(h *sdl.Surface) {
	self.HWND = h
}

func (self *Page) SetPsIndex( idx int) {
	self.PsIndex = idx
}

func (self *Page) GetPsIndex() int {
	return self.PsIndex
}

func (self *Page) SetIconIndex( idx int) {
	self.IconIndex = idx
}

func (self *Page) GetIconIndex() int {
	return self.IconIndex
}


func (self *Page) GetName() string {
	return self.Name
}

func (self *Page) SetName(n string) {
	self.Name = n
}

func (self *Page) SetIndex(idx int) {
	self.Index = idx
}

func (self *Page) SetAlign(al int) {
	inthere := false
	for _,v := range ALIGN {
		if v == al {
			inthere = true
			break
		}
	}

	if inthere {
		self.Align = al
	}
}

func (self *Page) GetAlign() int {
	return self.Align
}

