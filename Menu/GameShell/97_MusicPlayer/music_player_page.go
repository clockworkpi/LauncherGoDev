package MusicPlayer

import (
	"fmt"
	"path/filepath"
	"log"
	"strconv"
	"strings"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	"github.com/cuu/gogame/surface"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/color"
	"github.com/clockworkpi/LauncherGoDev/sysgo"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"

	"github.com/fhs/gompd/v2/mpd"
)

type MusicPlayerPage struct {
	UI.Page
	ListFontObj *ttf.Font
	URLColor    *color.Color
	TextColor   *color.Color
	Labels      map[string]UI.LabelInterface
	Icons       map[string]UI.IconItemInterface

	IP     string
	
	MyMusicLibListPage *MusicLibListPage //also use the MpdClient *mpd.Client

        //MyList   []UI.ListItemInterface
	MyStack        *UI.FolderStack
        BGwidth  int
        BGheight int //70	
        Scroller *UI.ListScroller
        Scrolled int

	MpdClient *mpd.Client

	CurSongTime string
	CurSongName string
}

func NewMusicPlayerPage() *MusicPlayerPage {
	p := &MusicPlayerPage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav","Remove","","Back","Play/Pause"}


	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")
	p.ListFontObj = UI.MyLangManager.TrFont("notosanscjk15")

	p.Labels = make(map[string]UI.LabelInterface)

	p.Icons = make(map[string]UI.IconItemInterface)

	p.BGwidth = 56
	p.BGheight = 70
	
	p.CurSongTime = "0:0"

	
	return p
}

func (self *MusicPlayerPage) OnLoadCb() {
	self.PosY = 0

	if self.MpdClient == nil {
	        conn, err := mpd.Dial("unix", sysgo.MPD_socket)
	        if err != nil {
                	log.Fatalln(err)
        	}
		self.MpdClient = conn

		fmt.Println("Start mpd client")
	}

	self.SyncList()
}

func (self *MusicPlayerPage) OnPopUpCb() {
	if self.MpdClient != nil {
		self.MpdClient.Close()
		self.MpdClient = nil
		fmt.Println("Close mpd client")
	}
}

func (self *MusicPlayerPage) OnReturnBackCb() {
	self.SyncList()
}

func (self *MusicPlayerPage) SetCoords() {

}

func (self *MusicPlayerPage) SetLabels() {

}

func (self *MusicPlayerPage) SyncList() {
	conn := self.MpdClient
	start_x := 0
	start_y := 0

	if conn == nil {
		return
	}
	
	self.MyList = nil

	play_list,_ := conn.PlaylistInfo(-1,-1)

	for i,v := range play_list {
		li := NewMusicPlayPageListItem()
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y + UI.DefaultInfoPageListItemHeight * i  
		li.Width = UI.Width
		li.Fonts["normal"] = self.ListFontObj
		
		if val,ok:=v["Title"]; ok {
			li.Init(val)

			if val2,ok2 := v["file"]; ok2 {
				li.Path = val2
			}
		}else {
			if val2,ok2 := v["file"]; ok2 {
				li.Init(filepath.Base(val2))
				li.Path = val2
			}else{
				li.Init("NoName")
			}
		}
		x,_ := li.Labels["Text"].Coord()	
		li.Labels["Text"].NewCoord(x,7)
		self.MyList = append(self.MyList, li)
	}
	
	self.SyncPlaying()
}

func (self *MusicPlayerPage) SyncPlaying() {
	conn := self.MpdClient

	for i,_ := range self.MyList {
		self.MyList[i].(*MusicPlayPageListItem).Active = false
		self.MyList[i].(*MusicPlayPageListItem).PlayingProcess = 0
	}
	current_song,_ := conn.CurrentSong()
	if len(current_song) > 0 {
		if val,ok := current_song["song"]; ok{
			posid, _ := strconv.Atoi(val)
			if posid < len(self.MyList) {
				if state,ok2 := current_song["state"]; ok2 {
					if state == "stop" {
						self.MyList[posid].(*MusicPlayPageListItem).Active = false
					}else{
						self.MyList[posid].(*MusicPlayPageListItem).Active = true
					}
				}

				if song_time,ok3 := current_song["time"]; ok3 {
					self.CurSongTime = song_time
					times := strings.Split(self.CurSongTime,":")
					if len(times) > 1{
						cur,_ := strconv.ParseFloat(times[0],64)
						end,_ := strconv.ParseFloat(times[1],64)
						pos := int( (cur/end)*100.0 )
						self.MyList[posid].(*MusicPlayPageListItem).PlayingProcess = pos
					}
				}
			}
		}
	}
}

func (self *MusicPlayerPage) InPlayList(path string) bool {
	if self.MyList == nil || len(self.MyList) == 0 {
		return false
	}

	for _,v := range self.MyList {
		///fmt.Println(v.(*MusicPlayPageListItem).Path, path)

		if v.(*MusicPlayPageListItem).Path == path {
			return true
		}
	}

	return false
}

func (self *MusicPlayerPage) Init() {
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

	//ps := UI.NewInfoPageSelector()
	ps := NewListPageSelector()
        //ps.Width = UI.Width - 12
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
	
	self.MyMusicLibListPage = NewMusicLibListPage()
        self.MyMusicLibListPage.Screen = self.Screen
        self.MyMusicLibListPage.Name = "Music library"
        self.MyMusicLibListPage.Parent = self
        self.MyMusicLibListPage.Init()

	self.MyStack = UI.NewFolderStack()
	self.MyStack.SetRootPath("/")
}

func (self *MusicPlayerPage) KeyDown(ev *event.Event) {
        if ev.Data["Key"] == UI.CurKeys["Right"] {
               	self.Screen.PushPage(self.MyMusicLibListPage) 
                self.Screen.Draw()
                self.Screen.SwapAndShow()
        }	
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

        if ev.Data["Key"] == UI.CurKeys["Up"] {

                self.ScrollUp()
                self.Screen.Draw()
                self.Screen.SwapAndShow()
        }

        if ev.Data["Key"] == UI.CurKeys["Down"] {

                self.ScrollDown()
                self.Screen.Draw()
                self.Screen.SwapAndShow()
        }
	if ev.Data["Key"] == UI.CurKeys["X"] {
		self.MpdClient.Delete(self.PsIndex,-1)
		self.SyncList()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}
	return
}

func (self *MusicPlayerPage) Draw() {
	self.ClearCanvas()

        if len(self.MyList) == 0 {
                self.Icons["bg"].NewCoord(self.Width/2, self.Height/2)
                self.Icons["bg"].Draw()
        }else {
                if len(self.MyList)*UI.DefaultInfoPageListItemHeight > self.Height {
                        self.Ps.(*ListPageSelector).Width = self.Width - 11
                        self.Ps.Draw()

                        for _, v := range self.MyList {
                                if v.(*MusicPlayPageListItem).PosY > self.Height+self.Height/2 {
                                        break
                                }

                                if v.(*MusicPlayPageListItem).PosY < 0 {
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
                                if v.(*MusicPlayPageListItem).PosY > self.Height+self.Height/2 {
                                        break
                                }

                                if v.(*MusicPlayPageListItem).PosY < 0 {
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
