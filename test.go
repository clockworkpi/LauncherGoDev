package main

import (
	"os"
	"fmt"
	
	"github.com/veandco/go-sdl2/sdl"

	//"./gogame"
	"./gogame/color"
	"./gogame/display"
	"./gogame/surface"
	"./gogame/event"
	"./gogame/rect"
	"./gogame/draw"
	
)

func run() int {
	
	display.Init()
	
	screen := display.SetMode(320,240,0,32)
	
	surface.FillRect(screen,nil, 0)

	rect := rect.Rect(0, 0, 200, 200)
	
	//surface.FillRect(screen,&rect, 0xffff0000)

	rect.X = 12
	
	draw.Rect(screen,color.Color{129,235,234,0},&rect,1)

	fmt.Println(screen.Pitch)
	fmt.Println( screen.BytesPerPixel() )

	/*
	for i:=1; i<319;i++ {
		draw.Point(screen, color.Color{255,44,255,0}, i,20)
	}
  */
	
	draw.Line(screen,color.Color{255,44,255,0}, 0,100, 320,100,3)

	draw.Line(screen,color.Color{255,44,255,0}, 10, 0, 10,250,4)
	display.Flip()

	
	running := true
	for running {
		ev := event.Wait()
		if ev.Type == event.QUIT {
			running = false
			break
		}

		if ev.Type == event.KEYDOWN {
			fmt.Println(ev)
			if ev.Data["Key"] == "Q" {
				return 0
			}
		}
	}

	return 0
}

func main() {
	var exitcode int

	sdl.Main(func() {
		exitcode = run()
	})

	os.Exit(exitcode)
}
