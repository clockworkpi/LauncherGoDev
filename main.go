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
  "os/exec"
  //"encoding/json"
	gotime "time"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/cuu/gogame"
	"github.com/cuu/gogame/display"
	"github.com/cuu/gogame/event"
//	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/font"
	//"github.com/cuu/gogame/time"
	
  "github.com/clockworkpi/LauncherGoDev/sysgo"
  
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"

  
)

var (
  flash_led1_counter  = 0
  last_brt = -1
  passout_time_stage = 0
  led1_proc_file = "/proc/driver/led1"
  
  everytime_keydown = gotime.Now()
  
  sound_patch *UI.SoundPatch
  
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
      
      if main_screen.TitleBar.InLowBackLight >= 0 {
        everytime_keydown = cur_time
        continue
      }
      
      if UI.FileExists(sysgo.BackLight) {
        lines,err := UI.ReadLines(sysgo.BackLight) 
        
        if err == nil {
          brt,err2 := strconv.Atoi(strings.Trim(lines[0],"\r\n "))
          if err2 == nil {
            if brt > 0 {
              if last_brt < 0 {
                last_brt = brt
              }
              d := []byte(fmt.Sprintf("%d",1)) // lowest backlight
              ioutil.WriteFile(sysgo.BackLight,d,0644)
            }
          }
        }
      
      
        main_screen.TitleBar.InLowBackLight = 0
        if time2 != 0 {
          passout_time_stage = 1 // next 
        }
      }
      everytime_keydown = cur_time
    }else if elapsed > gotime.Duration(time2) *gotime.Second && passout_time_stage == 1 {
      fmt.Println("timeout, close screen ", elapsed)
      
      if main_screen.Closed == true {
        everytime_keydown = cur_time
        continue
      }
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
      
      if UI.FileExists(sysgo.BackLight) { //hdmi does not have BackLight dev node
        d := []byte(fmt.Sprintf("%d",last_brt))
        ioutil.WriteFile(sysgo.BackLight,d,0644)
        
        main_screen.CounterScreen.Draw()
        main_screen.CounterScreen.SwapAndShow()
        main_screen.CounterScreen.StartCounter()        
        
      }
      
      main_screen.TitleBar.InLowBackLight = 0
      
      passout_time_stage = 4
      
    }
    
    gotime.Sleep(gotime.Duration(UI.DT) * gotime.Millisecond)
  }
}

//If not under awesomeWM, AutoRedraw improves the experience of gsnotify 
//awesomeWM can hold individual window's content from being polluted without redrawing
func AutoRedraw(main_screen *UI.MainScreen) {

	for {
    if main_screen.TitleBar.InLowBackLight < 0 {
      UI.SwapAndShow()
    }
    gotime.Sleep(650 * gotime.Millisecond)
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
    
    if pwl != ""{
      sysgo.CurPowerLevel = pwl
      if pwl == "supersaving" {
        UI.System("sudo iw wlan0 set power_save on >/dev/null")
      }else{
        UI.System("sudo iw wlan0 set power_save off >/dev/null")
      }
    }else {
      UI.System("sudo iw wlan0 set power_save off >/dev/null")
    }
  }
  
}

func run() int {	
	display.Init()
	font.Init()
	screen := display.SetMode(int32(UI.Width),int32(UI.Height),0,32)
    
	UI.Init()
  
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

	ReadTheDirIntoPages(main_screen,"Menu",0,nil)
  ReadTheDirIntoPages(main_screen,"/home/cpi/apps/Menu",1,main_screen.Pages[len(main_screen.Pages)-1])
  ReunionPagesIcons(main_screen)
  
	main_screen.FartherPages()
  
  sound_patch = UI.NewSoundPatch()
  sound_patch.Parent = main_screen
  sound_patch.Init()
  
	main_screen.Draw()
	main_screen.SwapAndShow()

	UI.SwapAndShow()
	
	//fmt.Println(main_screen)
  event.AllocEvents(5)
  event.AddCustomEvent(UI.RUNEVT)
  event.AddCustomEvent(UI.RUNSH)
  event.AddCustomEvent(UI.RUNSYS)
  event.AddCustomEvent(UI.RESTARTUI)
  event.AddCustomEvent(UI.POWEROPT)

  go FlashLed1(main_screen)
  go InspectionTeam(main_screen)
  go main_screen.TitleBar.RoundRobinCheck()
  //go AutoRedraw(main_screen)
  
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
          exec_app_cmd +="; sync & cd "+UI.GetExePath()+"; "+os.Args[0]+";"
          fmt.Println(exec_app_cmd)
          cmd := exec.Command("/bin/sh","-c",exec_app_cmd)
          err := cmd.Start()
          if err != nil {
            fmt.Println(err)
          }
          err = cmd.Process.Release()
          if err != nil {
            fmt.Println(err)
          }
          os.Exit(0)
          
        case UI.RUNSYS:
          main_screen.OnExitCb()      
          gogame.Quit()   
          exec_app_cmd := ev.Data["Msg"]
          cmd := exec.Command("/bin/sh","-c",exec_app_cmd)
          err := cmd.Start()
          if err != nil {
            fmt.Println(err)
          }
          err = cmd.Process.Release()
          if err != nil {
            fmt.Println(err)
          }
          os.Exit(0)

        case UI.RUNSH:
          main_screen.OnExitCb()      
          gogame.Quit()          
          
          fmt.Println("RUNSH")
          exec_app_cmd := ev.Data["Msg"]+";"
          fmt.Println(exec_app_cmd)
          cmd := exec.Command("/bin/sh","-c",exec_app_cmd)
          err := cmd.Start()
          if err != nil {
            fmt.Println(err)
          }
          err = cmd.Process.Release()
          if err != nil {
            fmt.Println(err)
          }
          os.Exit(0)          
        
        case UI.RESTARTUI:
          main_screen.OnExitCb()      
          gogame.Quit()
          exec_app_cmd :=" sync & cd "+UI.GetExePath()+"; "+os.Args[0]+";"
          fmt.Println(exec_app_cmd)
          cmd := exec.Command("/bin/sh","-c",exec_app_cmd)
          err := cmd.Start()
          if err != nil {
            fmt.Println(err)
          }
          err = cmd.Process.Release()
          if err != nil {
            fmt.Println(err)
          }
          os.Exit(0)
          
        case UI.POWEROPT:
          everytime_keydown = gotime.Now()

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
			}
      
      
      if ev.Data["Key"] == "Keypad +" {
        if main_screen.CurPage().GetName() != "Sound volume" {
          main_screen.Draw()
          sound_patch.VolumeUp()
          sound_patch.Draw()
          main_screen.SwapAndShow()
        }
        continue
      }
      
      if ev.Data["Key"] == "Keypad -" {
        if main_screen.CurPage().GetName() != "Sound volume" {
          main_screen.Draw()
          sound_patch.VolumeDown()
          sound_patch.Draw()
          main_screen.SwapAndShow()       
        
        }
        continue
      }
                  
			main_screen.KeyDown(ev)
			
		}
	}

	return 0
}

func main() {
	
	var exitcode int
  
  runtime.GOMAXPROCS(1)
  
	os.Setenv("SDL_VIDEO_CENTERED","1")
	
	sdl.Main(func() {
		exitcode = run()
	})

	os.Exit(exitcode)
}
