package draw

import (
	"fmt"
	"math"
	"github.com/veandco/go-sdl2/sdl"
	"../color"
	"../rect"
)

const (
	LEFT_EDGE=0x1
	RIGHT_EDGE=0x2
	BOTTOM_EDGE=0x4
	TOP_EDGE=0x8
)

func Line(surf *sdl.Surface, col color.Color,x1,y1,x2,y2 ,width int) sdl.Rect {
	pts := make([]int,4)
	pts[0] = x1
	pts[1] = y1
	pts[2] = x2
	pts[3] = y2

	if width < 1 {
		return rect.Rect(x1,y1,0,0)
	}

	err := surf.Lock()
	if err != nil {
		return rect.Rect(0,0,0,0)
	}
	anydraw := clip_and_draw_line_width(surf,&surf.ClipRect, col, width,pts)
	surf.Unlock()
	if anydraw == 0 {
		return rect.Rect(x1,y1,0,0)
	}
	rleft := 0
	rtop := 0

	if x1 < x2 {
		rleft = x1
	}else {
		rleft = x2
	}

	if y1 < y2 {
		rtop = y1
	}else {
		rtop = y2
	}

	dx := abs(x1-x2)
	dy := abs(y1-y2)

	rwidth := 0
	rheight := 0
	if dx > dy {
		rwidth = dx +1
		rheight = dy + width
	}else {
		rwidth = dx + width
		rheight = dy + 1
	}

	return rect.Rect(rleft,rtop,rwidth,rheight)	
}

func Rect(surf *sdl.Surface,color color.Color, _rect *sdl.Rect, border_width uint32) {
	color_hex := color.ToHex()
	fmt.Printf("%x\n",color_hex)

}

func clip_and_draw_line(surf *sdl.Surface, rect *sdl.Rect, col color.Color, pts []int)  int {
	
	if clipline(pts, int(rect.X),int(rect.Y),int(rect.X+ rect.W-1), int(rect.Y+rect.H-1) ) == 0 {
		return 0
	}

	if pts[1] == pts[3] {
		drawhorzline(surf, col, pts[0],pts[1],pts[2])
	}else if pts[0] == pts[2] {
		drawvertline(surf,col, pts[0],pts[1],pts[3])
	}else {
		drawline(surf, col, pts[0],pts[1],pts[2],pts[3])
	}

	return 1
}


func clip_and_draw_line_width(surf *sdl.Surface,rect *sdl.Rect,col color.Color,  width int, pts []int) int {
	loop := 0
	xinc :=0
	yinc :=0
	newpts :=make([]int,4)
  range_ := make([]int,4)
	anydraw := 0
	if abs(pts[0]-pts[2]) > abs(pts[1]-pts[3]) {
		yinc = 1
	}else{
		xinc = 1
	}
	copy(newpts, pts)
	if clip_and_draw_line(surf,rect,col, newpts) > 0 {
		anydraw = 1
		copy(range_,newpts)
	}else {
		range_[0] = 10000
		range_[1] = 10000
		range_[2] = -10000
		range_[3] = -10000
	}
	
	for loop = 1; loop < width; loop +=2 {
		newpts[0] = pts[0] + xinc*(loop/2+1)
		newpts[1] = pts[1] + yinc*(loop/2+1)
		newpts[2] = pts[2] + xinc*(loop/2+1)
		newpts[3] = pts[3] + yinc*(loop/2+1)
		if clip_and_draw_line(surf,rect,col,newpts) > 0 {
			anydraw = 1
			range_[0] = min(newpts[0],range_[0])
			range_[1] = min(newpts[1],range_[1])
			range_[2] = max(newpts[2],range_[2])
			range_[3] = max(newpts[3],range_[3])
		}
		if (loop + 1) < width {
			newpts[0] = pts[0] - xinc*(loop/2+1)
			newpts[1] = pts[1] - yinc*(loop/2+1)
			newpts[2] = pts[2] - xinc*(loop/2+1)
			newpts[3] = pts[3] - yinc*(loop/2+1)
			if clip_and_draw_line(surf,rect,col, newpts) > 0 {
				anydraw = 1
				range_[0] = min(newpts[0],range_[0])
				range_[1] = min(newpts[1],range_[1])
				range_[2] = max(newpts[2],range_[2])
				range_[3] = max(newpts[3],range_[3])
			}
		}
	}
	if anydraw > 0 {
		copy(pts,range_)
	}
	return anydraw
}

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

func encode(x,y,left,top,right,bottom int) int {
	code := 0
	if (x < left ) {
		code |= LEFT_EDGE
	}
	if (x > right) {
		code |= RIGHT_EDGE
	}
	if (y < top)   {
		code |= TOP_EDGE
	}
	if (y > bottom) {
		code |= BOTTOM_EDGE
	}
	return code
}

func inside(a int) bool {
	if a > 0 {
		return false
	}
	return true
}
		
func accept(a,b int) bool {
	ret := a | b
	if ret > 0 {
		return false
	}else {
		return true
	}
}

func reject(a,b int) bool {
	ret := a & b
	if ret > 0 {
		return true
	}
	return false
}

func clipline(pts []int, left,top,right,bottom int)  int {

	x1 := pts[0]
	y1 := pts[1]
	x2 := pts[2]
	y2 := pts[3]
	
	var code1 int
	var code2 int
	draw := 0
	var swaptmp int
	var m float64 /*slope*/

	for true {
		code1 = encode(x1,y1,left,top,right,bottom)
		code2 = encode(x2,y2,left,top,right,bottom)
		if ( accept(code1,code2) ) {
			draw = 1
			break
		} else if ( reject(code1,code2 ) ) {
			break
		}else {
			if inside(code1) {
				swaptmp = x2
				x2 = x1
				x1 = swaptmp
				swaptmp = y2
				y2 = y1
				y1 = swaptmp
				swaptmp = code2
				code2 = code1
				code1 = swaptmp
			}
			if x2 != x1 {
				m = float64(y2 - y1) / float64(x2-x1)
			}else {
				m = 1.0
			}
			if (code1 & LEFT_EDGE) > 0 {
				y1 += int(float64(left-x1)*m)
				x1 = left
			}else if (code1 & RIGHT_EDGE) > 0 {
				y1 += int(float64(right-x1)*m)
				x1 = right
			}else if (code1 & BOTTOM_EDGE) > 0 {
				if x2 != x1 {
					x1 += int(float64(bottom-y1) / m)
				}
				y1 = bottom
			}else if (code1 & TOP_EDGE) > 0 {
				if x2 != x1 {
					x1 += int( float64(top-y1) / m)
				}
				y1 = top
			}
		}
	}

	if draw > 0 {
		pts[0] = x1
		pts[1] = y1
		pts[2] = x2
		pts[3] = y2
	}
	
	return draw
}

func drawline(surf *sdl.Surface, col color.Color, x1,y1,x2,y2 int) {
	deltax := x2 - x1
	deltay := y2 - y1

	signx := 0
	signy := 0
	
	if deltax < 0 {
		signx = -1
	}else {
		signx = 1
	}

	if deltay < 0 {
		signy = -1
	}else {
		signy = 1
	}

	deltax = signx * deltax + 1
	deltay = signy * deltay + 1

	bytes_per_pixel := surf.BytesPerPixel()

	pixx := int(bytes_per_pixel)
	pixy := int(surf.Pitch)

	addr := pixy* y1 + x1 * bytes_per_pixel

	pixx *= int(signx)
	pixy *= int(signy)

  swaptmp := 0 
	if deltax < deltay {
		swaptmp = deltax
		deltax = deltay
		deltay = swaptmp
		swaptmp = pixx
		pixx = pixy
		pixy = swaptmp
	}

	x := 0
	y := 0

	color_bytes := col.ToBytes()
	pixels := surf.Pixels()
	
	switch bytes_per_pixel {
	case 1:
		for ; x < deltax; x++ {
			pixels[addr] = color_bytes[0]
			y += deltay
			if y >= deltax {
				y -= deltax
				addr += pixy
			}
			addr +=pixx
		}
		break
	case 2:
		for ; x < deltax;x++  {
			pixels[addr] = color_bytes[0]
			pixels[addr+1] = color_bytes[1]
			y+= deltay
			if  y >= deltax {
				y -= deltax
				addr += pixy
			}

			addr+=pixx
		}
		break
	case 3:
		for ; x < deltax; x++ {
			pixels[addr] = color_bytes[0]
			pixels[addr+1] = color_bytes[1]
			pixels[addr+2] = color_bytes[2]
			y+=deltay
			if y >= deltax {
				y-=deltax
				addr += pixy
			}
			addr+=pixx
		}
		break
	case 4:
		for ; x < deltax; x++ {
			pixels[addr] = color_bytes[0]
			pixels[addr+1] = color_bytes[1]
			pixels[addr+2] = color_bytes[2]
			pixels[addr+3] = color_bytes[3]
			y+=deltay
			if y >= deltax {
				y-=deltax
				addr += pixy
			}
			addr+=pixx
		}
		break		
	}
	
}

func drawhorzline(surf *sdl.Surface, col color.Color, x1,y1,x2 int) {
	if x1 == x2 {
		pixel(surf,col,x1,y1)
		return
	}
	
	bytes_per_pixel := surf.BytesPerPixel()	
	color_bytes := col.ToBytes()
	pixels := surf.Pixels()

	addr := int(surf.Pitch) * y1
	end := 0
	start := 0
	if x1 < x2 {
		end   = addr + x2*bytes_per_pixel
		start = addr+x1 *bytes_per_pixel
	}else {
		end  = addr + x1 *bytes_per_pixel
		start = addr + x2 * bytes_per_pixel
	}

	switch bytes_per_pixel {
	case 1:
		for ; start <=end; start++ {
			pixels[start] = color_bytes[0]
		}
	case 2:
		for ; start <= end; start+=2 {
			pixels[start] = color_bytes[0]
			pixels[start+1] = color_bytes[1]
		}
	case 3:
		for ; start <= end; start+=3 {
			pixels[start] = color_bytes[0]
			pixels[start+1] = color_bytes[1]
			pixels[start+2] = color_bytes[2]
		}
	case 4:
		for ; start <= end; start +=4 {
			pixels[start] = color_bytes[0]
			pixels[start+1] = color_bytes[1]
			pixels[start+2] = color_bytes[2]
			pixels[start+3] = color_bytes[3]
		}
	}

}

func drawvertline(surf *sdl.Surface, col color.Color, x1,y1,y2 int) {
	if y1 == y2 {
		pixel(surf,col, x1,y1)
	}
	bytes_per_pixel := surf.BytesPerPixel()	
	color_bytes := col.ToBytes()
	pixels := surf.Pixels()
	pitch  := int(surf.Pitch)
	
	addr := x1 * bytes_per_pixel
	end := 0
	start := 0
	if y1 < y2 {
		end = addr + y2* pitch
		start = addr + y1*pitch
	}else {
		end = addr + y1*pitch
		start = addr + y2*pitch
	}

	switch bytes_per_pixel {
	case 1:
		for ; start <=end; start+=pitch {
			pixels[start] = color_bytes[0]
		}
	case 2:
		for ; start <= end; start+=pitch {
			pixels[start] = color_bytes[0]
			pixels[start+1] = color_bytes[1]
		}
	case 3:
		for ; start <= end; start+=pitch {
			pixels[start] = color_bytes[0]
			pixels[start+1] = color_bytes[1]
			pixels[start+2] = color_bytes[2]
		}
	case 4:
		for ; start <= end; start +=pitch {
			pixels[start] = color_bytes[0]
			pixels[start+1] = color_bytes[1]
			pixels[start+2] = color_bytes[2]
			pixels[start+3] = color_bytes[3]
		}
	}
	
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

func Point(surf *sdl.Surface, c color.Color, x,y int) {
	pixels := surf.Pixels()
	bytes_per_pixel := surf.BytesPerPixel()
	
	addr := y * int(surf.Pitch) + x*bytes_per_pixel // 1 2 3 4

	color_bytes := c.ToBytes()

	surf.Lock()
	
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
	
	surf.Unlock()
}
