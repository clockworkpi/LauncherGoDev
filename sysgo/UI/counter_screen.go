package UI

import (
  
  gotime "time"
  
)
type CounterScreen struct {
  FullScreen
  

  CounterFont  *ttf.Font
  TextFont1    *ttf.Font
  TextFont2    *ttf.Font
  
  TopLabel    LabelInterface
  BottomLabel  LabelInterface
  NumberLabel  LabelInterface
  BGColor     *color.Color
  FGColor     *color.Color
  
  Counting    bool
  
  Number      int // 10
  
  inter_counter int //
  
  TheTicker *gotime.Ticker
  TickerStoped chan bool
  
}

func NewCounterScreen() *CounterScreen {
  p := &CounterScreen{}
  p.Number = 10
  p.CounterFont = Fonts["varela120"]
  p.TextFont1 = Fonts["varela15"]
  p.TextFont2 = Fonts["varela12"]
  
  p.BGColor = &color.Color{0,0,0,255}
  p.FGColor = &color.Color{255,255,255,255}
  
  return p
}

func (self *CounterScreen ) Interval() {

 for {
		select {
		case <-self.TheTicker.C:
      self.inter_counter += 1
      
      if self.Number == 0 {
        self.Counting = false
        self.TheTicker.Stop()
        fmt.Println("do the real shutdown")
        
        if sysgo.CurKeySet != "PC" {
          cmdpath := "feh --bg-center sysgo/gameshell/wallpaper/seeyou.png;"
          cmdpath = cmdpath + "sleep 3;"
          cmdpath = cmdpath + "sudo halt -p"
          event.Post(RUNEVT,cmdpath)
          
        }
        
        break
      }
      
      if self.inter_counter >= 2 {
        self.Number -= 1
        if self.Number < 0 {
          self.Number = 0
        }
        
        fmt.Println("sub Number ", self.Number)
        self.inter_counter = 0
        
        self.Draw()
        self.SwapAndShow()
      
      }
    case <- self.TickerStoped:
      break
    }
  }

}


func (self *CounterScreen) StartCounter() {
  if self.Counting == true {
    return
  }
  
  self.Number = 10
  self.inter_counter = 0
  
  self.Counting = true
  
  self.TheTicker.Start()
  
  go self.Interval()

}


func (self *CounterScreen) StopCounter() {
  if self.Counting == false {
    return
  }
  
  self.Counting = false
  self.Number = 0
  self.inter_counter = 0
  
  self.TheTicker.Stop()
  self.TickerStoped <- true
  
}

func (self *CounterScreen) Init() {
  
  self.CanvasHWND = surface.Surface(self.Width,self.Height)
  
  self.TopLabel = NewLabel()
  self.TopLabel.SetCanvasHWND( self.CanvasHWND)
  self.TopLabel.Init("System shutdown in", self.TextFont1,self.FGColor)
  
  self.BottomLabel = NewLabel()
  self.BottomLabel.SetCanvasHWND(self.CanvasHWND)
  self.BottomLabel.Init("Press any key to stop countdown",self.TextFont2,self.FGColor)
  
  
  self.NumberLabel  = NewLabel()
  self.NumberLabel.SetCanvasHWND(self.CanvasHWND)
  number_str := fmt.Sprintf("%d",self.Number)
  self.NumberLabel.Init(number_str,self.CounterFont,self.FGColor)
  
  self.TheTicker = gotime.NewTicker(500 * gotime.Millisecond)
  self.TickerStoped = make(chan bool,1)
  
}

func (self *CounterScreen) Draw() {
  surface.Fill(self.CanvasHWND, self.FGColor)
  
  self.TopLabel.NewCoord(Width/2,15)
  self.TopLabel.DrawCenter(false)
  
  self.BottomLabel.NewCoord(Width/2, Height-15)
  self.BottomLabel.DrawCenter(false)

  self.NumberLabel.NewCoord(Width/2,Height/2)
  number_str := fmt.Sprintf("%d",self.Number)
  self.NumberLabel.SetText(number_str)
  self.NumberLabel.DrawCenter(false)  

}

