package draw

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"../color"
)

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
	
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
	
func abs(n int) int {
	return int(math.Abs(float64(n)))
}

func pixel(surf *sdl.Surface, c color.Color, x,y int) int {
	pixels := surf.Pixels()
	bytes_per_pixel := surf.BytesPerPixel()
	
	addr := y * int(surf.Pitch) + x*bytes_per_pixel // 1 2 3 4

	color_bytes := c.ToBytes()

	if x < int(surf.ClipRect.X) || x >= int(surf.ClipRect.X + surf.ClipRect.W) ||
		y < int(surf.ClipRect.Y) || y >= int(surf.ClipRect.Y + surf.ClipRect.H) {
		return 0
	}
	
	if bytes_per_pixel == 1 {
		pixels[addr] = color_bytes[0]
	}

	if bytes_per_pixel == 2 {
		for i :=0; i < bytes_per_pixel; i++ {
			pixels[addr+i] = color_bytes[i]
		}	
	}

	if bytes_per_pixel == 3 {
		for i :=0; i < bytes_per_pixel; i++ {
			pixels[addr+i] = color_bytes[i]
		}	
	}

	if bytes_per_pixel == 4 {
		for i :=0; i < bytes_per_pixel; i++ {
			pixels[addr+i] = color_bytes[i]
		}
	}

	return 1
}
