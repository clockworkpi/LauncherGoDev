package surface

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"../color"
)


func Fill(surface *sdl.Surface, col color.Color) {

	rect := sdl.Rect{0,0,0,0}

	rect.W = surface.W
	rect.H = surface.H

	FillRect(surface, &rect, uint32(col.ToHex()))
}

func FillRect(surface *sdl.Surface,rect *sdl.Rect, color uint32) {
	
	sdl.Do(func() {
		surface.FillRect(rect,color)
	})
	
	return
}

// Create a New Surface
func Surface(w,h int) *sdl.Surface {
	//flags=0, depth=0, masks=None
	Rmask := 0x000000ff
	Gmask := 0x0000ff00
	Bmask := 0x00ff0000
	Amask := 0xff000000
	
	flags := 0
	depth := 32
	
	surf,err := sdl.CreateRGBSurface(uint32(flags),int32(w),int32(h), int32(depth), uint32(Rmask), uint32(Gmask), uint32(Bmask), uint32(Amask))
	if err != nil {
		panic( fmt.Sprintf("sdl.CreateRGBSurface failed %s",sdl.GetError()))
	}

	return surf
}
