package MusicPlayer

import (
        //"fmt"
//        "path/filepath"

//        "github.com/cuu/gogame/event"
	"github.com/cuu/gogame/draw"
        "github.com/cuu/gogame/rect"
//        "github.com/cuu/gogame/surface"
//        "github.com/veandco/go-sdl2/ttf"

        "github.com/cuu/gogame/color"

//        "github.com/clockworkpi/LauncherGoDev/sysgo"
        "github.com/clockworkpi/LauncherGoDev/sysgo/UI"


)

type ListPageSelector struct {
        UI.InfoPageSelector
}

func NewListPageSelector() *ListPageSelector {

        p := &ListPageSelector{}

        p.Width = UI.Width
        p.BackgroundColor = &color.Color{131, 199, 219, 255} //SkinManager().GiveColor('Front')

        return p

}

func (self *ListPageSelector) Draw() {

        idx := self.Parent.GetPsIndex()
        mylist := self.Parent.GetMyList()
        if idx < len(mylist) {
                x, y := mylist[idx].Coord()
                _, h := mylist[idx].Size()

                self.PosX = x + 2
                self.PosY = y + 1
                self.Height = h - 3

                canvas_ := self.Parent.GetCanvasHWND()
                rect_ := rect.Rect(self.PosX, self.PosY, self.Width-4, self.Height)

                draw.AARoundRect(canvas_, &rect_, self.BackgroundColor, 4, 0, self.BackgroundColor)
        }
}
