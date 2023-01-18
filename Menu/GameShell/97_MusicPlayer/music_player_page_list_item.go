package MusicPlayer

import (
        //"fmt"
        //"io/ioutil"
        //"path/filepath"
        "github.com/veandco/go-sdl2/ttf"
        //"runtime"
        //"strconv"
        //"strings"
        //"github.com/mitchellh/go-homedir"

        //"github.com/clockworkpi/LauncherGoDev/sysgo"
        "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
        //"github.com/cuu/gogame/color"
        "github.com/cuu/gogame/draw"
        //"github.com/cuu/gogame/event"
        //"github.com/cuu/gogame/rect"
        //"github.com/cuu/gogame/surface"
        //"github.com/cuu/gogame/time"
)

type MusicPlayPageListItem struct {
        UI.InfoPageListItem

        Active bool
	Value  string
	MyType int
	Path  string

	PlayingProcess int
}

func NewMusicPlayPageListItem() *MusicPlayPageListItem {

        p := &MusicPlayPageListItem{}
        p.Height = UI.DefaultInfoPageListItemHeight
        p.ReadOnly = false
        p.Labels = make(map[string]UI.LabelInterface)
        p.Icons = make(map[string]UI.IconItemInterface)
        p.Fonts = make(map[string]*ttf.Font)
	p.MyType = UI.ICON_TYPES["EXE"]
        return p
}

func (self *MusicPlayPageListItem) Draw() {

        x, _ := self.Labels["Text"].Coord()
        _, h := self.Labels["Text"].Size()

        self.Labels["Text"].NewCoord(x, self.PosY+(self.Height-h)/2)
	

	if self.MyType == UI.ICON_TYPES["DIR"] &&  self.Path != "[..]" {
		sys_icon := self.Parent.(*MusicPlayerPage).Icons["sys"]
		_,h := sys_icon.Size()
		sys_icon.SetIconIndex(0)
		sys_icon.NewCoord(self.PosX+12,self.PosY + ( self.Height - h)/2 + h/2)
                sys_icon.Draw()		
	}

	if self.MyType == UI.ICON_TYPES["FILE"] {
		sys_icon := self.Parent.(*MusicPlayerPage).Icons["sys"]
		_,h := sys_icon.Size()
		sys_icon.SetIconIndex(1)
		sys_icon.NewCoord(self.PosX+12,self.PosY + ( self.Height - h)/2 + h /2)
		sys_icon.Draw()
	}


	self.Labels["Text"].NewCoord(x, self.PosY+(self.Height-h)/2)
	
        self.Labels["Text"].SetBold(self.Active)
        self.Labels["Text"].Draw()
	
	/*
        if _, ok := self.Labels["Small"]; ok {
                x, _ = self.Labels["Small"].Coord()
                w, h = self.Labels["Small"].Size()

                self.Labels["Small"].NewCoord(self.Width-w-10, self.PosY+(self.Height-h)/2)
                self.Labels["Small"].Draw()

        }
	*/
        canvas_ := self.Parent.GetCanvasHWND()	
        draw.Line(canvas_, UI.MySkinManager.GiveColor("Line"),
                self.PosX, self.PosY+self.Height-1,
                self.PosX+self.Width, self.PosY+self.Height-1, 1)

	if self.PlayingProcess > 0 {
		seek_posx := int(self.Width * self.PlayingProcess/100.0)	
		draw.Line(canvas_, UI.MySkinManager.GiveColor("Active"),
                	self.PosX, self.PosY+self.Height-2,
	                self.PosX+seek_posx, self.PosY+self.Height-2, 2)
	}
}
