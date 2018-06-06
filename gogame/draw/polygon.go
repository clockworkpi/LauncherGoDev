package draw

import (
	"fmt"
//	"math"
	"github.com/veandco/go-sdl2/sdl"
	"../color"
	"../rect"
	"../qsort"
	
)

func draw_fillpoly(surf *sdl.Surface, vx []int, vy []int, numpoints int, col color.Color) {

	miny:=0
	maxy:=0
	y:=0
	x1:=0
	y1:=0
	x2:=0
	y2:=0
	ints := 0
	ind1:=0
	ind2:=0
	
	polyints := make([]int,numpoints)
	
	/* Determine Y maxima */
	miny = vy[0]
	maxy = vy[0]
	for i:=1; i < numpoints; i++ {
		miny = min(miny,vy[i])
		maxy = max(maxy,vy[i])
	}

	/* Draw, scanning y */
	for y=miny;y<=maxy;y++ {
		ints = 0
		for i:=0; i< numpoints;i++ {
			if i == 0 {
				ind1 = numpoints -1
				ind2 = 0
			}else {
				ind1 = i-1
				ind2 = i
			}

			y1 = vy[ind1]
			y2 = vy[ind2]
			if y1 < y2 {
				x1 = vx[ind1]
				x2 = vx[ind2]
			}else if y1 > y2 {
				y2 = vy[ind1]
				y1 = vy[ind2]
				x2 = vx[ind1]
				x1 = vx[ind2]
			}else if miny == maxy {
				/* Special case: polygon only 1 pixel high. */
				minx:= vx[0]
				maxx:= vx[0]
				
				for j:= 1; j < numpoints; j++ {
					minx = min(minx,vx[j])
					maxx = max(maxx,vx[j])
				}
				polyints[ints] = minx
				ints+=1
				polyints[ints] = maxx
				ints+=1
				break
			}else {
				continue
			}
			
			if ( y >= y1) && (y < y2 ) {
				fmt.Println("ints : ",ints)
				polyints[ints] = (y-y1) * (x2-x1) / (y2-y1) + x1
				ints+=1
				
			}else if (y == maxy) && (y > y1) && (y <= y2) {				
				polyints[ints] = (y-y1) * (x2-x1) / (y2-y1) + x1
				ints+=1
			}
		}

		new_polyints := make([]int, ints)
		copy(new_polyints,polyints)		
		new_polyints = qsort.QuickSort(new_polyints)
		
		for i:=0;i<ints;i+=2 {
			drawhorzlineclip(surf, col, new_polyints[i], y, new_polyints[i+1])
		}
	}
}

func Polygon(surf *sdl.Surface, color color.Color, points [][]int, border_width int) sdl.Rect {

	if border_width > 0 {
		ret := Lines(surf, color, true,points,border_width)
		return ret
	}

	bytes_per_pixel := surf.BytesPerPixel()
	if bytes_per_pixel <= 0 || bytes_per_pixel > 4 {
		panic("unsupport bit depth for line draw")
	}

	length := len(points)
	if length < 3 {
		panic("points argument must contain more than 2 points")
	}

	item := points[0] // Needs sequence_get_item to fetch from sorted points list
	if len(item) < 2 {
		panic("points should be a pair of coordinators")
	}

	x := item[0]
	y := item[1]

	left   := x
	right  := x
	top    := y
	bottom := y

	xlist := make([]int,length)
	ylist := make([]int,length)

	numpoints := 0
	for loop := 0; loop < length; loop++ {
		item = points[loop]
		if len(item) < 2 {
			panic("points should be a pair of coordinators")
		}
		x = item[0]
		y = item[1]
		
		xlist[numpoints] = x
		ylist[numpoints] = y
		numpoints+=1

		left   = min(x,left)
		top    = min(y,top)
		right  = max(x,right)
		bottom = max(y,bottom)
	}

	err := surf.Lock()
	if err != nil {
		return rect.Rect(0,0,0,0)
	}

	draw_fillpoly(surf,xlist,ylist,numpoints,color)

	surf.Unlock()

	left = max(left,int(surf.ClipRect.X))
	top  = max(top, int(surf.ClipRect.Y))
	right = min(right,int(surf.ClipRect.X + surf.ClipRect.W))
	bottom = min(bottom, int(surf.ClipRect.Y + surf.ClipRect.H))
	return rect.Rect(left,top,right-left+1, bottom-top+1)
	
}
