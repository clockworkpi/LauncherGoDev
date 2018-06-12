package UI

import (
	"log"
	"strconv"
	"bufio"
	"strings"
	
	"github.com/veandco/go-sdl2/sdl"
	
	"github.com/cuu/gogame/surface"
	
	"github.com/itchyny/volume-go"

	"../sysgo"
)


var TitleBar_BarHeight = 24

type TitleBar struct {

	PosX int
	PosY int
	Width int
	Height int
	BarHeight int
	LOffset int
	ROffset int
	Icons map[string]interface{}
	IconWidth
	IconHeight
	BorderWidth
	CanvasHWND *sdl.Surface
	HWND       interface{}
	Title string
	InLowBackLight int
	SkinManager interface{}

	icon_base_path string /// SkinMap("gameshell/titlebar_icons/")
}


func NewTitleBar() *TitleBar {
	t := &TitleBar{}

	
	t.BorderWidth = 1

	t.BarHeight = TitleBar_BarHeight
	t.Height = t.BarHeight + t.BorderWidth

	t.Width = Width
	
	t.Icons = make(map[string]interface{})
	
	//t.icon_base_path  = SkinMap("gameshell/titlebar_icons/")
}

func (t *TitleBar) RoundRobinCheck {
	
}

func (t *TitleBar) UpdateWifiStrength() {
	
}

func (t *TitleBar) GetWifiStrength(stren string) int {
	segs := [][]int{ []int{-2,-1}, []int{0,25}, []int{25,50}, []int{50,75},int{75,100}}
	stren_number,_ :=  strconv.Atoi( stren )
	ge := 0
	if stren_number == 0 {
		return ge
	}
	
	for i,v in range segs {
		if stren_number >= v[0] && stren_number <= v[1] {
			ge = i
			break
		}
	}

	return ge	
}

func (t *TitleBar) SyncSoundVolume() {
	
  vol, err := volume.GetVolume()
  if err != nil {
    log.Fatalf("get volume failed: %+v", err)
		vol = 0
  }
  fmt.Printf("current volume: %d\n", vol)

	snd_segs := [][]int{ []int{0,10}, []int{10,30}, []int{30,70},[]int{70,100} }
	ge := 0

	for i,v in range snd_segs {
		if vol >= v[0] && vol <= v[1] {
			ge = i
			break
		}
	}

	// 
}

func (t *TitleBar) SetSoundVolume(vol int) {
	//pass
}

func (t *TitleBar) CheckBatteryStat() {
	bat_segs:= [][]int{[]int{0,6},[]int{7,15},[]int{16,20},[]int{21,30},[]int{31,50},[]int{51,60},[]int{61,80},[]int{81,90},[]int{91,100}}
	
	file, err := os.Open( sysgo.Battery )
	if err != nil {
		fmt.Println("Could not open file ", sysgo.Battery)
		t.Icons["battery"] = t.Icons["battery_unknown"]
		return
	}

	defer file.Close()

	bat_uevent := make([string]string)
	
  scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines) 

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line," ")
		pis := strings.Split(line,"=")
		if len(pis) > 1 {
			bat_uevent[ pis[0] ] = pis[1]
		}
	}

	cur_cap := 0
	
	if val, ok := bat_uevent["POWER_SUPPLY_CAPACITY"]; ok {
		cur_cap = strings.Atoi(val)
	}else {
		cur_cap = 0
	}

	cap_ge := 0

	for i,v in range bat_segs {
		if cur_cap >= v[0] && cur_cap <= v[1] {
			cap_ge = i
			break
		}
	}

	if val, ok := bat_uevent["POWER_SUPPLY_STATUS"]; ok {
		if val == "Charging" {
			t.Icons["battery_charging"].IconIndex = cap_ge
			t.Icons["battery"] = t.Icons["battery_charging"]
		}else {
			t.Icons["battery_charging"].IconIndex = cap_ge
			t.Icons["battery"] = t.Icons["battery_discharging"]			
		}
	}
	
}

func (t *TitleBar) SetBatteryStat( bat int) {
	
}

func (t *TitleBar) Init(screen *MainScreen) {

	start_x := 0

	t.CanvasHWND = surface.Surface(t.Width,t.Height)
	t.HWND = screen

	icon_wifi_statu := NewMultiIconItem()
	
}
