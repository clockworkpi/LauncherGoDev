package Emulator

import (
  "fmt"
  "os"
  "strings"
  "path/filepath"
  "errors"
  
  "github.com/veandco/go-sdl2/ttf"
  
  "github.com/cuu/gogame/event"

  "github.com/cuu/gogame/color"
	"github.com/cuu/gogame/time"
  "github.com/cuu/LauncherGoDev/sysgo/UI"

)

type FavListPage struct {
  UI.Page
  Icons  map[string]UI.IconItemInterface
  ListFont *ttf.Font
  MyStack *EmuStack
  EmulatorConfig *ActionConfig
  
  RomSoConfirmDownloadPage *RomSoConfirmPage
  
  MyList []UI.ListItemInterface
  BGwidth int
  BGheight int //70
  Scroller *UI.ListScroller
  Scrolled int 
  Leader *MyEmulator
  
}

func NewFavListPage() *FavListPage {
  p := &FavListPage{}
  p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]
	
	p.FootMsg = [5]string{ "Nav","Scan","Remove","","Run" }
  
  p.Icons=make(map[string]UI.IconItemInterface)
  p.ListFont =  UI.Fonts["notosanscjk15"]
  
  p.MyStack = NewEmuStack()
  
  p.BGwidth = 75
  p.BGheight = 73
  
  return p
}
func (self *FavListPage) GetMapIcons() map[string]UI.IconItemInterface {
  return self.Icons
}

func (self *FavListPage) GetEmulatorConfig() *ActionConfig {
  return self.EmulatorConfig
}

func (self *FavListPage) GeneratePathList(path string) ([]map[string]string,error) {
  if UI.IsDirectory(path) == false {
    return nil,errors.New("Path is not a folder")
  }
  dirmap := make(map[string]string)
  var ret []map[string]string
  
  file_paths,err := filepath.Glob(path+"/*")//sorted
  if err != nil {
    fmt.Println(err)
    return ret,err
  }
  
  for _,v := range file_paths {
    if UI.IsDirectory(v) && self.EmulatorConfig.FILETYPE == "dir" { // like DOSBOX
      gameshell_bat := self.EmulatorConfig.EXT[0]
      if UI.GetGid(v) != FavGID { //only favs
        continue
      }
      
      if UI.FileExists( filepath.Join(v,gameshell_bat))  == true {
        dirmap["gamedir"] = v
        ret = append(ret,dirmap)
      }
    }
    
    if UI.IsAFile(v) && self.EmulatorConfig.FILETYPE == "file" {
      if UI.GetGid(v) != FavGID { //only favs
        continue
      }
      
      bname := filepath.Base(v)
      if len(bname) > 1 {        
        pieces := strings.Split(bname,".")
        if len(pieces) > 1 {
          pieces_ext := strings.ToLower( pieces[len(pieces)-1])
          for _,u := range self.EmulatorConfig.EXT {
            if pieces_ext == u {
              dirmap["file"] = v
              ret = append(ret,dirmap)
              break
            }
          }
        }
      }
    }
  }
  
  return ret,nil
  
}

func (self *FavListPage) SyncList( path string ) {
  
  alist,err := self.GeneratePathList(path) 

  if err != nil {
    fmt.Println(err)
    return
  }
  
  self.MyList = nil 
  
  start_x := 0 
  start_y := 0
  
  hasparent := 0 
  
  if self.MyStack.Length() > 0 {
    hasparent = 1
    
    li := NewEmulatorListItem()
    li.Parent = self
    li.PosX   = start_x
    li.PosY   = start_y
    li.Width  = UI.Width
    li.Fonts["normal"] = self.ListFont
    li.MyType = UI.ICON_TYPES["DIR"]
    li.Init("[..]")
    self.MyList = append(self.MyList,li)
  }
  
  for i,v := range alist {
    li := NewEmulatorListItem()
    li.Parent = self
    li.PosX   = start_x
    li.PosY   = start_y + (i+hasparent)*li.Height
    li.Fonts["normal"] = self.ListFont
    li.MyType = UI.ICON_TYPES["FILE"]
    
    init_val := "NoName"
    
    if val, ok := v["directory"]; ok {
      li.MyType = UI.ICON_TYPES["DIR"]
      init_val = val
    }
    
    if val, ok := v["file"]; ok {
      init_val = val
    }
    
    if val, ok := v["gamedir"]; ok {
      init_val = val
    }
    
    li.Init(init_val)
    
    self.MyList = append(self.MyList,li)
  }
}

func (self *FavListPage) Init() {
  self.PosX = self.Index *self.Screen.Width
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height
  
  self.CanvasHWND = self.Screen.CanvasHWND
  
  ps := UI.NewInfoPageSelector()
  ps.Width  = UI.Width - 12
  ps.PosX = 2
  ps.Parent = self
  
  self.Ps = ps
  self.PsIndex = 0
  
  self.SyncList( self.EmulatorConfig.ROM )
  
  self.MyStack.EmulatorConfig = self.EmulatorConfig
  
  icon_for_list := UI.NewMultiIconItem()
  icon_for_list.ImgSurf = UI.MyIconPool.GetImgSurf("sys")
  icon_for_list.MyType = UI.ICON_TYPES["STAT"]
  icon_for_list.Parent = self
  
  icon_for_list.Adjust(0,0,18,18,0)
        
  self.Icons["sys"] = icon_for_list  
  
  bgpng := UI.NewIconItem()
  bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("star")
  bgpng.MyType = UI.ICON_TYPES["STAT"]
  bgpng.Parent = self
  bgpng.AddLabel("my favourites games",UI.Fonts["varela22"])
  bgpng.SetLabelColor(&color.Color{204,204,204,255}  )
  bgpng.Adjust(0,0,self.BGwidth,self.BGheight,0)

  self.Icons["bg"] = bgpng
  
  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = self.Width - 10
  self.Scroller.PosY = 2
  self.Scroller.Init()
  
  rom_so_confirm_page := NewRomSoConfirmPage()
  rom_so_confirm_page.Screen = self.Screen
  rom_so_confirm_page.Name = "Download Confirm"
  rom_so_confirm_page.Parent = self
  rom_so_confirm_page.Init()

  self.RomSoConfirmDownloadPage = rom_so_confirm_page 
}


func (self *FavListPage) ScrollUp() {
  if len(self.MyList) == 0 {
    return
  }
  
  self.PsIndex -=1
  
  if self.PsIndex < 0 {
    self.PsIndex = 0
  }
  
  cur_li := self.MyList[self.PsIndex]
  x,y := cur_li.Coord()
  _,h := cur_li.Size()
  
  if y < 0 {
    for i,_ := range self.MyList{
      self.MyList[i].NewCoord(x, y + h)
    }
    
    self.Scrolled +=1
  }
}


func (self *FavListPage) ScrollDown(){
  if len(self.MyList) == 0 {
    return
  }
  self.PsIndex +=1
  
  if self.PsIndex >= len(self.MyList) {
    self.PsIndex = len(self.MyList) - 1
  }
  
  cur_li := self.MyList[self.PsIndex]
  x,y := cur_li.Coord()
  _,h := cur_li.Size()
  if y + h > self.Height { 
    for i,_ := range self.MyList{
      self.MyList[i].NewCoord(x,y-h)
    }
    self.Scrolled -=1    
  }

}

func (self *FavListPage) SyncScroll() {

  if self.Scrolled == 0 {
    return
  }
  
  if self.PsIndex < len(self.MyList) {
    cur_li := self.MyList[self.PsIndex]
    x,y := cur_li.Coord()
    _,h := cur_li.Size()
    
    if self.Scrolled > 0 {
      if y < 0 {
        for i,_ := range self.MyList{
          _,h = self.MyList[i].Size()
          self.MyList[i].NewCoord(x, y + self.Scrolled*h)
        }
      }
    }else if self.Scrolled < 0 {
      if y  + h > self.Height {
        for i,_ := range self.MyList {
          _,h = self.MyList[i].Size()
          self.MyList[i].NewCoord(x,y +  self.Scrolled*h)
        }
      }
    }
  
  }
}


func (self *FavListPage) Click() {

  if len(self.MyList) == 0 {
    return
  }
  
  
  if self.PsIndex > len(self.MyList) - 1 {
    return
  }
  
  
  cur_li := self.MyList[self.PsIndex]
  
  if cur_li.(*EmulatorListItem).MyType == UI.ICON_TYPES["DIR"] {
    if cur_li.(*EmulatorListItem).Path ==  "[..]" {
      self.MyStack.Pop()
      self.SyncList(self.MyStack.Last())
      self.PsIndex = 0
    }else{
      self.MyStack.Push(self.MyList[self.PsIndex].(*EmulatorListItem).Path)
      self.SyncList(self.MyStack.Last())
      self.PsIndex = 0
    }
  }
  
  if cur_li.(*EmulatorListItem).MyType == UI.ICON_TYPES["FILE"] {
    self.Screen.MsgBox.SetText("Launching")
    self.Screen.MsgBox.Draw()
    self.Screen.SwapAndShow()
    
    path := ""
    if self.EmulatorConfig.FILETYPE == "dir" {
      path = filepath.Join(cur_li.(*EmulatorListItem).Path,self.EmulatorConfig.EXT[0])
    }else{
      path  = cur_li.(*EmulatorListItem).Path
    }
    
    fmt.Println("Run ",path)
    
    escaped_path := UI.CmdClean(path)
    
    if self.EmulatorConfig.FILETYPE == "dir" {
      escaped_path = UI.CmdClean(path)
    }
    
    custom_config := ""
    
    if self.EmulatorConfig.RETRO_CONFIG != "" && len(self.EmulatorConfig.RETRO_CONFIG) > 5 {
      custom_config = " -c " + self.EmulatorConfig.RETRO_CONFIG
    }
    
    partsofpath := []string{self.EmulatorConfig.LAUNCHER,self.EmulatorConfig.ROM_SO,custom_config,escaped_path}
    
    cmdpath := strings.Join( partsofpath," ")
    
    if self.EmulatorConfig.ROM_SO =="" { //empty means No needs for rom so 
      event.Post(UI.RUNEVT,cmdpath)
    }else{
      
      if UI.FileExists(self.EmulatorConfig.ROM_SO) == true {
        event.Post(UI.RUNEVT,cmdpath)
      } else {
        self.Screen.PushCurPage()
        self.Screen.SetCurPage( self.RomSoConfirmDownloadPage)
        self.Screen.Draw()
        self.Screen.SwapAndShow()
      }
    }
    
    return
    
  }
  
  self.Screen.Draw()
  self.Screen.SwapAndShow() 
}

func (self *FavListPage) ReScan() {
  if self.MyStack.Length() == 0 {
    self.SyncList(self.EmulatorConfig.ROM)
  }else{
    self.SyncList(self.MyStack.Last())
  }
  
  
  idx := self.PsIndex
  
  if idx > len(self.MyList) - 1 {
    idx = len(self.MyList)
    if idx > 0 {
      idx -= 1
    }else if idx == 0 {
      //nothing in MyList
    }
  }
  
  self.PsIndex = idx //sync PsIndex
  
  self.SyncScroll()
}


func (self *FavListPage) OnReturnBackCb() {
  self.ReScan()
  self.Screen.Draw()
  self.Screen.SwapAndShow()
}

func (self *FavListPage) OnLoadCb() {
  self.ReScan()
  self.Screen.Draw()
  self.Screen.SwapAndShow()
}

func (self *FavListPage) KeyDown(ev *event.Event) {

  if ev.Data["Key"] == UI.CurKeys["Menu"] || ev.Data["Key"] == UI.CurKeys["Left"] {
    self.ReturnToUpLevelPage()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
    
  if ev.Data["Key"] == UI.CurKeys["Up"]{
    self.ScrollUp()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  if ev.Data["Key"] == UI.CurKeys["Down"] {
    self.ScrollDown()
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  }
  
  if ev.Data["Key"] == UI.CurKeys["Enter"] {
    self.Click()
  }
    
  if ev.Data["Key"] == UI.CurKeys["X"] { //Scan current
    self.ReScan()
    self.Screen.Draw()
    self.Screen.SwapAndShow()        
  }
  
  if ev.Data["Key"] == UI.CurKeys["Y"] {// del
    if len(self.MyList) == 0 {
      return
    }
    
    cur_li := self.MyList[self.PsIndex] 
    if cur_li.(*EmulatorListItem).IsFile() {
      uid := UI.GetUid(cur_li.(*EmulatorListItem).Path)
      os.Chown(cur_li.(*EmulatorListItem).Path,uid ,uid)
      self.Screen.MsgBox.SetText("Deleting")
      self.Screen.MsgBox.Draw()
      self.Screen.SwapAndShow()
      time.BlockDelay(600)
      self.ReScan()
      self.Screen.Draw()
      self.Screen.SwapAndShow()
    }
  }
}

func (self *FavListPage) Draw() {
  self.ClearCanvas()
  
  if len(self.MyList) == 0 {
    self.Icons["bg"].NewCoord(self.Width/2,self.Height/2)
    self.Icons["bg"].Draw()
  }else{
    _,h := self.Ps.Size()
    if len(self.MyList) * HierListItemDefaultHeight > self.Height {
      
      self.Ps.NewSize(self.Width - 10, h)
      self.Ps.Draw()
      
      
      for _,v := range self.MyList {
        _, y := v.Coord()
        if y > self.Height + self.Height/2 {
          break
        }
        
        if y < 0 {
          continue
        }
        
        v.Draw()
      }
      
      self.Scroller.UpdateSize( len(self.MyList)*HierListItemDefaultHeight, self.PsIndex*HierListItemDefaultHeight)
      self.Scroller.Draw()
      
      
      
    }else {
      self.Ps.NewSize(self.Width,h)
      self.Ps.Draw()
      for _,v := range self.MyList {
        v.Draw()
      }
    }
  }
}


