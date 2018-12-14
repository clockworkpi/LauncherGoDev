package Bluetooth

type BleForgetConfirmPage struct {
  
  UI.ConfirmPage


}

func NewBleForgetConfirmPage() *BleForgetConfirmPage {
  p := &BleForgetConfirmPage{}
  
  p.ConfirmText = "Confirm Forget?"
  p.UI.ConfirmPage.ConfirmText = p.ConfirmText
  
  return p
}

func (self *BleForgetConfirmPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == CurKeys["A"] || ev.Data["Key"] == CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	} 
  
  if ev.Data["Key"] == CurKeys["B"] {
    self.SnapMsg("Deleting")
    self.Screen.Draw()
    self.Screen.SwapAndShow()
    
    
    time.BlockDelay(400)
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()    
    
  }
}

func (self *BleForgetConfirmPage) Draw() {
  self.DrawBG()
  for _,v := range self.MyList{
    v.Draw()
  }
}


type BleInfoPageSelector struct {
  UI.InfoPageSelector
}

func NewBleInfoPageSelector() *BleInfoPageSelector{
  p := &BleInfoPageSelector{}

  return p 
}

func (self *BleInfoPageSelector) Draw() {
  idx := self.Parent.GetPsIndex()
  mylist := self.Parent.GetMyList()
 
  if idx < len(mylist) {
    _,y := mylist[idx].Coord()
    _,h := mylist[idx].Size()
    
    x := self.PosX+2
    self.PosY = y+1
    self.Height = h-3
    
    canvas_ := self.Parent.GetCanvasHWND()
    rect_   := rect.Rect(x,self.PosY,self.Width-4, self.Height)
    
    draw.AARoundRect(canvas_,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  }
}


type BleInfoPage struct {
  UI.Page
  
  ListFontObj *ttf.Font
  ListSmFontObj *ttf.Font
  ListSm2FontObj *ttf.Font
  
  MyList []UI.ListItemInterface
  
  ConfirmPage1 *BleForgetConfirmPage
  
}


func NewBleInfoPage() *BleInfoPage {
  p :=&BleInfoPage{}
  
  p.FootMsg = [5]string{"Nav","Disconnect","Forget","Back","" }  
  
  p.ListFontObj = UI.Fonts["varela15"]
  p.ListSmFontObj = UI.Fonts["varela12"]
  p.ListSm2FontObj = UI.Fonts["varela11"]
  
  return p
}


func (self *BleInfoPage) Init() {
  
  if self.Screen != nil {
    if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
      self.CanvasHWND = self.Screen.CanvasHWND
    }
  }
     
  self.PosX = self.Index*self.Screen.Width 
  self.Width = self.Screen.Width //  equals to screen width
  self.Height = self.Screen.Height
        
  ps := NewBleInfoPageSelector()
  ps.Parent = self
  self.Ps = ps
  self.PsIndex = 0  
  
  self.GenList()
  
  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = 2
  self.Scroller.PosY = 2
  self.Scroller.Init()
        
  self.ConfirmPage1 = BleForgetConfirmPage()
  self.ConfirmPage1.Screen = self.Screen
  self.ConfirmPage1.Name   = "Confirm Forget"
  self.ConfirmPage1.Parent = self
  self.ConfirmPage1.Init()   

}

func (self *BleInfoPage) GenList() {
  if len(self.AList) == 0 {
    return
  }
  
  
  self.MyList = nil
  
  self.PsIndex = 0
  
  start_x := 0 
  start_y := 0
  
  i := 0
  for k,v := range self.AList {
    li := UI.NewInfoPageListItem()
    li.Parent = self
    li.PosX   = start_x
    li.PosY   = start_y +i*NetItemDefaultHeight
    li.Width  = UI.Width
    
    li.Fonts["normal"] = self.ListFontObj
    if k =="UUIDs" {
      li.Fonts["small"] = self.ListSm2FontObj
    }else{
      li.Fonts["small"] = self.ListSmFontObj
    }
    
    li.Init(k)
    li.Flag = k
    
  }

}

func (self *BleInfoPage) ScrollUp() {
  if len(self.MyList) == 0 {
    return
  }
  
  self.PsIndex -= 1
  
  if self.PsIndex < 0 {
    self.PsIndex = 0
  }
  
  cur_li = self.MyList[self.PsIndex]
  
  x,y := cur_li.Coord()
  
  if y < 0 {
    for i,v := range self.MyList {
      x,y = v.Coord()
      _,h := v.Size()
      self.MyList[i].NewCoord(x,y+h)
    }
  }

}

func (self *BleInfoPage) ScrollDown() {
  

}

