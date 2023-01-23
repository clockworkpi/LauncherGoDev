package MusicPlayer

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/color"
	
//	"github.com/clockworkpi/LauncherGoDev/sysgo"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"

//	"github.com/fhs/gompd/v2/mpd"
	
)

type MusicLibListPage struct {
	UI.Page
	ListFontObj *ttf.Font
	URLColor    *color.Color
	TextColor   *color.Color
	Labels      map[string]UI.LabelInterface
	Icons       map[string]UI.IconItemInterface

	IP     string

        //MyList   []UI.ListItemInterface
	MyStack        *UI.FolderStack
        BGwidth  int
        BGheight int //70	
        Scroller *UI.ListScroller
        Scrolled int
	
        Parent *MusicPlayerPage// also use the MpdClient from 

}

func NewMusicLibListPage() *MusicLibListPage {
	p := &MusicLibListPage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav.", "", "Scan","Back","Add to Playlist"}


	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")
	p.ListFontObj = UI.MyLangManager.TrFont("notosanscjk15")

	p.Labels = make(map[string]UI.LabelInterface)

	p.Icons = make(map[string]UI.IconItemInterface)

	p.BGwidth = 56
	p.BGheight = 70
	
	
	return p
}


func (self *MusicLibListPage) OnLoadCb() {
//	self.PosY = 0
	if self.MyList == nil || len(self.MyList) == 0 {
		self.MyStack.Clear()
		self.SyncList("/")
	}
}

func (self *MusicLibListPage) SetCoords() {
}

func (self *MusicLibListPage) SyncList(path string) {
	fmt.Println("SyncList: ",path)
	fmt.Println(self.MyStack)
	conn := self.Parent.MpdClient
	
	self.MyList = nil

	start_x := 0
	start_y := 0
	hasparent :=0 

	atts, err := conn.ListInfo(path)
	if err != nil {
		log.Println(err)
		return
	}
	
	if self.MyStack.Length() > 0 {
		hasparent = 1
                li := NewMusicLibListPageListItem()
                li.Parent = self
                li.PosX = start_x
                li.PosY = start_y
                li.Width = UI.Width
                li.Fonts["normal"] = self.ListFontObj
		li.Path = "[..]"
                li.Init("[..]")
		li.MyType = UI.ICON_TYPES["DIR"]
                self.MyList = append(self.MyList, li)		
	}

	if len(atts) == 0 {
		log.Println("no songs")
		return
	}
	
	for i, m := range atts {

		li := NewMusicLibListPageListItem()
		li.Parent = self
                li.PosX = start_x
                li.PosY = start_y + (i+hasparent)*li.Height
                li.Width = UI.Width
                li.Fonts["normal"] = self.ListFontObj
                li.MyType = UI.ICON_TYPES["FILE"]
		
		init_val := "NoName"

		if val, ok := m["directory"] ; ok {
			li.MyType = UI.ICON_TYPES["DIR"]
			init_val = filepath.Base(val)
			li.Path = val
		}

		if val, ok := m["file"]; ok {
			li.MyType = UI.ICON_TYPES["FILE"]
			li.Path = val

			val2, ok2 := m["Title"]
			if ok2  && len(val2) > 4{
				init_val = val2
			}else{
				init_val = val
			}
		}

		li.Init(init_val)
		self.MyList = append(self.MyList, li)

	}
	
}

func (self *MusicLibListPage) Init() {
	if self.Screen == nil {
		panic("No Screen")
	}

	if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
		self.HWND = self.Screen.CanvasHWND
		self.CanvasHWND = surface.Surface(self.Screen.Width, self.Screen.Height)
	}

	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	ps := NewListPageSelector()
        ps.Parent = self

        ps.Width = UI.Width - 12
        ps.PosX = 2
        ps.Parent = self

        self.Ps = ps
        self.PsIndex = 0

	bgpng := UI.NewIconItem()
        bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("empty")
        bgpng.MyType = UI.ICON_TYPES["STAT"]
        bgpng.Parent = self
        bgpng.AddLabel("Please upload data over Wi-Fi", UI.Fonts["varela22"])
        bgpng.SetLabelColor(&color.Color{204, 204, 204, 255})
        bgpng.Adjust(0, 0, self.BGwidth, self.BGheight, 0)

        self.Icons["bg"] = bgpng

        icon_for_list := UI.NewMultiIconItem()
        icon_for_list.ImgSurf = UI.MyIconPool.GetImgSurf("sys")
        icon_for_list.MyType = UI.ICON_TYPES["STAT"]
        icon_for_list.Parent = self

        icon_for_list.Adjust(0, 0, 18, 18, 0)

        self.Icons["sys"] = icon_for_list

        self.Scroller = UI.NewListScroller()
        self.Scroller.Parent = self
        self.Scroller.PosX = self.Width - 10
        self.Scroller.PosY = 2
        self.Scroller.Init()

	self.MyStack = UI.NewFolderStack()
	self.MyStack.SetRootPath("/")

}

func (self *MusicLibListPage) Click() {
	self.RefreshPsIndex()
	
        if  len(self.MyList) == 0 {
                return
        }

	cur_li := self.MyList[self.PsIndex].(*MusicLibListPageListItem)
	
        if cur_li.MyType == UI.ICON_TYPES["DIR"] {
            if cur_li.Path == "[..]" {
                self.MyStack.Pop()
                self.SyncList( self.MyStack.Last() )
                self.PsIndex = 0
	    } else {
                self.MyStack.Push( cur_li.Path )
                self.SyncList( self.MyStack.Last() )
                self.PsIndex = 0
            }
    	}

        if cur_li.MyType == UI.ICON_TYPES["FILE"] {
		conn := self.Parent.MpdClient
		conn.Add(cur_li.Path)
		self.Parent.SyncList()	
            	fmt.Println("add" , cur_li.Path)

	}

        self.Screen.Refresh()

}

func (self *MusicLibListPage) KeyDown(ev *event.Event) {
	
        if UI.IsKeyMenuOrA(ev.Data["Key"]) || ev.Data["Key"] == UI.CurKeys["Left"] {
                self.ReturnToUpLevelPage()
                self.Screen.Refresh()
        }

        if ev.Data["Key"] == UI.CurKeys["Up"] {

                self.ScrollUp()
                self.Screen.Refresh()
        }

        if ev.Data["Key"] == UI.CurKeys["Down"] {

                self.ScrollDown()
                self.Screen.Refresh()
        }

	if ev.Data["Key"] == UI.CurKeys["B"] {
		self.Click()
	}

	if ev.Data["Key"] == UI.CurKeys["Y"] {
		self.Screen.ShowMsg("Scan...",300)
		self.OnLoadCb()
	}
	return
}

func (self *MusicLibListPage) Draw() {
	self.ClearCanvas()

        if len(self.MyList) == 0 {
                self.Icons["bg"].NewCoord(self.Width/2, self.Height/2)
                self.Icons["bg"].Draw()
	}else {
		if len(self.MyList)*UI.DefaultInfoPageListItemHeight > self.Height {
			self.Ps.(*ListPageSelector).Width = self.Width - 11
	                self.Ps.Draw()

			for _, v := range self.MyList {
				v.(*MusicLibListPageListItem).Active = false

				if self.Parent.InPlayList( v.(*MusicLibListPageListItem).Path) {
					v.(*MusicLibListPageListItem).Active = true
					fmt.Println("in PlayList: ",v.(*MusicLibListPageListItem).Path)
				}

	                        if v.(*MusicLibListPageListItem).PosY > self.Height+self.Height/2 {
        	                        break
	                        }

        	                if v.(*MusicLibListPageListItem).PosY < 0 {
                	                continue
                        	}

	                        v.Draw()				
			}

                self.Scroller.UpdateSize( len(self.MyList)*UI.DefaultInfoPageListItemHeight, self.PsIndex*UI.DefaultInfoPageListItemHeight)
                self.Scroller.Draw()

		} else{
	                self.Ps.(*ListPageSelector).Width = self.Width
        	        self.Ps.Draw()
                	for _, v := range self.MyList {
                                v.(*MusicLibListPageListItem).Active = false

                                if self.Parent.InPlayList( v.(*MusicLibListPageListItem).Path) {
                                        v.(*MusicLibListPageListItem).Active = true
                                        fmt.Println("in PlayList: ",v.(*MusicLibListPageListItem).Path)
                                }

                        	if v.(*MusicLibListPageListItem).PosY > self.Height+self.Height/2 {
                                	break
	                        }

        	                if v.(*MusicLibListPageListItem).PosY < 0 {
                	                continue
                       	 	}
	                        v.Draw()
	                }
		}
	}

        if self.HWND != nil {
                surface.Fill(self.HWND, UI.MySkinManager.GiveColor("White"))
                rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
                surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
        }

}
