package main

import (

	"os"
	"fmt"
  "log"
  "io/ioutil"
  "strconv"
  "strings"
  "runtime"
  "path/filepath"
  
	gotime "time"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/cuu/gogame"
	"github.com/cuu/gogame/display"
	"github.com/cuu/gogame/event"
//	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/font"
	//"github.com/cuu/gogame/time"
	
  "github.com/cuu/LauncherGoDev/sysgo"
  
	"github.com/cuu/LauncherGoDev/sysgo/UI"
)

var (
  flash_led1_counter  = 0
  last_brt = 0
  passout_time_stage = 0
  led1_proc_file = "/proc/driver/led1"
  
  everytime_keydown = gotime.Now()
  
)
// flash the Led1 on the GS back
func FlashLed1(main_screen *UI.MainScreen) {
  
  for {
    if UI.FileExists(led1_proc_file) {
      if main_screen.Closed == false {
        if flash_led1_counter > 0 {
          d := []byte(fmt.Sprintf("%d",0))
          err := ioutil.WriteFile(led1_proc_file, d, 0644) // turn off led1
          if err != nil {
            fmt.Println(err)
          }
          flash_led1_counter = 0
        }
      
      } else {
        flash_led1_counter +=1
        if flash_led1_counter == 3 {
          d := []byte(fmt.Sprintf("%d",1))
          err := ioutil.WriteFile(led1_proc_file, d, 0644)
          if err != nil {
            fmt.Println(err)
          }
        }
        
        if flash_led1_counter == 5 {
          d := []byte(fmt.Sprintf("%d",0))
          err := ioutil.WriteFile(led1_proc_file, d, 0644)
          if err != nil {
            fmt.Println(err)
          }
        }
        
        if flash_led1_counter == 11 {
          flash_led1_counter = 1
        }
      }
    }
    
    gotime.Sleep(200 * gotime.Millisecond)
  }
}

//happens everytime when KeyDown occurs
func RestoreLastBackLightBrightness(main_screen *UI.MainScreen) bool {
  
  passout_time_stage = 0
  main_screen.TitleBar.InLowBackLight = -1
  main_screen.Closed = false
  
  if last_brt == -1 {
    return true
  }
  
  if UI.FileExists(sysgo.BackLight) {
    lines,err := UI.ReadLines(sysgo.BackLight)
    if err == nil {
      brt,err2 := strconv.Atoi(strings.Trim(lines[0],"\r\n "))
      if err2 == nil {
        if brt < last_brt {
          d := []byte(fmt.Sprintf("%d",last_brt))
          ioutil.WriteFile(sysgo.BackLight,d,0644)
          last_brt = -1
        }
      }
    }else {
      fmt.Println(err)
    }
    
  }else {
    
  }
  
  if UI.FileExists(led1_proc_file) {
    d := []byte(fmt.Sprintf("%d",0))
    err := ioutil.WriteFile(led1_proc_file, d, 0644)
    if err != nil {
      fmt.Println(err)
    }
  }

  //Stop CounterScreen here
  
  if main_screen.CounterScreen.Counting == true {
    main_screen.CounterScreen.StopCounter()
    main_screen.Draw()
    main_screen.SwapAndShow()
    return false
  }
  
  return true
  
}

//power stuff dealer
func InspectionTeam(main_screen *UI.MainScreen) {

  for {
    cur_time := gotime.Now()
    elapsed := cur_time.Sub(everytime_keydown)
    
    time1 := sysgo.PowerLevels[sysgo.CurPowerLevel].Dim
    time2 := sysgo.PowerLevels[sysgo.CurPowerLevel].Close
    time3 := sysgo.PowerLevels[sysgo.CurPowerLevel].PowerOff
    
    if elapsed > gotime.Duration(time1) *gotime.Second && passout_time_stage == 0 {
      fmt.Println("timeout, dim screen ",elapsed)
      
      if UI.FileExists(sysgo.BackLight) {
        lines,err := UI.ReadLines(sysgo.BackLight) 
        
        if err == nil {
          brt,err2 := strconv.Atoi(strings.Trim(lines[0],"\r\n "))
          if err2 == nil {
            if brt > 0 {
              if last_brt < 0 {
                last_brt = brt
              }
              d := []byte(fmt.Sprintf("%d",1))
              ioutil.WriteFile(sysgo.BackLight,d,0644)
            }
          }
        }
      }
      
      main_screen.TitleBar.InLowBackLight = 0
      if time2 != 0 {
        passout_time_stage = 1 // next 
      }
      everytime_keydown = cur_time
    }else if elapsed > gotime.Duration(time2) *gotime.Second && passout_time_stage == 1 {
      fmt.Println("timeout, close screen ", elapsed)
      
      if UI.FileExists(sysgo.BackLight) {
        d := []byte(fmt.Sprintf("%d",0))
        ioutil.WriteFile(sysgo.BackLight,d,0644)
      }      
      
      main_screen.TitleBar.InLowBackLight = 0
      main_screen.Closed = true
      if time3 != 0 {
        passout_time_stage = 2 // next
      }
      
      everytime_keydown = cur_time
    }else if elapsed > gotime.Duration(time3) * gotime.Second && passout_time_stage  == 2{
      
      fmt.Println("Power Off counting down")
      
      main_screen.CounterScreen.Draw()
      main_screen.CounterScreen.SwapAndShow()
      main_screen.CounterScreen.StartCounter()
      
      if UI.FileExists(sysgo.BackLight) {
        d := []byte(fmt.Sprintf("%d",last_brt))
        ioutil.WriteFile(sysgo.BackLight,d,0644)
      }
      
      main_screen.TitleBar.InLowBackLight = 0
      
      passout_time_stage = 4
      
    }
        
    gotime.Sleep(gotime.Duration(UI.DT) * gotime.Millisecond)
  }
}

func PreparationInAdv(){
  
  if strings.Contains(runtime.GOARCH,"arm") == false {
    return
  }
  
  if UI.FileExists("sysgo/.powerlevel") == false {
    UI.System("touch sysgo/.powerlevel")
    UI.System("sudo iw wlan0 set power_save off >/dev/null")
    
  }else{
    b, err := ioutil.ReadFile("sysgo/.powerlevel")
    if err != nil {
        log.Fatal(err)
    }
    
    pwl := strings.Trim(string(b),"\r\n ")
    
    if pwl == "supersaving" {
      UI.System("sudo iw wlan0 set power_save on >/dev/null")
    }else{
      UI.System("sudo iw wlan0 set power_save off >/dev/null")
    }
  }
  
}

func run() int {	
	display.Init()
	font.Init()
	screen := display.SetMode(int32(UI.Width),int32(UI.Height),0,32)
    
	UI.Init()
	UI.MyIconPool.Init()
  
  PreparationInAdv()
  
	main_screen := UI.NewMainScreen()
	main_screen.HWND = screen
	main_screen.Init()
		
	title_bar := UI.NewTitleBar()
	foot_bar := UI.NewFootBar()

	title_bar.Init(main_screen)
	foot_bar.Init(main_screen)
	
	main_screen.TitleBar = title_bar
	main_screen.FootBar  = foot_bar

	main_screen.ReadTheDirIntoPages("Menu",0,nil)
	main_screen.FartherPages()

	main_screen.Draw()
	main_screen.SwapAndShow()

	UI.SwapAndShow()
	
	//fmt.Println(main_screen)
    
	event.AddCustomEvent(UI.RUNEVT)
  
  go FlashLed1(main_screen)
  go InspectionTeam(main_screen)
  
	running := true
	for running {
		ev := event.Wait()
		if ev.Type == event.QUIT {
			running = false
			break
		}
		if ev.Type == event.USEREVENT {

      
			fmt.Println("UserEvent: ",ev.Data["Msg"])
      
      switch ev.Code {
        case UI.RUNEVT:
          main_screen.OnExitCb()      
          gogame.Quit()          
          
          fmt.Println("RUNEVT")
          exec_app_cmd := "cd " + filepath.Dir(ev.Data["Msg"])+";"
          exec_app_cmd += ev.Data["Msg"]
          exec_app_cmd +="; sync & cd "+UI.GetExePath()+"; "+os.Args[0]
          fmt.Println(exec_app_cmd)
          
          
      }
      
      
		}
		if ev.Type == event.KEYDOWN {
      everytime_keydown = gotime.Now()
      if RestoreLastBackLightBrightness(main_screen) == false {
        continue
      }
      
			if ev.Data["Key"] == "Q" {
				main_screen.OnExitCb()
				return 0
			}else if ev.Data["Key"] == "P" {				
				event.Post(UI.RUNEVT,"GODEBUG=cgocheck=0 sucks") // just id and string, simplify the stuff
				
			}else {
				main_screen.KeyDown(ev)
			}
		}
	}

	return 0
}

func main() {
	
	var exitcode int

	os.Setenv("SDL_VIDEO_CENTERED","1")
	
	sdl.Main(func() {
		exitcode = run()
	})

	os.Exit(exitcode)
}
