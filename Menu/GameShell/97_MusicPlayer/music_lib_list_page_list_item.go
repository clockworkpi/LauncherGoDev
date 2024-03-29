package MusicPlayer

import (
//        "fmt"
        //"io/ioutil"
        //"path/filepath"
        "github.com/veandco/go-sdl2/ttf"
//        "runtime"
//        "strconv"
//        "strings"
        //"github.com/mitchellh/go-homedir"

 ///       "github.com/clockworkpi/LauncherGoDev/sysgo"
        "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
        "github.com/cuu/gogame/color"
        "github.com/cuu/gogame/draw"
//        "github.com/cuu/gogame/event"
//        "github.com/cuu/gogame/rect"
//        "github.com/cuu/gogame/surface"
//        "github.com/cuu/gogame/time"
)

type MusicLibListPageListItem struct {
        UI.InfoPageListItem

        Active bool
	Value  string
	MyType int
	Path  string
}

func NewMusicLibListPageListItem() *MusicLibListPageListItem {

        p := &MusicLibListPageListItem{}
        p.Height = UI.DefaultInfoPageListItemHeight
        p.ReadOnly = false
        p.Labels = make(map[string]UI.LabelInterface)
        p.Icons = make(map[string]UI.IconItemInterface)
        p.Fonts = make(map[string]*ttf.Font)
	p.MyType = UI.ICON_TYPES["EXE"]
        return p
}

func (self *MusicLibListPageListItem) Init(text string) {
        l := UI.NewLabel()
        l.PosX = 10
        l.SetCanvasHWND(self.Parent.GetCanvasHWND())
        l.Init(text, self.Fonts["normal"], nil)
        self.Labels["Text"] = l	
}

func (self *MusicLibListPageListItem) Draw() {

        x, _ := self.Labels["Text"].Coord()
        w, h := self.Labels["Text"].Size()

        self.Labels["Text"].NewCoord(x, self.PosY+(self.Height-h)/2)

	if self.MyType == UI.ICON_TYPES["DIR"] && self.Path != "[..]" {
	    sys_icon := self.Parent.(*MusicLibListPage).Icons["sys"]
            sys_icon.SetIconIndex(0)
	    _,h := sys_icon.Size()
            sys_icon.NewCoord(self.Width - 12,self.PosY+ (self.Height - h)/2+h/2)
            sys_icon.Draw()
	}

        if self.MyType == UI.ICON_TYPES["FILE"] {
	    sys_icon := self.Parent.(*MusicLibListPage).Icons["sys"]
            sys_icon.SetIconIndex(1)
	    _,h := sys_icon.Size()
            sys_icon.NewCoord(self.Width-12,self.PosY+ (self.Height - h)/2+h/2)
            sys_icon.Draw()
	}
	/*
        if self.Active == true {
                self.Parent.(*MusicLibListPage).Icons["sys"].NewCoord(self.Parent.(*MusicLibListPage).Width-30, self.PosY+5)
                self.Parent.(*MusicLibListPage).Icons["sys"].Draw()
        }
	*/

        self.Labels["Text"].SetBold(self.Active)
        self.Labels["Text"].Draw()

        if _, ok := self.Labels["Small"]; ok {
                x, _ = self.Labels["Small"].Coord()
                w, h = self.Labels["Small"].Size()

                self.Labels["Small"].NewCoord(self.Width-w-10, self.PosY+(self.Height-h)/2)
                self.Labels["Small"].Draw()

        }

        canvas_ := self.Parent.GetCanvasHWND()
        draw.Line(canvas_, &color.Color{169, 169, 169, 255},
                self.PosX, self.PosY+self.Height-1,
                self.PosX+self.Width, self.PosY+self.Height-1, 1)

}
