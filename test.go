package main

import (
	"os"
	"fmt"
	
	"github.com/veandco/go-sdl2/sdl"

	//"github.com/cuu/gogame"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/display"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/draw"
	"github.com/cuu/gogame/image"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/time"	
)

const (
    RUNEVT=1
)
func run() int {

	width := 320
	height := 240
	
	display.Init()
	
	screen := display.SetMode(int32(width),int32(height),0,32)
	
	surface.Fill(screen, &color.Color{255,255,255,255} ) 

	rect1 := rect.Rect(0,10, 12, 10)
	
	//surface.FillRect(screen,&rect, 0xffff0000)
	rect1.X = 12
	draw.Rect(screen,&color.Color{129,235,234,255},&rect1,1)

	fmt.Println(screen.Pitch)
	fmt.Println( screen.BytesPerPixel() )

	img_surf := image.Load("skin/default/sysgo/gameshell/icons/roundcorners.png")

	
	fmt.Println("WxH: ", img_surf.W,img_surf.H)

	portion := rect.Rect(0,0,10,10)
	surface.Blit(screen,img_surf, draw.MidRect(5,5,10,10,width,height), &portion)
	portion.Y = 10
	surface.Blit(screen,img_surf, draw.MidRect(315,5,10,10,width,height), &portion)	
	portion.Y = 20
	surface.Blit(screen,img_surf, draw.MidRect(5,235,10,10,width,height), &portion)
	portion.Y = 30
	surface.Blit(screen,img_surf, draw.MidRect(315,235,10,10,width,height), &portion)	

	
	/*
	for i:=1; i<319;i++ {
		draw.Point(screen, color.Color{255,44,255,0}, i,20)
	}
  */
	
	draw.Line(screen,&color.Color{255,44,255,255}, 0,100, 320,100,3)
	draw.Line(screen,&color.Color{255,44,255,255}, 10, 0, 10,250,4)

	rect2 := rect.Rect(3,120,200,30)
	draw.AARoundRect(screen,&rect2,&color.Color{0,213,222,255},10,0, &color.Color{0,213,222,255})

	rect3 := rect.Rect(300,12,7,200)
	draw.AARoundRect(screen,&rect3,&color.Color{0,213,222,255},3,0, &color.Color{0,213,222,255})
	
	font.Init()
	
	font_path := "skin/default/truetype/NotoSansCJK-Regular.ttf"
	
	notocjk15 := font.Font(font_path,15)

	fmt.Println( font.LineSize( notocjk15 ))

 	my_text := font.Render(notocjk15,"Test ㆑ ㆒ ㆓ ㆔ ㆕ ㆖ 豈 更 車 賈 滑 串 句 龜 龜 契 金 ",true, &color.Color{234,123,12,255},nil)

	surface.Blit(screen,my_text,draw.MidRect(width/2,100,surface.GetWidth(my_text),surface.GetHeight(my_text),width,height),nil)

 	my_text2 := font.Render(notocjk15,"Test ㆑ ㆒ ㆓ ㆔ ㆕ ㆖ 豈 更 車 賈 滑 串 句 龜 龜 契 金 ",true, &color.Color{234,123,12,255},&color.Color{0,0,111,255})	
	surface.Blit(screen,my_text2,draw.MidRect(width/2,100+font.LineSize(notocjk15),surface.GetWidth(my_text),surface.GetHeight(my_text),width,height),nil)	
	
	display.Flip()

	event.AddCustomEvent(RUNEVT)
    
	running := true
	for running {
		ev := event.Wait()
		if ev.Type == event.QUIT {
			running = false
			break
		}
		if ev.Type == event.USEREVENT {
			
			fmt.Println(ev.Data["Msg"])
		}
		if ev.Type == event.KEYDOWN {
			fmt.Println(ev)
			if ev.Data["Key"] == "Q" {
				return 0
			}
			if ev.Data["Key"] == "Escape" {
				return 0
			}
			if ev.Data["Key"] == "T" {
				time.Delay(1000)
			}
			if ev.Data["Key"] == "P" {				
				event.Post(RUNEVT,"GODEBUG=cgocheck=0 sucks") // just id and string, simplify the stuff
			}
		}
	}

	return 0
}

func main() {
	var exitcode int

	os.Setenv("SDL_VIDEO_CENTERED","1")
	os.Setenv("GODEBUG", "cgocheck=0")
	
	sdl.Main(func() {
		exitcode = run()
	})

	os.Exit(exitcode)
}
