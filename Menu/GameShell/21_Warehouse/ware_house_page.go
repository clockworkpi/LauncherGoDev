package Warehouse

import (
	"context"
	"fmt"
	"log"
	gotime "time"
	"strconv"
	"strings"
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"reflect"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	
	"github.com/zyxar/argo/rpc"

	//"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/time"
	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/rect"
	
	"github.com/cuu/grab"
	"github.com/clockworkpi/LauncherGoDev/sysgo"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

type WareHouse struct {

	UI.Page

	ListFontObj15 *ttf.Font
	ListFontObj12 *ttf.Font
	Icons  map[string]UI.IconItemInterface
	
	BGwidth     int
	BGheight    int
	DrawOnce    bool
	Scroller    *UI.ListScroller
	RemovePage *UI.YesCancelConfirmPage
	Keyboard   *UI.Keyboard
	PreviewPage *ImageDownloadProcessPage
	LoadHousePage *LoadHousePage

	WareHouseDB string
	MyStack  *WareHouseStack
	
	rpcc               rpc.Client
	rpcSecret          string
	rpcURI             string

	Aria2DownloadingGid     string // the Downloading Gid of aria2c
	
	Downloading  chan bool        

	Downloader *grab.Client
	resp       *grab.Response
	req        *grab.Request

	ScrolledCnt   int
	
}

func NewWareHouse() *WareHouse {

	p := &WareHouse{}
	p.ListFontObj12 = UI.MyLangManager.TrFont("notosanscjk12")
	p.ListFontObj15 = UI.MyLangManager.TrFont("varela15")
	p.Icons = make(map[string]UI.IconItemInterface)
	
	p.FootMsg = [5]string{"Nav","Update","Up","Back","Select"}

	p.WareHouseDB = sysgo.SQLDB

	p.BGwidth = 320
	p.BGheight = 240-24-20
	
	p.MyStack = NewWareHouseStack()

	repo := make(map[string]string)
	repo["title"] = "github.com/clockworkpi/warehouse"
	repo["file"]  = "https://raw.githubusercontent.com/clockworkpi/warehouse/master/index.json"
	repo["type"]  = "source"

	p.MyStack.Push(repo)

	p.rpcURI = sysgo.Aria2Url
	
	return p
}

func (self *WareHouse) GetAria2DownloadingPercent(url string) int {
	
	if resp,err := self.rpcc.TellActive();err == nil {
		for _,v := range resp {
			if uris,err := self.rpcc.GetURIs(v.Gid); err == nil {
				for _,x := range uris {
					if x.URI == url {
						comp_len,_ := strconv.ParseInt(v.CompletedLength,10,64)
						totl_len,_ := strconv.ParseInt(v.TotalLength,10,64)
						pct := float64(comp_len)/float64(totl_len)
						pct = pct * 100.0
						return int(pct)
					}
				}	
			}
		}
	}
	return -1;///None
}
func (self *WareHouse) UpdateProcessInterval(ms int) {
	dirty := false
	RefreshTicker := gotime.NewTicker(gotime.Duration(ms)*gotime.Millisecond)
	defer RefreshTicker.Stop()
L:
	for {
		select {
		case <- RefreshTicker.C:
			for _,i := range self.MyList {
				x := i.(*WareHouseListItem)
				if x.Type == "launcher" || x.Type == "pico8" || x.Type == "tic80" {
					percent := self.GetAria2DownloadingPercent(x.Value["file"])
					if percent < 0 {
						x.SetSmallText("")
					}else {
						x.SetSmallText(fmt.Sprintf("%d%%",percent))
						dirty = true
					}
				}
			}

			if self.Screen.CurPage() == self && dirty == true {
				self.Screen.Refresh()
			}
			dirty = false
		case v:= <- self.Downloading:
			if v== false {
				break L
			}
		}
	}
}


func (self *WareHouse) SyncWareHouse() []map[string]string {
	db, err := sql.Open("sqlite3", self.WareHouseDB)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close()

	
	//id,title,file,type
	rows, err := db.Query("select * from warehouse")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ret []map[string]string
	
	for rows.Next() {
		var id int
		var title string
		var file string
		var type_ string
		
		err = rows.Scan(&id, &title,&file,&type_)
		if err != nil {
			log.Fatal(err)
		}

		w_dbt :=  make(map[string]string)
		w_dbt["title"] = title
		w_dbt["file"]  = file
		w_dbt["type"] = type_
		ret = append(ret,w_dbt)

	}
	return ret	
}

func (self *WareHouse) SyncTasks() []map[string]string {
	db, err := sql.Open("sqlite3", self.WareHouseDB)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close()

	
	//id,gid,title,file,type,status,totalLength,completedLength,fav
	rows, err := db.Query("select * from tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ret []map[string]string
	
	for rows.Next() {
		var id int
		var gid string
		var title string
		var file string
		var type_ string
		var status string
		var totalLength string
		var completedLength string
		var fav     string
		
		err = rows.Scan(&id,&gid, &title,&file,&type_,&status,&totalLength,&completedLength,&fav)
		if err != nil {
			log.Fatal(err)
		}

		w_dbt := make(map[string]string)
		w_dbt["gid"] = gid
		w_dbt["title"] = title
		w_dbt["file"] =  file
		w_dbt["type"] = type_
		w_dbt["status"] = status
		w_dbt["totalLength"] = totalLength
		w_dbt["completedLength"] = completedLength
		
		ret = append(ret,w_dbt)

	}
	return ret	
}

func IsSlice(v interface{}) bool {
	if reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array {
		return true
	}
	return false
}

func (self *WareHouse) SyncList()  {

	self.MyList = self.MyList[:0]

	start_x := 0
	start_y := 0

	last_height := 0

	var repos []map[string]string

	
	fmt.Printf("SyncList: %+v\n", self.MyStack)
	
	stk := self.MyStack.Last()
	stk_len := self.MyStack.Length()

	if IsSlice(stk) {
		repos = append(repos, stk.([]map[string]string)...)
	}else {
		repos = append(repos, stk.(map[string]string))
	}
	
	add_new_house := make(map[string]string)
	add_new_house["title"]  = "Add new warehouse..."
	add_new_house["file"]   = "master/index.json"
	add_new_house["type"]   = "add_house"
	add_new_house["status"] = "complete"
	
	if stk_len == 1 {//on top
		ware_menu := self.SyncWareHouse()
		if len(ware_menu) > 0 {
			repos = append(repos,ware_menu...)
		}

		tasks_menu := self.SyncTasks()
		if len(tasks_menu) > 0 {
			repos = append(repos,tasks_menu...)
		}

		repos = append(repos,add_new_house)
	}

	for _, u := range repos {
		fmt.Printf("%+v\n",u)
		li := NewWareHouseListItem()
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y + last_height
		li.Width  = UI.Width
		li.Fonts["normal"] = self.ListFontObj15
		li.Fonts["small"] = self.ListFontObj12
		li.ReadOnly = true
		li.Type = u["type"]
		li.Value = u
		li.Init(u["title"])

		if stk_len > 1 {
			remote_file_url := u["file"]
			menu_file := strings.Split(remote_file_url,"raw.githubusercontent.com")[1]
			home_path, err := os.UserHomeDir()
			if err != nil {
				log.Fatal( err )
			}
			local_menu_file := fmt.Sprintf(aria2dl_folder,home_path,menu_file)
			fmt.Println("for loop ",local_menu_file)
			if UI.FileExists(local_menu_file) {
				li.ReadOnly = false
			}else {
				li.ReadOnly = true
			}
		} else if stk_len == 1 {
			if _,ok := u["status"];ok {
				if u["status"] == "complete" {
					li.ReadOnly = false
				}
			}

			if u["type"] == "source" {
				li.ReadOnly = false
			}
		}

		last_height += li.Height
		if li.Type == "launcher" || li.Type == "pico8" || li.Type == "tic80" {
			li.SetSmallText("")
		}
		self.MyList = append(self.MyList,li)
	}

	self.RefreshPsIndex()
}


func (self *WareHouse) Init()  {

	if self.Screen != nil {
		if self.Screen.CanvasHWND != nil && self.CanvasHWND == nil {
			self.HWND = self.Screen.CanvasHWND
			self.CanvasHWND = surface.Surface(self.Screen.Width, self.BGheight)
		}

		self.PosX = self.Index * self.Screen.Width
		self.Width = self.Screen.Width //equal to screen width
		self.Height = self.Screen.Height

		done := UI.NewIconItem()
		done.ImgSurf = UI.MyIconPool.GetImgSurf("done")
		done.MyType = UI.ICON_TYPES["STAT"]
		done.Parent = self

		self.Icons["done"] = done

		ps := UI.NewInfoPageSelector()
		ps.Parent = self
		self.Ps = ps
		self.PsIndex = 0

		self.Scroller = UI.NewListScroller()
		self.Scroller.Parent = self
		self.Scroller.PosX = self.Width - 10
		self.Scroller.PosY = 2
		self.Scroller.Init()
		self.Scroller.SetCanvasHWND(self.CanvasHWND)

		self.RemovePage = UI.NewYesCancelConfirmPage()
		self.RemovePage.Screen = self.Screen
		self.RemovePage.StartOrAEvent = self.RemoveGame
		self.RemovePage.Name = "Are you sure?"
		
		self.RemovePage.Init()
		
		self.Keyboard = UI.NewKeyboard()
		self.Keyboard.Name = "Enter warehouse addr"
		self.Keyboard.FootMsg = [5]string{"Nav.","Add","ABC","Backspace","Enter"}
		self.Keyboard.Screen = self.Screen
		self.Keyboard.Init()
		self.Keyboard.SetPassword("github.com/clockworkpi/warehouse")
		self.Keyboard.Caller = self
		
		self.PreviewPage = NewImageDownloadProcessPage()
		self.PreviewPage.Screen = self.Screen
		self.PreviewPage.Name ="Preview"
		self.PreviewPage.Init()

		self.LoadHousePage = NewLoadHousePage()
		self.LoadHousePage.Screen = self.Screen
		self.LoadHousePage.Name = "Warehouse"
		self.LoadHousePage.Parent = self
		self.LoadHousePage.Init()
		
		rpcc, err := rpc.New(context.Background(),
			self.rpcURI,
			self.rpcSecret,
			gotime.Second, AppNotifier{Parent:self})
		
	    	if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		self.rpcc = rpcc
		self.Downloader = grab.NewClient()
		self.Downloading = make(chan bool,1)

		self.Screen.HookExitCb(self)

	}

}

func (self *WareHouse) SetDownloading(v bool) {
	for len(self.Downloading) > 0 {
		<- self.Downloading
	}

	self.Downloading <- v
}

func (self *WareHouse) ResetHouse() {
	if self.PsIndex > len(self.MyList) -1 {
		return
	}
	cur_li := self.MyList[self.PsIndex].(*WareHouseListItem)
	home_path, _ := os.UserHomeDir()
	
	if cur_li.Value["type"] == "source" {
		remote_file_url := cur_li.Value["file"]
		parts := strings.Split(remote_file_url,"raw.githubusercontent.com")
		menu_file := parts[1]
		local_menu_file := fmt.Sprintf(aria2dl_folder,home_path,menu_file)
		local_menu_file_path := filepath.Dir(local_menu_file)
		
		fmt.Println(local_menu_file)
		local_jsons,err := filepath.Glob(local_menu_file_path+"/**/*.json")
		if err != nil {
			fmt.Println(err)
		}
		if UI.FileExists(local_menu_file) {
			os.Remove(local_menu_file)
		}
		if UI.FileExists(local_menu_file+".aria2") {
			os.Remove(local_menu_file+".aria2")
		}

		for _,x :=  range local_jsons {
			os.Remove(x)
		}

		self.Screen.MsgBox.SetText("Done")
		self.Screen.MsgBox.Draw()
		self.Screen.SwapAndShow()
		
	}
}

func (self *WareHouse) LoadHouse() {
	if self.PsIndex > len(self.MyList) -1 {
		return
	}

	cur_li := self.MyList[self.PsIndex].(*WareHouseListItem)
	if cur_li.Value["type"] == "source" || cur_li.Value["type"] == "dir" {
		self.LoadHousePage.URL = cur_li.Value["file"]
		self.Screen.PushPage(self.LoadHousePage)
		self.Screen.Refresh()
	}
	
}

func (self *WareHouse) PreviewGame() {
	if self.PsIndex > len(self.MyList) -1 {
		return
	}

	cur_li := self.MyList[self.PsIndex].(*WareHouseListItem)

	if cur_li.Value["type"] == "launcher" ||
		cur_li.Value["type"] == "pico8" ||
		cur_li.Value["type"] == "tic80" {

		if _,ok := cur_li.Value["shots"];ok {
			fmt.Println(cur_li.Value["shots"])
			self.PreviewPage.URL = cur_li.Value["shots"]
			self.Screen.PushPage(self.PreviewPage)
			self.Screen.Refresh()
		}
	}
}
//check if an Url is downloading in aria2c
func (self *WareHouse) UrlIsDownloading(url string) (string,bool) {
	if resp,err := self.rpcc.TellActive();err == nil {
		for _,v := range resp {
			if uris,err := self.rpcc.GetURIs(v.Gid);err == nil {
				for _,x := range uris {
					if x.URI == url {
						fmt.Println(x.URI," ",url)
						return v.Gid,true
					}
				}
			}
		}
	}else {
		log.Println(err)
	}
	return "",false
}

func (self *WareHouse) RemoveGame() {
	if self.PsIndex > len(self.MyList) -1 {
		return
	}
	fmt.Println("RemoveGame")
	cur_li := self.MyList[self.PsIndex].(*WareHouseListItem)

	fmt.Println("Remove cur_li._Value",cur_li.Value)
	home_path, _ := os.UserHomeDir()
	
	if cur_li.Value["type"] == "source" {
		db, err := sql.Open("sqlite3", self.WareHouseDB)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		_, err = db.Exec(fmt.Sprintf("DELETE FROM warehouse WHERE file = '%s'",
			cur_li.Value["file"]))
		
		if err != nil {
			log.Println(err)
		}
		
	} else if cur_li.Value["type"] == "launcher" ||
		cur_li.Value["type"] == "pico8" ||
		cur_li.Value["type"] == "tic80" {

		remote_file_url := cur_li.Value["file"]
		parts := strings.Split(remote_file_url,"raw.githubusercontent.com")
		menu_file := parts[1]
		local_menu_file := fmt.Sprintf(aria2dl_folder,home_path,menu_file)
		local_menu_file_path := filepath.Dir(local_menu_file)

		gid,ret := self.UrlIsDownloading(remote_file_url)
		if ret == true {
			self.rpcc.Remove(gid)
		}

		if UI.FileExists(local_menu_file)  {
			os.Remove(local_menu_file)
		}
		if UI.FileExists(local_menu_file+".aria2") {
			os.Remove(local_menu_file+".aria2")
		}
		if UI.FileExists(filepath.Join(local_menu_file_path,cur_li.Value["title"])) {
			os.RemoveAll(filepath.Join(local_menu_file_path,cur_li.Value["title"]))
		}
		
	}
}

func (self *WareHouse) Click() {
	if self.PsIndex > len(self.MyList) -1 {
		return
	}
	cur_li := self.MyList[self.PsIndex].(*WareHouseListItem)
	home_path, _ := os.UserHomeDir()
	fmt.Println("Click cur_li._Value",cur_li.Value)

	if cur_li.Value["type"] == "source" || cur_li.Value["type"] == "dir" {
		remote_file_url := cur_li.Value["file"]
		parts := strings.Split(remote_file_url,"raw.githubusercontent.com")//assume master branch
		menu_file := parts[1]
		local_menu_file := fmt.Sprintf(aria2dl_folder,home_path,menu_file)
		fmt.Println("warehouse click: ",local_menu_file)
		if UI.FileExists(local_menu_file) == false {
			self.LoadHouse()
		}else {
			//read the local_menu_file,push into stack,display menu
			self.Aria2DownloadingGid = ""
			var result WareHouseIndex
			jsonFile, err := os.Open(local_menu_file)
			
			if err != nil {
        fmt.Println(err)
				self.Screen.MsgBox.SetText("Open House failed")
				self.Screen.MsgBox.Draw()
				self.Screen.SwapAndShow()
				
			}else {
				defer jsonFile.Close()
				
				byteValue, _ := ioutil.ReadAll(jsonFile)
				json.Unmarshal(byteValue, &result)
				self.MyStack.Push(result.List)
				
				self.SyncList()
				self.Screen.Refresh()
			}

			
		}
	} else if cur_li.Value["type"] == "add_house" {
		fmt.Println("show keyboard to add warehouse")
		self.Screen.PushCurPage()
		self.Screen.SetCurPage(self.Keyboard)
		
	} else {
		//download the game probably
		remote_file_url := cur_li.Value["file"]
		parts := strings.Split(remote_file_url,"raw.githubusercontent.com")//assume master branch
		menu_file := parts[1]
		local_menu_file := fmt.Sprintf(aria2dl_folder,home_path,menu_file)
		fmt.Println("Click on game ", local_menu_file)
		
		if UI.FileExists(local_menu_file) == false {
			gid,ret := self.UrlIsDownloading(remote_file_url)
			if ret == false {
				
				outfile := struct {
					Out    string `json:"out"`
					
				}{Out:menu_file}
				
				gid,err := self.rpcc.AddURI([]string{remote_file_url},outfile)
				
				if err != nil {
					log.Println(err)
				}else {
					fmt.Println("Warehouse Click game is downloading, ",gid)
					fmt.Println(remote_file_url)
					self.Aria2DownloadingGid = gid
				}
				
			} else {
				fmt.Println(self.rpcc.TellStatus(gid,"status","totalLength","completedLength"))
				self.Screen.MsgBox.SetText("Getting the game now")
				self.Screen.MsgBox.Draw()
				self.Screen.SwapAndShow()
				time.BlockDelay(800)
				self.Screen.TitleBar.Redraw()
			}
		}else {
			fmt.Println("file downloaded ", cur_li.Value) //maybe check it if is installed fst,then execute it
			if cur_li.Value["type"] == "launcher" && cur_li.ReadOnly == false {
				local_menu_file_path := filepath.Dir(local_menu_file)
				game_sh := filepath.Join(local_menu_file_path,cur_li.Value["title"],cur_li.Value["title"]+".sh")

				fmt.Println("run game: ",game_sh, UI.FileExists(game_sh))
				self.Screen.RunEXE(game_sh)
					
			}
			if cur_li.Value["type"] == "pico8" &&  cur_li.ReadOnly == false {
				if UI.FileExists("/home/cpi/games/PICO-8/pico-8/pico8") {
					game_sh := "/home/cpi/launchergo/Menu/GameShell/50_PICO-8/PICO-8.sh"
					self.Screen.RunEXE(game_sh) //pico8 manages its games self
				}
			}
			if cur_li.Value["type"] == "tic80" && cur_li.ReadOnly == false {
				game_sh := "/home/cpi/apps/Menu/51_TIC-80/TIC-80.sh"
				self.Screen.RunEXE(game_sh)
			}
		}
	}
}

func (self *WareHouse) OnAria2CompleteCb(gid string) {
	fmt.Println("OnAria2CompleteCb", gid)
	self.SyncList()
	self.Screen.Refresh()
	
	if gid == self.Aria2DownloadingGid {
		self.Aria2DownloadingGid = ""
	}
}

func (self *WareHouse) raw_github_com(url string) (bool,string) {
	if strings.HasPrefix(url,"github.com") == false {
		return false,""
	}

	parts := strings.Split(url,"/")

	if len(parts) != 3 {
		return false, ""
	}
	str := []string{"https://raw.githubusercontent.com",
		parts[1],
		parts[2],
		"master/index.json"}
	
	return true,strings.Join(str,"/")
	
}

	
func (self *WareHouse)  OnKbdReturnBackCb() {
	
	inputed:= strings.Join(self.Keyboard.Textarea.MyWords,"")
	inputed = strings.Replace(inputed,"http://","",-1)
	inputed = strings.Replace(inputed,"https://","",-1)

	if strings.HasSuffix(inputed,".git") {
		inputed = inputed[:len(inputed)-4]
	}
	if strings.HasSuffix(inputed,"/") {
		inputed = inputed[:len(inputed)-1]
	}

	fmt.Println("last: ",inputed)
	db, err := sql.Open("sqlite3", self.WareHouseDB)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	
	stmt, err := db.Prepare("SELECT count(*) FROM warehouse WHERE title= ?")
	defer stmt.Close()
	
	if err != nil {
		log.Fatal(err)
	}
	
	var cnt_str string
	cnt := 0
	err = stmt.QueryRow(inputed).Scan(&cnt_str)
	if err != nil {
		log.Println(err)
		cnt_str = "0"
	}else {
		cnt,_= strconv.Atoi(cnt_str)
	}
	
	if cnt > 0 {
		self.Screen.MsgBox.SetText("Warehouse existed")
		self.Screen.MsgBox.Draw()
		self.Screen.SwapAndShow()
	} else {
		if strings.Contains(inputed,"github.com/clockworkpi/warehouse") {
			self.Screen.MsgBox.SetText("Warehouse existed")
			self.Screen.MsgBox.Draw()
			self.Screen.SwapAndShow()
		}else {
			valid_,_url := self.raw_github_com(inputed)
			
			if valid_ == false {
				self.Screen.MsgBox.SetText("Warehouse url error!")
				self.Screen.MsgBox.Draw()
				self.Screen.SwapAndShow()
			} else {
				sql_insert := fmt.Sprintf("INSERT INTO warehouse(title,file,type) VALUES('%s','%s','source');",
					inputed,_url)
				_, err = db.Exec(sql_insert)
				if err != nil {
					log.Println(err)
				}

				self.SyncList()
				self.Screen.Refresh()
			}
		}
	}	
}

func (self *WareHouse) OnExitCb() {
	self.SetDownloading(false)
	self.rpcc.Close()
	
}

func (self *WareHouse) OnLoadCb() {

	if self.MyStack.Length() == 1 {
		self.FootMsg[2] = "Remove"
		self.FootMsg[1] = "Update"
	}else {
		self.FootMsg[2] = "Remove"
		self.FootMsg[1] = "Preview"		
	}

	self.SetDownloading(true)
	go self.UpdateProcessInterval(500)
	
	self.SyncList()
}


func (self *WareHouse)  OnReturnBackCb() {
	if self.MyStack.Length() == 1 {
		self.FootMsg[2] = "Remove"
		self.FootMsg[1] = "Update"
	}else {
		self.FootMsg[2] = "Remove"
		self.FootMsg[1] = "Preview"
	}
	
	self.SyncList()
	self.RestoreScrolled()
	
	self.Screen.Refresh()
	
}

func (self *WareHouse) ScrollDown() {
	if len(self.MyList) == 0 {
		return
	}
	self.PsIndex += 1

	if self.PsIndex >= len(self.MyList) {
		self.PsIndex = len(self.MyList) - 1
	}

	cur_li := self.MyList[self.PsIndex]
	x, y := cur_li.Coord()
	_, h := cur_li.Size()
	
	if y+h > self.Height {
		for i, _ := range self.MyList {
			x, y = self.MyList[i].Coord()
			_, h = self.MyList[i].Size()
			self.MyList[i].NewCoord(x, y-h)
		}
		
		self.ScrolledCnt -= h
	}
}

func (self *WareHouse) ScrollUp() {
	if len(self.MyList) == 0 {
		return
	}

	self.PsIndex -= 1

	if self.PsIndex < 0 {
		self.PsIndex = 0
	}

	cur_li := self.MyList[self.PsIndex]
	x, y := cur_li.Coord()
	_, h := cur_li.Size()
	if y < 0 {
		for i, _ := range self.MyList {
			x, y = self.MyList[i].Coord()
			_, h = self.MyList[i].Size()
			self.MyList[i].NewCoord(x, y+h)
		}

		self.ScrolledCnt += h
	}

}

func (self *WareHouse) RestoreScrolled() {
	for i,_ := range self.MyList {
		x,y := self.MyList[i].Coord()
		self.MyList[i].NewCoord(x, y+ self.ScrolledCnt)
	}
}

func (self *WareHouse) KeyDown(ev *event.Event) {
	if UI.IsKeyMenuOrB(ev.Data["Key"]) {
		if self.MyStack.Length() > 1 {
			self.MyStack.Pop()
			if self.MyStack.Length() == 1 {
				self.FootMsg[2] = "Remove"
				self.FootMsg[1] = "Update"
				
			}else {
				self.FootMsg[2] = "Remove"
				self.FootMsg[1] = "Preview"
				if self.MyStack.Length() == 2 {
					self.FootMsg[2] = ""
					self.FootMsg[1] = ""
				}
			}

			self.SyncList()
			self.Screen.Refresh()
		}else if self.MyStack.Length() == 1 {
			self.ReturnToUpLevelPage()
			self.Screen.Refresh()
			self.SetDownloading(false)//shutdown UpdateProcessInterval
		}
	}
	
	if UI.IsKeyStartOrA(ev.Data["Key"]) {
		self.Click()
		if self.MyStack.Length() == 1 {
			self.FootMsg[2] = "Remove"
			self.FootMsg[1] = "Update"
		}else {
			self.FootMsg[2] = "Remove"
			self.FootMsg[1] = "Preview"
			if self.MyStack.Length() == 2 {
				self.FootMsg[2] = ""
				self.FootMsg[1] = ""
			}
		}

		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["X"] {
		if self.PsIndex <= len(self.MyList) -1 {
			cur_li := self.MyList[self.PsIndex].(*WareHouseListItem)
			if cur_li.Type != "dir" {
				if self.MyStack.Length() ==1 && self.PsIndex == 0 {
					//pass
				}else {
					self.Screen.PushPage(self.RemovePage)
					self.RemovePage.StartOrAEvent = self.RemoveGame
					self.Screen.Refresh()
				}
			}
			return
		}
		self.SyncList()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Y"] {
		if self.MyStack.Length() == 1 {
			self.ResetHouse()
		}else {
			self.PreviewGame()
		}
	}

	if ev.Data["Key"] == UI.CurKeys["Up"] {
		self.ScrollUp()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == UI.CurKeys["Down"] {
		self.ScrollDown()
		self.Screen.Refresh()
	}	
	
}

func (self *WareHouse) Draw() {

	self.ClearCanvas()
	if self.PsIndex > len(self.MyList) -1 {
		self.PsIndex = len(self.MyList) -1
	}
	if self.PsIndex < 0 {
		self.PsIndex = 0
	}
	if len(self.MyList) == 0 {
		return
	} else {
		if len(self.MyList) * UI.DefaultInfoPageListItemHeight > self.Height {
			_,h := self.Ps.Size()
			self.Ps.NewSize(self.Width - 11,h)
			self.Ps.Draw()
			for _,v := range self.MyList {
				_,y := v.Coord()
				if y > (self.Height + self.Height/2) {
					break
				}
				if y < 0 {
					continue
				}
				v.Draw()
			}

			self.Scroller.UpdateSize(len(self.MyList)*UI.DefaultInfoPageListItemHeight,self.PsIndex*UI.DefaultInfoPageListItemHeight)
			self.Scroller.Draw()
		}else {
			_,h := self.Ps.Size()
			self.Ps.NewSize(self.Width,h)
			self.Ps.Draw()
			for _,v := range self.MyList {
				_,y := v.Coord()
				if y > self.Height + self.Height/2 {
					break
				}
				if y < 0 {
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
