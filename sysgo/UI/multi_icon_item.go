package UI

import (
	"github.com/cuu/gogame/image"
)


type MultiIconItem struct {
	IconItem
	
	IconWidth int
	IconHeight int
	IconIndex int
}

func NewMultiIconItem() *MultiIconItem {
	m := &MultiIconItem{}
	m.IconIndex = 0
	m.IconWidth = 18
	m.IconHeight = 18

	return m 
}



func (m * MultiIconItem) CreateImageSurf() {
	if m.ImgSurf == nil and m.ImageName != "" {
		m.ImgSurf = image.Load(m.ImageName)
	}
}

func (m *MultiIconItem) Draw() {
	
}
