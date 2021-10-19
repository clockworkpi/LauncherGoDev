package Warehouse

import (
	"fmt"
	"log"
	
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	
)

type WareHouse struct {

	UI.Page

	ListFontObj15 *ttf.Font
	ListFontObj12 *ttf.Font
	
	BGwidth     int
	BGheight    int
	DrawOnce    bool
	Scroller    *UI.ListScroller
	RemovePage *UI.YesCancelConfirmPage
	Keyboard   *UI.Keyboard
	
	WareHouseDB string
	MyStack  *WareHouseStack
}

func NewWareHouse() *WareHouse {

	p := &WareHouse{}
	p.ListFontObj12 = UI.MyLangManager.TrFont("notosanscjk12")
	p.ListFontObj15 = UI.MyLangManager.TrFont("varela15")

	p.FootMsg = [5]string{"Nav","Update","Up","Back","Select"}

	p.WareHouseDB = "foo.db"

	p.BGWidth = 320
	p.BGheight = 240-24-20
	
	p.MyStack = NewWareHouseStack()

	repo := make(map[string]string)
	repo["title"] = "github.com/clockworkpi/warehouse"
	repo["file"]  = "https://raw.githubusercontent.com/clockworkpi/warehouse/master/index.json"
	repo["type"]  = "source"

	p.MyStack.Push(repo)
	
	return p
}

func (self*WareHouse) UpdateProcessInterval(ms int) {
	dirty := false
	
}

func (self *WareHouse) SyncWareHouse() []map[string]string {
	db, err := sql.Open("sqlite3", self.WareHouseDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	//id,title,file,type
	rows, err = db.Query("select * from warehouse")
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

		w_dbt := m = make(map[string]string)
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
	}
	defer db.Close()

	
	//id,gid,title,file,type,status,totalLength,completedLength,fav
	rows, err = db.Query("select * from tasks")
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


func (self *WareHouse) SyncList()  {

	self.MyList = self.MyList[:0]

	start_x := 0
	start_y := 0

	last_height := 0

	var repos []map[string]string

	stk := self.MyStack.Last()
	stk_len := self.MyStack.Length()

	repos = append(repos, stk)
	
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

	for i, u := range repos {
		li := WareHouseListItem{}
		li.Parent = self
		li.PosX = start_x
		li.PosY = start_y + last_height
		li.Width  = UI.Width
		li.Fonts["normal"] = self.ListFontObj15
		li.Fonts["small"] = self.ListFontObj12
		li.ReadOnly = true
		li.Type = u["type"]
		li.init(u["title"])

		if stk_len > 1 {
			remote_file_url := u["file"]
			menu_file := strings.Split(remote_file_url,"raw.githubusercontent.com")[1]
			home_path, err := os.UserHomeDir()
			if err != nil {
				log.Fatal( err )
			}
			local_menu_file := fmt.Sprintf("%s/aria2download%s",home_path,menu_file)
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

		done := NewIconItem()
		done.ImgSurf = UI.MyIconPool.GetImgSurf()
		done.MyType = UI.ICON_TYPES["STAT"]
		done.Parent = self

		self.Icons["done"] = done

		ps := UI.NewInfoPageSelector()
		ps.Parent = self
		self.Ps = ps
		self.PsIndex = 0

		self.SyncList()

		self.Scroller = UI.NewListScroller()
		self.Scroller.Parent = self
		self.Scroller.PosX = self.Width - 10
		self.Scroller.PosY = 2
		self.Scroller.Init()
		self.Scroller.SetCanvasHWND(self.CanvasHWND)

		
		
		
	}

}
