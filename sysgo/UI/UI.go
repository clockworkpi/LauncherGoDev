package UI

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/clockworkpi/LauncherGoDev/sysgo"
	"github.com/cuu/gogame/font"
)

type WidgetInterface interface {
	Size() (int, int)
	NewSize(w, h int)
	Coord() (int, int)
	NewCoord(x, y int)
}

type Coord struct {
	X int
	Y int
}

type Plane struct {
	W int
	H int
}

type Widget struct {
	PosX   int
	PosY   int
	Width  int
	Height int
}

func (self *Widget) Size() (int, int) {
	return self.Width, self.Height
}

func (self *Widget) NewSize(w, h int) {
	self.Width = w
	self.Height = h
}

func (self *Widget) Coord() (int, int) {
	return self.PosX, self.PosY
}

func (self *Widget) NewCoord(x, y int) {
	self.PosX = x
	self.PosY = y
}

func FontRW(font_data [] byte, font_size int) *ttf.Font {

	mem,err := sdl.RWFromMem(font_data)
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFontRW(mem, 0, font_size)
	if err != nil {
		panic(fmt.Sprintf("OpenRW font failed %s", err))
	}

	return font
}

type FontData struct {
	Data []byte
	Size int
}
//invoked in main.go
func Init() {
	font.Init()
	
	fonts_name := [4]string{"VarelaRound-Regular.ttf",
				"VeraMono.ttf",
				"NotoSansMono-Regular.ttf",
				"NotoSansCJK-Regular.ttf"}
	
	skinpath := sysgo.SKIN + "/truetype"

	Fonts = make(map[string]*ttf.Font)

	fonts_path := make(map[string]string)
	
	fonts_path["varela"] = fmt.Sprintf("%s/%s", skinpath,fonts_name[0])
	fonts_path["veramono"] = fmt.Sprintf("%s/%s", skinpath,fonts_name[1])
	fonts_path["noto"] = fmt.Sprintf("%s/%s", skinpath,fonts_name[2])
	fonts_path["notocjk"] = fmt.Sprintf("%s/%s", skinpath,fonts_name[3])
	
	fonts_data := make(map[string]FontData)
	d,s := sdl.LoadFile(fonts_path["varela"])	
	fonts_data["varela"] = FontData{d,s}
	d,s = sdl.LoadFile(fonts_path["notocjk"])
	fonts_data["notocjk"] = FontData{d,s}
	d,s = sdl.LoadFile(fonts_path["veramono"])
	fonts_data["veramono"] = FontData{d,s}
	d,s = sdl.LoadFile(fonts_path["noto"])
	fonts_data["noto"] = FontData{d,s}

	for i := 10; i < 41; i++ {
		keyname := fmt.Sprintf("varela%d", i)
		Fonts[keyname] = FontRW(fonts_data["varela"].Data, i)
	}

	Fonts["varela120"] = FontRW(fonts_data["varela"].Data, 120)

	for i := 10; i < 26; i++ {
		keyname := fmt.Sprintf("veramono%d", i)
		Fonts[keyname] = FontRW(fonts_data["veramono"].Data, i)
	}

	for i := 10; i < 28; i++ {
		keyname := fmt.Sprintf("notosansmono%d", i)
		Fonts[keyname] = FontRW(fonts_data["noto"].Data, i)
	}

	for i := 10; i < 28; i++ {
		keyname := fmt.Sprintf("notosanscjk%d", i)
		Fonts[keyname] = FontRW(fonts_data["notocjk"].Data, i)
	}

	//
	keys_def_init()

	//// global variables Init
	if MyIconPool == nil {
		MyIconPool = NewIconPool()
		MyIconPool.Init()
	}
	if MyLangManager == nil {

		MyLangManager = NewLangManager()
		MyLangManager.Init()

	}
	if MySkinManager == nil {
		MySkinManager = NewSkinManager()
		MySkinManager.Init()
	}
}
