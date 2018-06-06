package draw

import (
//	"fmt"
//	"math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"

	"../color"
	"../rect"
)


func _aa_render_region(image *sdl.Renderer, _rect *sdl.Rect, col color.Color, rad int) {
	corners := rect.Inflate(_rect,-2*rad-1, -2*rad-1)
	topleft := []int{ int(corners.X),int(corners.Y)}
	topright := []int{int(corners.X+corners.W-1), int(corners.Y)}
	bottomleft := []int{int(corners.X),  int(corners.Y+corners.H-1)}
	bottomright := []int{int(corners.X+corners.W -1), int(corners.Y+corners.H-1)}

	attributes :=[][]int{topleft, topright, bottomleft, bottomright }

	r,g,b,a := col.RGBA()
	
	image.SetDrawColor( uint8(r),uint8(g),uint8(b),uint8(a) )

	
	for i:=0; i< len(attributes);i++ {
		x,y := attributes[i][0],attributes[i][1]
		
		gfx.AACircleRGBA(image,int32(x),int32(y),int32(rad),uint8(r),uint8(g),uint8(b),uint8(a))
		gfx.FilledCircleRGBA(image,int32(x),int32(y), int32(rad),uint8(r),uint8(g),uint8(b),uint8(a))
	}


	r1 := rect.Inflate(_rect,-2*rad,0)
	r2 := rect.Inflate(_rect,0,-2*rad)
	
	image.FillRect( &r1 ) // main body except four circles in corners
	image.FillRect( &r2 ) // fix gap between circles of up and down vertical
}

//alpha of color should be 255 
func AARoundRect(surf *sdl.Surface,_rect *sdl.Rect,col color.Color,rad,border int, inside color.Color) {

	image,_ := sdl.CreateSoftwareRenderer(surf)

	/*
	image.SetDrawColor(233,100,200,0)
	image.DrawLine(10,20,100,200)
  */
	
//	image.Clear()
	_aa_render_region(image,_rect,col,rad)
	if border > 0 {
		rect.InflateIp(_rect,-2*border,-2*border)
		_aa_render_region(image,_rect,inside,rad)
	}

	//image.Present()
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
