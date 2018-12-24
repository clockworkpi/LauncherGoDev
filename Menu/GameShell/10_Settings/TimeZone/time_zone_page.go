package TimeZone

import (
  "fmt"
  "os/exec"
  
  "path/filepath"
  
  "github.com/veandco/go-sdl2/ttf"

  "github.com/cuu/gogame/draw"
  "github.com/cuu/gogame/rect"
  "github.com/cuu/gogame/color"
  "github.com/cuu/gogame/event"
  "github.com/cuu/gogame/time"

  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"

)
var TimeZonePath = "/usr/share/zoneinfo/posix"

type ListPageSelector struct {
  UI.InfoPageSelector
}

func NewListPageSelector() *ListPageSelector {
  
  p := &ListPageSelector{}
  
  p.Width = UI.Width
  p.BackgroundColor = &color.Color{131,199,219,255} //SkinManager().GiveColor('Front')
  
  return p

}


func (self *ListPageSelector) Draw() {
  
  idx := self.Parent.GetPsIndex()
  mylist := self.Parent.GetMyList()
 
  if idx < len(mylist) {
    x,y := mylist[idx].Coord()
    _,h := mylist[idx].Size()
    
    self.PosX = x+2
    self.PosY = y+1
    self.Height = h-3
    
    canvas_ := self.Parent.GetCanvasHWND()
    rect_   := rect.Rect(self.PosX,self.PosY,self.Width-4, self.Height)
    
    draw.AARoundRect(canvas_,&rect_,self.BackgroundColor,4,0,self.BackgroundColor)
  } 
}



type TimeZoneListPage struct {
  UI.Page
  
  Scroller *UI.ListScroller
  
  Icons  map[string]UI.IconItemInterface

  ListFont *ttf.Font
  MyStack *UI.FolderStack
  BGpng *UI.IconItem
  BGwidth int
  BGheight int
  
  SwapMyList []UI.ListItemInterface
}

type ListEle struct {
  
  Name string
  FilePath string
  IsFile bool
}




func NewTimeZoneListPage() *TimeZoneListPage {
  p := &TimeZoneListPage{}
  
  p.BGwidth = 56
  p.BGheight = 70
  
  p.FootMsg = [5]string{ "Nav","","","Back","Select" }
  
  p.ListFont = UI.Fonts["notosanscjk15"]
  
  p.MyStack = UI.NewFolderStack()
  p.MyStack.SetRootPath( TimeZonePath )
  
  p.Icons = make(map[string]UI.IconItemInterface )
  return p
}

func (self *TimeZoneListPage) GetMapIcons() map[string]UI.IconItemInterface {
  return self.Icons
}

func (self *TimeZoneListPage) buildDirectoryList(path string) []*ListEle  {
  
  //[*ListEle{},*ListEle{}]
  var ret []*ListEle 
  
  file_paths,err := filepath.Glob(path+"/*")//sorted
  if err == nil {
    for _, u := range file_paths {
      e := &ListEle{}
      e.Name = filepath.Base(u)
      e.FilePath = u
      if UI.IsAFile(u) {
        e.IsFile = true
      }else {
        e.IsFile = false
      }
      
      ret = append(ret,e)
    }
  }
  return ret
}

func (self *TimeZoneListPage) SyncList(path string) {
  
  alist := self.buildDirectoryList(path)
  if len(alist) == 0 {
    fmt.Println("buildDirectoryList empty")
    return
  }
  
  self.MyList = nil
  self.SwapMyList = nil
  
  
  start_x := 0
  start_y := 0
  
  hasparent := 0
  
  if self.MyStack.Length() > 0 {
    hasparent = 1
    
    li := NewTimeZoneListPageListItem()
    li.Parent = self
    li.PosX = start_x
    li.PosY = start_y
    li.Width  = UI.Width
    li.Fonts["normal"] = self.ListFont
    li.MyType = UI.ICON_TYPES["DIR"]
    li.Init("[..]")
    
    self.MyList = append(self.MyList,li)
  
  }
  
  for i,v := range alist{
    li := NewTimeZoneListPageListItem()
    li.Parent = self
    li.PosX = start_x
    li.PosY = start_y + (i+hasparent) *TimeZoneListPageListItemDefaultHeight
    li.Width  = UI.Width
    li.Fonts["normal"] = self.ListFont
    li.MyType = UI.ICON_TYPES["FILE"]    
    
    if v.IsFile == false {
      li.MyType = UI.ICON_TYPES["DIR"]
    }else{
      li.MyType = UI.ICON_TYPES["FILE"]
    }
    
    li.Init(v.Name)
    li.Path = v.FilePath
    
    self.MyList = append(self.MyList,li)
  }
  
  
  for _,v := range self.MyList {
    self.SwapMyList = append(self.SwapMyList,v)
  }

}

func (self *TimeZoneListPage) Init() {
  self.PosX = self.Index * self.Screen.Width
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height

  self.CanvasHWND = self.Screen.CanvasHWND
  
  ps := NewListPageSelector()
  ps.Parent = self
  
  self.Ps = ps
  self.PsIndex = 0
  
  self.SyncList( TimeZonePath )
  
  icon_for_list := UI.NewMultiIconItem()
  icon_for_list.ImgSurf = UI.MyIconPool.GetImgSurf("sys")
  icon_for_list.MyType = UI.ICON_TYPES["STAT"]
  icon_for_list.Parent = self
  
  icon_for_list.Adjust(0,0,18,18,0)        
  self.Icons["sys"] = icon_for_list
  
  bgpng := UI.NewIconItem()
  bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("empty")
  bgpng.MyType = UI.ICON_TYPES["STAT"]
  bgpng.Parent = self
  bgpng.AddLabel("No timezones found on system!", UI.MyLangManager.TrFont("varela22"))
  bgpng.SetLabelColor( UI.MySkinManager.GiveColor("Disabled") )
  bgpng.Adjust(0,0,self.BGwidth,self.BGheight,0)
  
  self.BGpng = bgpng
  
  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = self.Width - 10
  self.Scroller.PosY = 2
  self.Scroller.Init()  

}


func (self *TimeZoneListPage) Click() {
  if len(self.MyList) == 0 {
    return
  }
  
  cur_li := self.MyList[self.PsIndex].(*TimeZoneListPageListItem)
  
  if cur_li.MyType == UI.ICON_TYPES["DIR"] {
    if cur_li.Path == "[..]" {
      self.MyStack.Pop()
      self.SyncList(self.MyStack.Last())
      self.PsIndex = 0
    }else {
      self.MyStack.Push( self.MyList[self.PsIndex].(*TimeZoneListPageListItem).Path)
      self.SyncList(self.MyStack.Last())
      self.PsIndex = 0
    }
    
  }
  
  if cur_li.MyType == UI.ICON_TYPES["FILE"] { //set the current timezone
    self.Screen.MsgBox.SetText("Applying")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow()  
    time.BlockDelay(300)
    cpCmd := exec.Command("sudo","cp", cur_li.Path,"/etc/localtime")
    err := cpCmd.Run()
    if err != nil{
      fmt.Println(err)
    }else {
      
      self.Screen.TitleBar.UpdateTimeLocation()
    
    }
    fmt.Println("add ",cur_li.Path)
  }
  
  self.Screen.Draw()
  self.Screen.SwapAndShow() 
}


func (self *TimeZoneListPage) Rescan() {
  self.SyncList(TimeZonePath)
  self.PsIndex = 0
}

func (self *TimeZoneListPage) KeyDown(ev *event.Event) {

  if ev.Data["Key"] == UI.CurKeys["Menu"] || ev.Data["Key"] == UI.CurKeys["A"] {
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  if ev.Data["Key"] == UI.CurKeys["Up"] {
    self.ScrollUp()
    self.Screen.Draw()
    self.Screen.SwapAndShow()    
  }
  
  if ev.Data["Key"] == UI.CurKeys["Down"] {
    self.ScrollDown()
    self.Screen.Draw()
    self.Screen.SwapAndShow()    
  }
  
  if ev.Data["Key"] == UI.CurKeys["Right"] {
    self.FastScrollDown(5)
    self.Screen.Draw()
    self.Screen.SwapAndShow()    
  }
  
  if ev.Data["Key"] == UI.CurKeys["Left"] {
    self.FastScrollUp(5)
    self.Screen.Draw()
    self.Screen.SwapAndShow()    
  }
  
  if ev.Data["Key"] == UI.CurKeys["Enter"] {
    self.Click()
  }
}

func (self *TimeZoneListPage) Draw() {
  self.ClearCanvas()
  
  if len(self.MyList) == 0 {
    self.BGpng.NewCoord(self.Width/2,self.Height/2)
    self.BGpng.Draw()
  }
  
  if len(self.MyList) *TimeZoneListPageListItemDefaultHeight > self.Height {
    
    self.Ps.(*ListPageSelector).Width = self.Width - 11
    self.Ps.Draw()
    
    for _,v := range self.MyList {
      if v.(*TimeZoneListPageListItem).PosY > self.Height + self.Height/2 {
        break
      }
      if v.(*TimeZoneListPageListItem).PosY < 0 {
        continue
      }
      
      v.Draw()
    }
    self.Scroller.UpdateSize( len(self.MyList)*TimeZoneListPageListItemDefaultHeight,
                            self.PsIndex*TimeZoneListPageListItemDefaultHeight)
    self.Scroller.Draw()
  
  }else {
    self.Ps.(*ListPageSelector).Width = self.Width
    self.Ps.Draw()
    
    for _,v := range self.MyList {
      if v.(*TimeZoneListPageListItem).PosY > self.Height + self.Height/2 {
        break
      }
      if v.(*TimeZoneListPageListItem).PosY < 0 {
        continue
      }
      
      v.Draw()
    }
  }
  
  
}



