package rect

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Rect(top ,left, width,height int) sdl.Rect {
	return sdl.Rect{int32(top),int32(left),int32(width),int32(height)}
}
