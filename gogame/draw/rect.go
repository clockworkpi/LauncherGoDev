package draw

import (
//	"math"
	"github.com/veandco/go-sdl2/sdl"
	"../color"
	
)

func Rect(surf *sdl.Surface,color color.Color, _rect *sdl.Rect, border_width int) {
	l := int(_rect.X)
	r := int(_rect.X + _rect.W - 1)
	t := int(_rect.Y)
	b := int(_rect.Y + _rect.H - 1)
	
	points := [][]int{ []int{l,t}, []int{r,t}, []int{r,b},[]int{l,b} }
	Polygon(surf, color, points, border_width)
	
}
