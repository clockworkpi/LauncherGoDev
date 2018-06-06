package draw

import (
//	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"../color"
	"../rect"
)

//closed true => an additional line segment is drawn between the first and last points.
//pointlist should be [][2] 
func Lines(surf *sdl.Surface, col color.Color,closed bool, pointlist [][]int,width int) sdl.Rect {
	length := len(pointlist)
	if length < 2 {
		panic("draw lines at least contains more than 1 points pair")
	}
	pts := make([]int,4)

	if len(pointlist[0]) < 2 {
		panic("start points should be more than 1 at least")
	}
	
	x := pointlist[0][0]
	y := pointlist[0][1]
	
	startx := x
	pts[0]  = x
	left   := x
	right  := x

	starty := y
	pts[1]  = y
	top    := y
	bottom := y

	if width < 1 {
		return rect.Rect(x,y,0,0)
	}

	err := surf.Lock()
	if err != nil {
		return rect.Rect(0,0,0,0)
	}

	drawn := 1
	for loop := 1; loop < length; loop++ {
		item := pointlist[loop]
		if len(item) < 2 {
			continue
		}

		x = item[0]
		y = item[1]
		drawn += 1
		pts[0] = startx
		pts[1] = starty
		startx = x
		starty = y
		pts[2] = x
		pts[3] = y
		if clip_and_draw_line_width(surf, &surf.ClipRect, col, width, pts) > 0 {
			left =   min(min(pts[0],pts[2]),left)
			top  =   min(min(pts[1],pts[3]),top)
			right =  max(max(pts[0],pts[2]),right)
			bottom = max(max(pts[1],pts[3]),bottom)
		}
	}

	if closed == true && drawn > 2 {
		item := pointlist[0]
		x = item[0]
		y = item[1]

		pts[0] = startx
		pts[1] = starty
		pts[2] = x
		pts[3] = y
		clip_and_draw_line_width(surf, &surf.ClipRect, col, width, pts)
	}
	surf.Unlock()

	return rect.Rect(left,top,right-left+1, bottom-top+1)
	
}

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

func drawhorzlineclip(surf *sdl.Surface, col color.Color, x1 , y1, x2 int ) {
	if y1 < int(surf.ClipRect.Y) || y1 >= int( surf.ClipRect.Y + surf.ClipRect.H) {
		return
	}

	if x2 < x1 {
		temp := x1
		x1 = x2
		x2 = temp
	}
	x1 = max(x1,int(surf.ClipRect.X))
	x2 = min(x2,int(surf.ClipRect.X+surf.ClipRect.W-1))
	if x2 < int(surf.ClipRect.X) || x1 >= int( surf.ClipRect.X + surf.ClipRect.W) {
		return
	}

	if x1 == x2 {
		pixel(surf,col, x1,y1)
	}else {
		drawhorzline(surf,col,x1,y1,x2)
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

