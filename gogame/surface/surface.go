package surface

import (
	"github.com/veandco/go-sdl2/sdl"
)


func FillRect(surface *sdl.Surface,rect *sdl.Rect, color uint32) {
	
	sdl.Do(func() {
		surface.FillRect(rect,color)
	})
	
	return
}

