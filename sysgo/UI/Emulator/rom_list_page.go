package Emulator

import (
	"fmt"
	"os"
	"strings"
	//"regexp"
	"errors"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/event"
	//"github.com/cuu/gogame/time"
	"github.com/veandco/go-sdl2/ttf"
	"os/exec"
	"path/filepath"
	gotime "time"
)

type RomListPage struct {
	UI.Page
	Icons          map[string]UI.IconItemInterface
	ListFont       *ttf.Font
	MyStack        *UI.FolderStack
	EmulatorConfig *ActionConfig

	RomSoConfirmDownloadPage *RomSoConfirmPage

	MyList   []UI.ListItemInterface
	BGwidth  int
	BGheight int //70
	Scroller *UI.ListScroller
	Scrolled int

	Leader *MyEmulator
}

func NewRomListPage() *RomListPage {
	p := &RomListPage{}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]

	p.FootMsg = [5]string{"Nav", "Scan", "Del", "AddFav", "Run"}

	p.Icons = make(map[string]UI.IconItemInterface)
	p.ListFont = UI.Fonts["notosanscjk15"]

	p.MyStack = UI.NewFolderStack()

	p.BGwidth = 56
	p.BGheight = 70

	p.ScrollStep = 1
	return p
}

func (self *RomListPage) GetMyList() []UI.ListItemInterface {
	return self.MyList
}

func (self *RomListPage) GetMapIcons() map[string]UI.IconItemInterface {
	return self.Icons
}

func (self *RomListPage) GetEmulatorConfig() *ActionConfig {
	return self.EmulatorConfig
}

func (self *RomListPage) GeneratePathList(path string) ([]map[string]string, error) {
	if UI.IsDirectory(path) == false {
		return nil, errors.New("Path is not a folder")
	}

	var ret []map[string]string

	file_paths, err := filepath.Glob(path + "/*") //sorted
	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	for _, v := range file_paths {
		dirmap := make(map[string]string)
		if UI.IsDirectory(v) && self.EmulatorConfig.FILETYPE == "dir" { // like DOSBOX
			gameshell_bat := self.EmulatorConfig.EXT[0]
			if UI.GetGid(v) == FavGID { // skip fav roms
				continue
			}

			if UI.FileExists(filepath.Join(v, gameshell_bat)) == true {
				dirmap["gamedir"] = v
				ret = append(ret, dirmap)
			}
		}

		if UI.IsAFile(v) && self.EmulatorConfig.FILETYPE == "file" {
			if UI.GetGid(v) == FavGID {
				continue
			}

			bname := filepath.Base(v)

			if len(bname) > 1 {
				is_excluded := false
				for _, exclude_ext := range self.EmulatorConfig.EXCLUDE {
					exclude_ext2 := strings.Trim(exclude_ext, "\r\n ")
					if len(exclude_ext2) > 1 && strings.HasSuffix(bname, exclude_ext2) {
						is_excluded = true
						break
					}
				}

				if is_excluded == false {
					pieces := strings.Split(bname, ".")

					if len(pieces) > 1 {
						pieces_ext := strings.ToLower(pieces[len(pieces)-1])
						for _, u := range self.EmulatorConfig.EXT {
							if pieces_ext == u {
								dirmap["file"] = v
								ret = append(ret, dirmap)
								break
							}
						}
					}
				}
			}
		}
	}

	return ret, nil

}

func (self *RomListPage) SyncList(path string) {

	alist, err := self.GeneratePathList(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	self.MyList = nil

	start_x := 0
	start_y := 0

	hasparent := 0

	if self.MyStack.Length() > 0 {
		hasparent = 1

		li := NewEmulatorListItem()
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y
		li.Width = UI.Width
		li.Fonts["normal"] = self.ListFont
		li.MyType = UI.ICON_TYPES["DIR"]
		li.Init("[..]")
		self.MyList = append(self.MyList, li)
	}

	for i, v := range alist {
		li := NewEmulatorListItem()
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y + (i+hasparent)*li.Height
		li.Width = UI.Width
		li.Fonts["normal"] = self.ListFont
		li.MyType = UI.ICON_TYPES["FILE"]

		init_val := "NoName"

		if val, ok := v["directory"]; ok {
			li.MyType = UI.ICON_TYPES["DIR"]
			init_val = val
		}

		if val, ok := v["file"]; ok {
			init_val = val
		}

		if val, ok := v["gamedir"]; ok {
			init_val = val
		}

		li.Init(init_val)
		self.MyList = append(self.MyList, li)
	}
}

func (self *RomListPage) Init() {
	self.PosX = self.Index * self.Screen.Width
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.CanvasHWND = self.Screen.CanvasHWND

	ps := UI.NewInfoPageSelector()
	ps.Width = UI.Width - 12
	ps.PosX = 2
	ps.Parent = self

	self.Ps = ps
	self.PsIndex = 0

	self.MyStack.SetRootPath(self.EmulatorConfig.ROM)

	self.SyncList(self.EmulatorConfig.ROM)

	err := os.MkdirAll(self.EmulatorConfig.ROM+"/.Trash", 0700)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(self.EmulatorConfig.ROM+"/.Fav", 0700)
	if err != nil {
		panic(err)
	}

	icon_for_list := UI.NewMultiIconItem()
	icon_for_list.ImgSurf = UI.MyIconPool.GetImgSurf("sys")
	icon_for_list.MyType = UI.ICON_TYPES["STAT"]
	icon_for_list.Parent = self

	icon_for_list.Adjust(0, 0, 18, 18, 0)

	self.Icons["sys"] = icon_for_list

	bgpng := UI.NewIconItem()
	bgpng.ImgSurf = UI.MyIconPool.GetImgSurf("empty")
	bgpng.MyType = UI.ICON_TYPES["STAT"]
	bgpng.Parent = self
	bgpng.AddLabel("Please upload data over Wi-Fi", UI.Fonts["varela22"])
	bgpng.SetLabelColor(&color.Color{204, 204, 204, 255})
	bgpng.Adjust(0, 0, self.BGwidth, self.BGheight, 0)

	self.Icons["bg"] = bgpng

	self.Scroller = UI.NewListScroller()
	self.Scroller.Parent = self
	self.Scroller.PosX = self.Width - 10
	self.Scroller.PosY = 2
	self.Scroller.Init()

	rom_so_confirm_page := NewRomSoConfirmPage()
	rom_so_confirm_page.Screen = self.Screen
	rom_so_confirm_page.Name = "Download Confirm"
	rom_so_confirm_page.Parent = self
	rom_so_confirm_page.Init()

	self.RomSoConfirmDownloadPage = rom_so_confirm_page
}

func (self *RomListPage) ScrollUp() {
	if len(self.MyList) == 0 {
		return
	}
	tmp := self.PsIndex
	self.PsIndex -= self.ScrollStep
	dy := 0

	if self.PsIndex < 0 {
		self.PsIndex = len(self.MyList) - 1
	}

	dy = tmp - self.PsIndex

	cur_li := self.MyList[self.PsIndex]
	x, y := cur_li.Coord()
	_, h := cur_li.Size()
	{
		for i, _ := range self.MyList {
			x, y = self.MyList[i].Coord()
			_, h = self.MyList[i].Size()
			self.MyList[i].NewCoord(x, y+h*dy)
		}

		self.Scrolled += dy
	}
}

func (self *RomListPage) ScrollDown() {
	if len(self.MyList) == 0 {
		return
	}
	tmp := self.PsIndex
	self.PsIndex += self.ScrollStep

	if self.PsIndex >= len(self.MyList) {
		self.PsIndex = 0
	}

	dy := self.PsIndex - tmp

	cur_li := self.MyList[self.PsIndex]
	x, y := cur_li.Coord()
	_, h := cur_li.Size()

	{
		for i, _ := range self.MyList {
			x, y = self.MyList[i].Coord()
			_, h = self.MyList[i].Size()
			self.MyList[i].NewCoord(x, y-h*dy)
		}
		self.Scrolled -= dy
	}

}

func (self *RomListPage) SyncScroll() {

	if self.Scrolled == 0 {
		return
	}

	if self.PsIndex < len(self.MyList) {
		cur_li := self.MyList[self.PsIndex]
		x, y := cur_li.Coord()
		_, h := cur_li.Size()

		if self.Scrolled > 0 {
			if y < 0 {
				for i, _ := range self.MyList {
					x, y = self.MyList[i].Coord()
					_, h = self.MyList[i].Size()
					self.MyList[i].NewCoord(x, y+self.Scrolled*h)
				}
			}
		} else if self.Scrolled < 0 {
			if y+h > self.Height {
				for i, _ := range self.MyList {
					x, y = self.MyList[i].Coord()
					_, h = self.MyList[i].Size()
					self.MyList[i].NewCoord(x, y+self.Scrolled*h)
				}
			}
		}

	}
}

func (self *RomListPage) Click() {
	if len(self.MyList) == 0 {
		return
	}

	if self.PsIndex > len(self.MyList)-1 {
		return
	}

	cur_li := self.MyList[self.PsIndex]

	if cur_li.(*EmulatorListItem).MyType == UI.ICON_TYPES["DIR"] {
		if cur_li.(*EmulatorListItem).Path == "[..]" {
			self.MyStack.Pop()
			self.SyncList(self.MyStack.Last())
			self.PsIndex = 0
		} else {
			self.MyStack.Push(self.MyList[self.PsIndex].(*EmulatorListItem).Path)
			self.SyncList(self.MyStack.Last())
			self.PsIndex = 0
		}
	}

	if cur_li.(*EmulatorListItem).MyType == UI.ICON_TYPES["FILE"] {
		self.Screen.MsgBox.SetText("Launching")
		self.Screen.MsgBox.Draw()
		self.Screen.SwapAndShow()

		path := ""
		if self.EmulatorConfig.FILETYPE == "dir" {
			path = filepath.Join(cur_li.(*EmulatorListItem).Path, self.EmulatorConfig.EXT[0])
		} else {
			path = cur_li.(*EmulatorListItem).Path
		}

		fmt.Println("Run ", path)

		escaped_path := UI.CmdClean(path)

		if self.EmulatorConfig.FILETYPE == "dir" {
			escaped_path = UI.CmdClean(path)
		}

		custom_config := ""

		if self.EmulatorConfig.RETRO_CONFIG != "" && len(self.EmulatorConfig.RETRO_CONFIG) > 5 {
			custom_config = " -c " + self.EmulatorConfig.RETRO_CONFIG
		}

		partsofpath := []string{self.EmulatorConfig.LAUNCHER, self.EmulatorConfig.ROM_SO, custom_config, escaped_path}

		cmdpath := strings.Join(partsofpath, " ")

		if self.EmulatorConfig.ROM_SO == "" { //empty means No needs for rom so
			event.Post(UI.RUNEVT, cmdpath)
		} else {

			if UI.FileExists(strings.Split(self.EmulatorConfig.ROM_SO, " ")[0]) == true {
				event.Post(UI.RUNEVT, cmdpath)
			} else {
				self.Screen.PushCurPage()
				self.Screen.SetCurPage(self.RomSoConfirmDownloadPage)
				self.Screen.Refresh()
			}
		}

		return

	}

	self.Screen.Refresh()
}

func (self *RomListPage) ReScan() {
	//fmt.Println("RomListPage ReScan ",self.EmulatorConfig.ROM)
	if self.MyStack.Length() == 0 {
		self.SyncList(self.EmulatorConfig.ROM)
	} else {
		self.SyncList(self.MyStack.Last())
	}

	self.PsIndex = 0 //sync PsIndex
	self.Scrolled = 0

	self.SyncScroll()
}

func (self *RomListPage) OnReturnBackCb() {
	self.ReScan()
	self.Screen.Refresh()
}

func (self *RomListPage) SpeedScroll(thekey string) {
	if self.Screen.LastKey == thekey {
		self.ScrollStep += 1
		if self.ScrollStep >= self.Leader.SpeedMax {
			self.ScrollStep = self.Leader.SpeedMax
		}
	} else {
		self.ScrollStep = 1
	}
	cur_time := gotime.Now()

	if cur_time.Sub(self.Screen.LastKeyDown) > gotime.Duration(self.Leader.SpeedTimeInter)*gotime.Millisecond {
		self.ScrollStep = 1
	}
}

func (self *RomListPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Right"] {
		self.Screen.PushCurPage()
		self.Screen.SetCurPage(self.Leader.FavPage)
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Up"] {
		self.SpeedScroll(ev.Data["Key"])
		self.ScrollUp()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Down"] {
		self.SpeedScroll(ev.Data["Key"])
		self.ScrollDown()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Enter"] {
		self.Click()
	}

	if ev.Data["Key"] == UI.CurKeys["A"] {
		if len(self.MyList) == 0 {
			return
		}

		cur_li := self.MyList[self.PsIndex]

		if cur_li.(*EmulatorListItem).IsFile() {
			cmd := exec.Command("chgrp", FavGname, UI.CmdClean(cur_li.(*EmulatorListItem).Path))
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}

			self.Screen.ShowMsg("Add to favourite list",600)
			self.ReScan()
			self.Screen.Refresh()
		}
	}

	if ev.Data["Key"] == UI.CurKeys["X"] { //Scan current
		self.ReScan()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Y"] { // del
		if len(self.MyList) == 0 {
			return
		}

		cur_li := self.MyList[self.PsIndex]
		if cur_li.(*EmulatorListItem).IsFile() {
			self.Leader.DeleteConfirmPage.SetFileName(cur_li.(*EmulatorListItem).Path)
			self.Leader.DeleteConfirmPage.SetTrashDir(filepath.Join(self.EmulatorConfig.ROM, "/.Trash"))

			self.Screen.PushCurPage()
			self.Screen.SetCurPage(self.Leader.DeleteConfirmPage)

			self.Screen.Refresh()
		}
	}
}

func (self *RomListPage) Draw() {
	self.ClearCanvas()

	if len(self.MyList) == 0 {
		self.Icons["bg"].NewCoord(self.Width/2, self.Height/2)
		self.Icons["bg"].Draw()
	} else {
		_, h := self.Ps.Size()
		if len(self.MyList)*UI.HierListItemDefaultHeight > self.Height {

			self.Ps.NewSize(self.Width-10, h)
			self.Ps.Draw()
			for _, v := range self.MyList {
				_, y := v.Coord()
				if y > (self.Height + self.Height/2) {
					break
				}
				v.Draw()
			}

			self.Scroller.UpdateSize(len(self.MyList)*UI.HierListItemDefaultHeight,
				self.PsIndex*UI.HierListItemDefaultHeight)
			self.Scroller.Draw()

		} else {
			self.Ps.NewSize(self.Width, h)
			self.Ps.Draw()
			for _, v := range self.MyList {
				v.Draw()
			}
		}
	}
}
