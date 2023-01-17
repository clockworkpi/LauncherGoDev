package MusicPlayer

import (
	"fmt"
	"log"
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
}

func NewMusicPlayerPage() *MusicPlayerPage {
	p := &MusicPlayerPage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav","Remove","RTA","Back","Play/Pause"}


	p.URLColor = UI.MySkinManager.GiveColor("URL")
	p.TextColor = UI.MySkinManager.GiveColor("Text")
	p.ListFontObj = UI.MyLangManager.TrFont("notosanscjk15")

	p.Labels = make(map[string]UI.LabelInterface)

	p.Icons = make(map[string]UI.IconItemInterface)

	p.BGwidth = 56
	p.BGheight = 70
	
	
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
}

func (self *MusicPlayerPage) OnPopUpCb() {
	if self.MpdClient != nil {
		self.MpdClient.Close()
		self.MpdClient = nil
		fmt.Println("Close mpd client")
	}
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
		
		li.Labels["Text"].PosX = 7
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

	ps := UI.NewInfoPageSelector()
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
	
	self.MyMusicLibListPage = NewMusicLibListPage()
        self.MyMusicLibListPage.Screen = self.Screen
        self.MyMusicLibListPage.Name = "Music Library"
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
	return
}

func (self *MusicPlayerPage) Draw() {
	self.ClearCanvas()

        if len(self.MyList) == 0 {
                self.Icons["bg"].NewCoord(self.Width/2, self.Height/2)
                self.Icons["bg"].Draw()
        }

        if self.HWND != nil {
                surface.Fill(self.HWND, UI.MySkinManager.GiveColor("White"))
                rect_ := rect.Rect(self.PosX, self.PosY, self.Width, self.Height)
                surface.Blit(self.HWND, self.CanvasHWND, &rect_, nil)
        }
}
