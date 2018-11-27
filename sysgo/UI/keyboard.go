package UI

import (
  "fmt"
	"strings"
  
  "github.com/cuu/gogame/font"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"


	"github.com/cuu/LauncherGo/sysgo/easings"

)
//sysgo/UI/keyboard_keys.layout
type KeyboardIcon struct {
	TextItem // IconItem->TextItem->KeyboardIcon
}

func NewKeyboardIcon() *KeyboardIcon {
	p := &KeyboardIcon{}

	p.Color = &color.Color{83,83,83,255}//SkinManager().GiveColor('Text')
	
	p.MyType = ICON_TYPES["NAV"]

	return p
}

func (self *KeyboardIcon) Draw() {

	rect_ := draw.MidRect(self.PosX,self.PosY,self.Width,self.Height,Width,Height)
	
	surface.Blit(self.Parent.GetCanvasHWND(),self.ImgSurf,rect_,nil)
	
}


type KeyboardSelector struct {
	PageSelector
	Parent *Keyboard
}


func NewKeyboardSelector() * KeyboardSelector {
	p := &KeyboardSelector{}

	return p
}

func (self *KeyboardSelector) Draw() {
	sec_idx := self.Parent.SectionIndex
	row_idx := self.Parent.RowIndex
	idx     := self.Parent.PsIndex

	x, y    := self.Parent.SecsKeys[sec_idx][row_idx][idx].Coord()
	w, h    := self.Parent.SecsKeys[sec_idx][row_idx][idx].Size()

	rect_   := draw.MidRect(x,y,w+6,h+1,self.Parent.Width,self.Parent.Height)

	if rect_.W <= 0 || rect_.H <= 0 {
		return
	}
	
	color_ := &color.Color{126,206,244,255}
	draw.AARoundRect(self.Parent.CanvasHWND,rect_,color_,3,0,color_)

}

type Keyboard struct {
	Page

	Secs map[int][][]string
	
	SecsKeys map[int][][]TextItemInterface
	
	SectionNumbers int
	SectionIndex int
	Icons  map[string]IconItemInterface

	KeyboardLayoutFile string ///sysgo/UI/keyboard_keys.layout

	LeftOrRight int

	RowIndex int

	Textarea *Textarea
	Selector *KeyboardSelector

	
}

func NewKeyboard() *Keyboard {
	p := &Keyboard{}

	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20

	p.Align = ALIGN["SLeft"]
  
	p.EasingDur = 10
	
	p.SectionNumbers = 3
	p.SectionIndex = 1

	p.Icons =  make( map[string]IconItemInterface )

	p.LeftOrRight = 1

	p.RowIndex = 0
	
	p.FootMsg = [5]string{"Nav.","ABC","Done","Backspace","Enter"}

	p.Secs = make(map[int][][]string)
	p.SecsKeys = make(map[int][][]TextItemInterface)
	
	p.KeyboardLayoutFile = "sysgo/UI/keyboard_keys.layout"

	
	return p
	
}

func (self *Keyboard) ReadLayoutFile( fname string) {

	
	LayoutIndex := 0

	content ,err := ReadLines(fname)

	Assert(err)

	var tmp [][]string
	for i, v := range content {
		content[i] = strings.TrimSpace(v)

		stmp := strings.Split(content[i], " ")
		for j, u := range stmp {
			stmp[j] = strings.TrimSpace(u)
		}

		tmp = append(tmp, stmp)
	}

	for _, v := range tmp {
		if len(v) > 2 {

			if _, ok := self.Secs[LayoutIndex]; ok {
				self.Secs[LayoutIndex] = append(self.Secs[LayoutIndex], v)
			} else {

				self.Secs[LayoutIndex] = [][]string{}
				self.Secs[LayoutIndex] = append(self.Secs[LayoutIndex], v)

			}

		} else { //empty []
			LayoutIndex += 1
		}
	}	
}


func (self *Keyboard) SetPassword(pwd string) {
	pwd_seq_list := strings.SplitAfter(pwd,"")

	self.Textarea.ResetMyWords()
	for _,v := range pwd_seq_list {
		self.Textarea.AppendText(v)
	}
}


func (self *Keyboard) Init() {
	self.CanvasHWND = self.Screen.CanvasHWND
	self.ReadLayoutFile(self.KeyboardLayoutFile) //assign to self.Secs

	self.SectionNumbers = len(self.Secs)

	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	fontobj := Fonts["veramono24"]

	word_margin := 15

	secs_zero := strings.Join(self.Secs[0][0],"")
	fw,_:= font.Size(fontobj,secs_zero)

	start_x := (self.Width - fw - len(self.Secs[0][0])*word_margin)/2+word_margin/2
	start_y := 0

//	cnt := 0

	for i:=0; i < self.SectionNumbers; i++ {
		self.SecsKeys[i] = [][]TextItemInterface{}
		for j:=0; j < len(self.Secs[i]); j++ {
			self.SecsKeys[i] = append(self.SecsKeys[i],[]TextItemInterface{})
			secs_ij := strings.Join(self.Secs[i][j],"")
			fw,_ := font.Size(fontobj,secs_ij)
			start_x = (self.Width-fw- len(self.Secs[i][j])*word_margin)/2+word_margin/2
			start_x = start_x + i*self.Width

			start_y = 84 + j * (word_margin+14)

			for _,val := range self.Secs[i][j] {
				ti := NewTextItem()
				ti.FontObj = fontobj
				ti.Parent = self

				if val == "_L" || val == "_R" {
					it := NewKeyboardIcon()
					it.ImgSurf = MyIconPool.GetImgSurf(val)
					it.Parent = self
					it.Str = val
					it.Init(start_x+surface.GetWidth(it.ImgSurf)/2,start_y,surface.GetWidth(it.ImgSurf),surface.GetHeight(it.ImgSurf),0)
					self.SecsKeys[i][j] = append(self.SecsKeys[i][j],it)
					self.IconNumbers += 1
					start_x = start_x + surface.GetWidth(it.ImgSurf)+word_margin
				}else {

					if val ==  "_S" {
						val = "Space"
						ti.FontObj = Fonts["veramono15"]
						ti.Bold = true
					}

					cur_alpha_w,cur_alpha_h := font.Size(ti.FontObj,val)
					ti.Init(start_x + cur_alpha_w/2,start_y,cur_alpha_w,cur_alpha_h,0)
					ti.Str = val
					start_x = start_x + cur_alpha_w+word_margin // prepare for next alphabet
					self.SecsKeys[i][j] = append(self.SecsKeys[i][j],ti)
				}
			}
		}
	}

	self.SectionIndex = 0

	self.Textarea = NewTextarea()

	self.Textarea.PosX = 4
	self.Textarea.PosY = 4
	self.Textarea.Width = self.Width - 4*2
	self.Textarea.Height = 60

	self.Textarea.CanvasHWND = self.CanvasHWND
	self.Textarea.Init()

	ps := NewKeyboardSelector()

	ps.Parent = self
	ps.Init(start_x,start_y,25,25,128)
  ps.OnShow = true
  
	self.Ps = ps
	self.PsIndex = 0

}

func (self *Keyboard) SelectUpChar() {
	sec_idx := self.SectionIndex

	self.RowIndex -=1
	if self.RowIndex < 0 {
		self.RowIndex = len(self.SecsKeys[sec_idx])-1
	}

	if self.PsIndex >= len(self.SecsKeys[sec_idx][self.RowIndex]) {
		self.PsIndex = len(self.SecsKeys[sec_idx][self.RowIndex])-1
	}

	self.ClearCanvas()
	self.Draw()
	self.Screen.SwapAndShow()
}

func (self *Keyboard) SelectDownChar() {
	sec_idx := self.SectionIndex

	self.RowIndex += 1

	if self.RowIndex >= len(self.SecsKeys[sec_idx]) {
		self.RowIndex = 0
	}

	if self.PsIndex >=len(self.SecsKeys[sec_idx][self.RowIndex]) {
		self.PsIndex = len(self.SecsKeys[sec_idx][self.RowIndex])-1
	}

	self.ClearCanvas()
	self.Draw()
	self.Screen.SwapAndShow()
}

func (self *Keyboard) SelectNextChar() {

	sec_idx := self.SectionIndex
	row_idx := self.RowIndex
	self.PsIndex+=1
	
	if self.PsIndex >= len(self.SecsKeys[sec_idx][row_idx]) {
		self.PsIndex = 0
		self.RowIndex+=1
	
		if self.RowIndex >= len(self.SecsKeys[sec_idx]) {
			self.RowIndex = 0
		}

	}
	
	self.ClearCanvas()
	self.Draw()
	self.Screen.SwapAndShow()
	
}

func (self *Keyboard) SelectPrevChar() {

	sec_idx := self.SectionIndex    
	self.PsIndex-=1
	if self.PsIndex < 0 {
		self.RowIndex-=1
		if self.RowIndex <=0 {
			self.RowIndex = len(self.SecsKeys[sec_idx])-1
		}
		self.PsIndex = len(self.SecsKeys[sec_idx][self.RowIndex]) -1
	}

	self.ClearCanvas()
	self.Draw()
	self.Screen.SwapAndShow()
}

func (self *Keyboard) ClickOnChar() {
	sec_idx := self.SectionIndex        
	alphabet := self.SecsKeys[sec_idx][self.RowIndex][self.PsIndex].GetStr()
  
	if alphabet == "Space"{
		alphabet = " "
	}

	if alphabet == "_L" || alphabet == "_R" {
		if alphabet == "_L" {
			self.Textarea.SubTextIndex()
		}else if alphabet == "_R"{
			self.Textarea.AddTextIndex()
		}
	}else {
		self.Textarea.AppendText(alphabet)
	}

	self.Textarea.Draw()
	self.Screen.SwapAndShow()
}

func (self *Keyboard) KeyboardShift() {
	distance := self.Width //320
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

	for _,v := range all_last_posx {
		for j:=0;j<self.SectionNumbers;j++ {
			for _,u := range self.SecsKeys[j] {
				for _,x := range u {
          x_,y_ := x.Coord()
          x.NewCoord(x_+self.LeftOrRight*v,y_)
				}
			}
		}

		self.ResetPageSelector()
		self.ClearCanvas()
		self.Draw()
		self.Screen.SwapAndShow()
	}
}

func (self *Keyboard) ShiftKeyboardPage() {
	self.KeyboardShift()
	self.SectionIndex -= self.LeftOrRight
	self.Draw()
	self.Screen.SwapAndShow()
}


func (self *Keyboard) KeyDown( ev *event.Event) {
	if ev.Data["Key"] == CurKeys["Up"] {
		self.SelectUpChar()
		return
	}

	if ev.Data["Key"] == CurKeys["Down"] {
		self.SelectDownChar()
		return
	}

	if ev.Data["Key"] == CurKeys["Right"] {
		self.SelectNextChar()
		return
	}

	if ev.Data["Key"] == CurKeys["Left"] {
		self.SelectPrevChar()
		return
	}

	if ev.Data["Key"] == CurKeys["B"] || ev.Data["Key"] == CurKeys["Enter"] {
		self.ClickOnChar()
		return
	}
  
  if ev.Data["Key"] == CurKeys["X"] {
    if self.SectionIndex <= 0 {
      self.LeftOrRight = -1
    }
    
    if self.SectionIndex >= (self.SectionNumbers - 1) {
      self.LeftOrRight = 1
    }
    
    self.ShiftKeyboardPage()
    
  }

	if ev.Data["Key"] == CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
    
	}

	if ev.Data["Key"] == CurKeys["Y"] { // done
		fmt.Println(strings.Join(self.Textarea.MyWords,""))
		self.ReturnToUpLevelPage()
		self.Screen.SwapAndShow()
		//Uplevel/Parent page invoke OnReturnBackCb,eg: ConfigWireless
		
	}

	if ev.Data["Key"] == CurKeys["A"] {
		self.Textarea.RemoveFromLastText()
		self.Textarea.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == CurKeys["LK1"] {
		if self.SectionIndex < self.SectionNumbers -1 {
			self.LeftOrRight = -1
			self.ShiftKeyboardPage()
		}
	}

	if ev.Data["Key"] == CurKeys["LK5"] {
		if self.SectionIndex > 0 {
			self.LeftOrRight = 1
			self.ShiftKeyboardPage()
		}
	}

}

func (self *Keyboard) Draw() {
	self.ClearCanvas()
	self.Ps.Draw()

	for i:=0; i < self.SectionNumbers; i++ {
		for _,j := range self.SecsKeys[i] {
			for _,u := range j {
				u.Draw()
			}
		}
	}

	self.Textarea.Draw()
}
