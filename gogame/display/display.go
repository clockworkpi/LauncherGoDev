package display

import (
	"github.com/veandco/go-sdl2/sdl"
	"../../gogame"
)

var Inited =  false
var window *sdl.Window

func AssertInited() {
	if Inited == false {
		panic("run gogame.DisplayInit first")
	}
}

func Init() bool {
	sdl.Do(func() {
		
		if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
			panic(err)
		}
	
		Inited = true
	})
	
  return Inited 
}


func SetMode(w,h,flags,depth int32) *sdl.Surface {
	var err error
	var surface *sdl.Surface
	AssertInited()
	
	sdl.Do(func() {
		window, err = sdl.CreateWindow("gogame", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			w, h, uint32( gogame.SHOWN | flags))
	
		if err != nil {
			panic(err)
		}

		surface,err = window.GetSurface()
		if err != nil {
			panic(err)
		}
	})

	return surface
}

func Flip() {
	sdl.Do(func() {

		if window != nil {
			window.UpdateSurface()
		}
	})
}
		


