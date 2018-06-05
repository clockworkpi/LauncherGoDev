package rect

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Rect(top ,left, width,height int32) sdl.Rect {
	return sdl.Rect{top,left,width,height}
}
