package Warehouse

import (
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"

)

//GameStoreListItem in py 
type WareHouseListItem struct {
    UI.InfoPageListItem

    Type  string
}




func NewWareHouseListItem() *WareHouseListItem {

	p := &WareHouseListItem{}
	p.Height = UI.DefaultInfoPageListItemHeight
	
	p.Labels = make(map[string]UI.LabelInterface)
	p.Icons = make(map[string]UI.IconItemInterface)
	p.Fonts = make(map[string]*ttf.Font)

	return p
}

func (self *WareHouseListItem) Init( text string) {

	l := UI.NewLabel()
	l.CanvasHWND = self.Parent.GetCanvasHWND()
	l.PosX = 10
	l.Init(text,self.Fonts["normal"])

	self.Labels["text"] = l

	add_icon := NewIconItem()
	add_icon.ImgSurf = UI.MyIconPool.GetImgSurf("add")
	add_icon.Parent = self
	add_icon.Init(0,0,UI.MyIconPool.Width("add"),UI.MyIconPool.Height("add"),0)

	ware_icon := NewIconItem()
	ware_icon.ImgSurf = UI.MyIconPool.GetImgSurf("ware")
	ware_icon.Parent = self
	ware_icon.Init(0,0,UI.MyIconPool.Width("ware"),UI.MyIconPool.Height("ware"),0)

	app_icon := NewIconItem()
	app_icon.ImgSurf = UI.MyIconPool.GetImgSurf("app")
	app_icon.Parent = self
	app_icon.Init(0,0,UI.MyIconPool.Width("app"),UI.MyIconPool.Height("app"),0)

	appdling_icon := NewIconItem()
	appdling_icon.ImgSurf = UI.MyIconPool.GetImgSurf("appdling")
	appdling_icon.Parent = self
	appdling_icon.Init(0,0,UI.MyIconPool.Width("appdling"),UI.MyIconPool.Height("appdling"),0)

	blackheart_icon := NewIconItem()
	blackheart_icon.ImgSurf = UI.MyIconPool.GetImgSurf("blackheart")
	blackheart_icon.Parent = self
	blackheart_icon.Init(0,0,UI.MyIconPool.Width("blackheart"),UI.MyIconPool.Height("blackheart"),0)

	self.Icons["add"] = add_icon
	self.Icons["ware"] = ware_icon
	self.Icons["app"] = app_icon
	self.Icons["appdling"] = appdling_icon
	self.Icons["blackheart"] = blackheart_icon
		
}

func (self *WareHouseListItem) Draw() {
	if self.ReadOnly == true {
		self.Labels["text"].SetColor( UI.MySkinManager.GiveColor("ReadOnlyText"))
	} else {
		self.Labels["text"].SetColor( UI.MySkinManager.GiveColor("Text"))
	}

	padding := 17

	if self.Type == "" {
		padding = 0
	}

	if self.Type == "source" || self.Type == "dir" {
		self.Icons["ware"].NewCoord(4,(self.Height - self.Icons["ware"].Height)/2)
		self.Icons["ware"].DrawTopLeft()
	}

	if self.Type == "launcher" || self.Type == "pico8" || self.Type == "tic80" {
		_icon :=  app
		if self.ReadOnly == true {
			_icon = "appdling"
		}
		self.Icons[_icon].NewCoord(4,(self.Height - self.Icons[_icon].Height)/2)
		self.Icons[_icon].DrawTopLeft()
	}

	if self.Type == "add_house" {
		self.Icons["add"].NewCoord(4,(self.Height-self.Icons["add"].Height)/2)
		self.Icons["add"].DrawTopLeft()
	}
	
	self.Labels["text"].PosX = self.Labels["text"].PosX + self.PosX + padding
	self.Labels["text"].PosY = self.PosY + (self.Height-self.Labels["text"].Height)/2
	self.Labels["text"].Draw()
	self.Labels["text"].PosX = self.Labels["text"].PosX - self.PosX - padding
	
	if _, ok := self.Labels["Small"]; ok {
		x, _ = self.Labels["Small"].Coord()
		w, h = self.Labels["Small"].Size()

		self.Labels["Small"].NewCoord(self.Width-w-5, self.PosY+(self.Height-h)/2)
		self.Labels["Small"].Draw()	
	}
	
	canvas_ := self.Parent.GetCanvasHWND()
	draw.Line(canvas_,UI.MySkinManager.GiveColor("Line"),
		self.PosX,self.PosY + self.Height -1,
		self.PosX+self.Width,self.PosY+self.Height-1,
		1)
	
}
