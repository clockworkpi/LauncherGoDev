package ButtonsLayout

import (
  "fmt"
  "io/ioutil"
  "os/exec"
  //"path/filepath"
  "strings"
  
  "github.com/veandco/go-sdl2/ttf"


  //"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	//"github.com/cuu/gogame/time"
  
  //"github.com/clockworkpi/LauncherGoDev/sysgo"

  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"

)

type UpdateConfirmPage struct {
  UI.ConfirmPage
  
  RetroArchConf string
  
  LayoutMode string
  
}

func NewUpdateConfirmPage() *UpdateConfirmPage {

  p := &UpdateConfirmPage{}
  
  p.ListFont = UI.MyLangManager.TrFont("veramono20")
  
  p.FootMsg = [5]string{"Nav","","","Cancel","Yes"}
    
  p.ConfirmText = "Apply to RetroArch?"
  p.RetroArchConf = "/home/cpi/.config/retroarch/retroarch.cfg"
  p.LayoutMode = "Unknown"
  
  return p
}

func (self *UpdateConfirmPage) ModifyRetroArchConf( keys []string ) string {

  if UI.FileExists(self.RetroArchConf) {
    
    confarr,err := UI.ReadLines(self.RetroArchConf)
    var bka  = false 
    var bkb  = false
    var bkx  = false
    var bky  = false
    
    if err == nil {
      for i,ln := range confarr {
        parts := strings.Split(ln,"=")
        if len(parts) < 1 {
          fmt.Println("retroarch.cfg cannot parse.")
          return "retroarch.cfg cannot parse."
        }
        lnk := strings.Trim(parts[0],"\r\n ")
        if lnk == "input_player1_a" {
          confarr[i] = "input_player1_a = \"" + keys[0] + "\"\n"
          bka=true
        }
        if lnk == "input_player1_b" {
          confarr[i] = "input_player1_b = \"" + keys[1] + "\"\n"
          bkb = true
        }
        if lnk == "input_player1_x" {
          confarr[i] = "input_player1_x = \"" + keys[2] + "\"\n"
          bkx = true
        }
        
        if lnk == "input_player1_y" {
          confarr[i] = "input_player1_y = \"" + keys[3] + "\"\n"
          bky = true
        }
      }
      
      if bka == false || bkb == false || bkx == false || bky == false {
        fmt.Println("retroarch.cfg validation error.")
        return "retroarch.cfg validation error."
      }
    }
    
    err = UI.WriteLines(confarr,self.RetroArchConf)
    if err != nil {
      fmt.Println(err)
      return "retroarch.cfg cannot write."
    }
  
  }

  fmt.Println( "Completed! Your RA keymap: " + strings.ToUpper(self.LayoutMode)  )
  return "Completed! Your RA keymap: " + strings.ToUpper(self.LayoutMode)
}

func (self *UpdateConfirmPage) finalizeWithDialog(msg string) {
  self.Screen.MsgBox.SetText(msg)
  self.Screen.MsgBox.Draw()
  self.Screen.SwapAndShow()
  return
}

func (self *UpdateConfirmPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  
  if ev.Data["Key"] == UI.CurKeys["B"] {
  
    keymap := []string{"j","k","u","i"}
    
    if self.LayoutMode == "xbox" {
    
      keymap = []string{"j","k","u","i"}
      
    }else if self.LayoutMode == "snes" {
      
      keymap = []string{ "k","j","i","u" }
      
    }else {
      self.finalizeWithDialog("Internal error.")
      return
    }
    
    fmt.Println( "mode: ",self.LayoutMode)
    
    if UI.IsAFile(self.RetroArchConf) == false {
      self.finalizeWithDialog("retroarch.cfg was not found.")
      return
    }
    
    cpCmd := exec.Command("cp", "-rf", self.RetroArchConf,self.RetroArchConf+".blbak")
    err := cpCmd.Run()
    if err != nil {
      fmt.Println(err)
      self.finalizeWithDialog("Cannot create .blbak")
      return
    }
    
    self.finalizeWithDialog(self.ModifyRetroArchConf(keymap))
    return
  }
}

func (self *UpdateConfirmPage) OnReturnBackCb() {
  self.ReturnToUpLevelPage()
  self.Screen.Draw()
	self.Screen.SwapAndShow()
}

func (self *UpdateConfirmPage) Draw() {
  self.ClearCanvas()
  self.DrawBG()
  for _,v := range self.MyList {
    v.Draw()
  }
  
  self.Reset()
}

type ButtonsLayoutPage struct {
  
  UI.Page
  ListFontObj *ttf.Font
  BGwidth int 
  BGheight int
  
  DrawOnce bool
  Scrolled int
  Scroller *UI.ListScroller
  ConfirmPage *UpdateConfirmPage
  
  dialog_index int 
  Icons map[string]UI.IconItemInterface
  
  ConfigFilename string
}


func NewButtonsLayoutPage() *ButtonsLayoutPage {
  p := &ButtonsLayoutPage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]
	
	p.FootMsg = [5]string{"Nav","UpdateRetroArch","","Back","Toggle"} 
  p.Icons = make( map[string]UI.IconItemInterface  )

  p.BGwidth = UI.Width
  p.BGheight = UI.Height - 24 -20
  
  p.ConfigFilename = "sysgo/.buttonslayout"
  
  return p
  
}

func (self *ButtonsLayoutPage) Init() {
  
  if self.Screen != nil {
    if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
      self.HWND = self.Screen.CanvasHWND
      self.CanvasHWND = surface.Surface(self.Screen.Width,self.Screen.Height)
    }
  }
  
  self.PosX = self.Index*self.Screen.Width 
  self.Width = self.Screen.Width
  self.Height = self.Screen.Height
  
    
  DialogBoxs := UI.NewMultiIconItem()
  DialogBoxs.ImgSurf = UI.MyIconPool.GetImgSurf("buttonslayout")
  DialogBoxs.MyType = UI.ICON_TYPES["STAT"]
  DialogBoxs.Parent = self
  DialogBoxs.IconWidth = 300
  DialogBoxs.IconHeight = 150
  DialogBoxs.Adjust(0,0,134,372,0)
  self.Icons["DialogBoxs"] = DialogBoxs  

  self.Scroller = UI.NewListScroller()
  self.Scroller.Parent = self
  self.Scroller.PosX = self.Width - 10
  self.Scroller.PosY = 2
  self.Scroller.Init()
  self.Scroller.SetCanvasHWND(self.HWND)      
  
  
  self.ConfirmPage = NewUpdateConfirmPage()
  self.ConfirmPage.LayoutMode = self.GetButtonsLayoutMode()
  self.ConfirmPage.Screen = self.Screen
  self.ConfirmPage.Name  = "Overwrite RA conf"
  self.ConfirmPage.Init()
  
}

func (self *ButtonsLayoutPage) ScrollUp() {
  dis := 10
  
  if self.PosY < 0 {
    self.PosY += dis
    self.Scrolled += dis
  }
}

func (self *ButtonsLayoutPage) ScrollDown() {
  dis := 10
  
  if UI.Abs(self.Scrolled) < (self.BGheight - self.Height) / 2 + 0 {
    self.PosY -= dis
    self.Scrolled -=dis 
  }
}


func (self *ButtonsLayoutPage) GetButtonsLayoutMode() string {
  lm := "xbox"
 
  lm_bytes,err := ioutil.ReadFile(self.ConfigFilename)
  
  if err == nil {
    
    for _,v := range []string{"xbox","snes"} {
      if v == string(lm_bytes) {
        lm = string(lm_bytes)
        break
      }
    }
  }

  return lm
}

func (self *ButtonsLayoutPage) ToggleMode() {

  
  if self.GetButtonsLayoutMode() == "xbox" {
    d := []byte("snes")
    err := ioutil.WriteFile(self.ConfigFilename,d,0644)
    if err != nil {
      fmt.Println(err)
    }
    
    self.dialog_index = 1
    self.Screen.Draw()
    self.Screen.SwapAndShow()
  
  }else {
    d := []byte("xbox")
    err := ioutil.WriteFile(self.ConfigFilename,d,0644)
    if err != nil {
      fmt.Println(err)
    }
    
    self.dialog_index = 0
    self.Screen.Draw()
    self.Screen.SwapAndShow()  
  
  }
}

func (self *ButtonsLayoutPage) OnLoadCb() {

  self.Scrolled = 0
  self.PosY = 0
  self.DrawOnce = false
  
  
  if self.GetButtonsLayoutMode() == "xbox" {
    self.dialog_index = 0
  }else {
    self.dialog_index = 1
  }
}

func (self *ButtonsLayoutPage) OnReturnBackCb() {

  self.ReturnToUpLevelPage()
  self.Screen.Draw()
  self.Screen.SwapAndShow()

}


func (self *ButtonsLayoutPage) KeyDown(ev *event.Event) {
  
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
  
  if ev.Data["Key"] == UI.CurKeys["B"] {
    self.ToggleMode()
  }
  
  if ev.Data["Key"] == UI.CurKeys["X"] {
    self.ConfirmPage.LayoutMode = self.GetButtonsLayoutMode()
    self.Screen.PushPage(self.ConfirmPage)
    self.Screen.Draw()
		self.Screen.SwapAndShow()
  }
  
}

func (self *ButtonsLayoutPage) Draw() {

  self.ClearCanvas()
  
  self.Icons["DialogBoxs"].NewCoord(0,30)
  
  self.Icons["DialogBoxs"].SetIconIndex(self.dialog_index)
  self.Icons["DialogBoxs"].DrawTopLeft()
  
  if self.HWND != nil {
    surface.Fill(self.HWND, &color.Color{255,255,255,255})
    rect_ := rect.Rect(self.PosX,self.PosY,self.Width,self.Height)
    surface.Blit(self.HWND,self.CanvasHWND,&rect_,nil)
  }
}


