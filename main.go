package main

import (

	"os"
	"fmt"
	
	"github.com/veandco/go-sdl2/sdl"
	
	"github.com/cuu/gogame/display"
	"github.com/cuu/gogame/event"
//	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/time"
	
	"github.com/cuu/LauncherGoDev/sysgo/UI"
)


func run() int {	
	display.Init()
	font.Init()
	screen := display.SetMode(int32(UI.Width),int32(UI.Height),0,32)
    
	UI.Init()
	UI.MyIconPool.Init()

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
	
	fmt.Println(main_screen)
    
	event.AddCustomEvent(UI.RUNEVT)

	running := true
	for running {
		ev := event.Wait()
		if ev.Type == event.QUIT {
			running = false
			break
		}
		if ev.Type == event.USEREVENT {
			
			fmt.Println("UserEvent: ",ev.Data["Msg"])
		}
		if ev.Type == event.KEYDOWN {
			if ev.Data["Key"] == "Q" {
				main_screen.OnExitCb()
				return 0
			}else if ev.Data["Key"] == "D" {
				time.Delay(1000)
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
