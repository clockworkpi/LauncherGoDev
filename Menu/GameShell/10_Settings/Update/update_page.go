package Update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cuu/gogame/time"
	"github.com/veandco/go-sdl2/ttf"
	"net/http"
	"os/exec"
	"strings"
	gotime "time"
	//	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/event"
	//"github.com/cuu/gogame/rect"
	//"github.com/cuu/gogame/color"
	//	"github.com/cuu/gogame/font"
	//"github.com/cuu/gogame/draw"

	"github.com/clockworkpi/LauncherGoDev/sysgo"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

var InfoPageListItemDefaultHeight = 30
var launchergo_path = "/home/cpi/launchergo"

type UpdateConfirmPage struct {
	UI.ConfirmPage

	URL     string
	MD5     string
	Version string
	GIT     bool
}

func NewUpdateConfirmPage() *UpdateConfirmPage {
	p := &UpdateConfirmPage{}

	p.ListFont = UI.Fonts["veramono20"]
	p.FootMsg = [5]string{"Nav", "", "", "Cancel", "Yes"}
	p.ConfirmText = "Confirm Update?"

	return p
}

func (self *UpdateConfirmPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["B"] {
		fmt.Println("Update Confirm Page B", self.GIT)
		if self.GIT == true {
			//go exec requires FullPath of script filename
			cmdpath := fmt.Sprintf("%s/update.sh %s", UI.GetExePath(), self.Version)
			event.Post(UI.RUNSH, cmdpath)
			return
		}
	}
}

func (self *UpdateConfirmPage) OnReturnBackCb() {
	self.ReturnToUpLevelPage()
	self.Screen.Draw()
	self.Screen.SwapAndShow()
}

func (self *UpdateConfirmPage) Draw() {
	self.ClearCanvas()
	self.DrawBG()
	for _, v := range self.MyList {
		v.Draw()
	}
	self.Reset()
}

type UpdatePage struct {
	UI.Page
	AList       map[string]map[string]string
	ListFontObj *ttf.Font
	MyList      []*UI.InfoPageListItem
	ConfirmPage *UpdateConfirmPage
}

func NewUpdatePage() *UpdatePage {
	p := &UpdatePage{}
	p.FootMsg = [5]string{"Nav", "Check Update", "", "Back", ""}
	p.PageIconMargin = 20
	p.SelectedIconTopOffset = 20
	p.EasingDur = 10

	p.Align = UI.ALIGN["SLeft"]
	p.ListFontObj = UI.Fonts["varela15"]

	p.AList = make(map[string]map[string]string)

	return p
}

func (self *UpdatePage) GenList() {
	self.MyList = nil
	self.MyList = make([]*UI.InfoPageListItem, 0)

	start_x := 0
	start_y := 0
	i := 0

	for k, _ := range self.AList {
		li := UI.NewInfoPageListItem()
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y + i*InfoPageListItemDefaultHeight
		li.Width = UI.Width
		li.Fonts["normal"] = self.ListFontObj
		li.Fonts["small"] = UI.Fonts["varela12"]

		if self.AList[k]["label"] != "" {
			li.Init(self.AList[k]["label"])
		} else {
			li.Init(self.AList[k]["key"])
		}

		li.Flag = self.AList[k]["key"]

		li.SetSmallText(self.AList[k]["value"])

		self.MyList = append(self.MyList, li)

		i += 1
	}
}

func (self *UpdatePage) Init() {
	self.CanvasHWND = self.Screen.CanvasHWND
	self.Width = self.Screen.Width
	self.Height = self.Screen.Height

	self.ConfirmPage = NewUpdateConfirmPage()
	self.ConfirmPage.Screen = self.Screen
	self.ConfirmPage.Name = "Update Confirm"
	self.ConfirmPage.Init()

	it := make(map[string]string)
	it["key"] = "version"
	it["label"] = "Version"
	it["value"] = sysgo.VERSION

	self.AList["version"] = it

	self.GenList()
}

func (self *UpdatePage) CheckUpdate() bool {
	self.Screen.MsgBox.SetText("Checking Update")
	self.Screen.MsgBox.Draw()
	self.Screen.SwapAndShow()

	type Response struct {
		GitVersion string `json:"gitversion"`
	}

	timeout := gotime.Duration(8 * gotime.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(sysgo.UPDATE_URL)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var ret Response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &ret); err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("got git version :", ret.GitVersion)

	launchergo_git_rev_parse := exec.Command("git", "rev-parse", "--short", "HEAD")
	launchergo_git_rev_parse.Dir = launchergo_path
	var out bytes.Buffer
	launchergo_git_rev_parse.Stdout = &out
	err = launchergo_git_rev_parse.Run()

	if err != nil {
		fmt.Println(err)
		return false
	}

	git_revision_short_hash := strings.Trim(out.String(), "\r\n ")

	if git_revision_short_hash != ret.GitVersion {
		self.ConfirmPage.Version = ret.GitVersion
		self.ConfirmPage.GIT = true
		self.Screen.PushCurPage()
		self.Screen.SetCurPage(self.ConfirmPage)

		self.Screen.Draw()

		self.ConfirmPage.SnapMsg(fmt.Sprintf("Update to %s?", ret.GitVersion))
		self.Screen.SwapAndShow()

	} else {
		self.Screen.Draw()
		self.Screen.MsgBox.SetText("Launchergo is up to date")
		self.Screen.MsgBox.Draw()
		self.Screen.SwapAndShow()
		time.BlockDelay(765)
	}

	defer resp.Body.Close()

	return true

}

func (self *UpdatePage) KeyDown(ev *event.Event) {
	if ev.Data["Key"] == UI.CurKeys["A"] || ev.Data["Key"] == UI.CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == UI.CurKeys["X"] {
		if self.Screen.IsWifiConnectedNow() == true {
			if self.CheckUpdate() == true {
				self.Screen.Draw()
				self.Screen.SwapAndShow()
			} else {
				self.Screen.Draw()
				self.Screen.MsgBox.SetText("Check Update Failed")
				self.Screen.MsgBox.Draw()
				self.Screen.SwapAndShow()
			}
		} else {
			self.Screen.Draw()
			self.Screen.MsgBox.SetText("Please Check your Wi-Fi connection")
			self.Screen.MsgBox.Draw()
			self.Screen.SwapAndShow()
		}
	}
}

func (self *UpdatePage) Draw() {
	self.ClearCanvas()
	for _, v := range self.MyList {
		v.Draw()
	}

}
