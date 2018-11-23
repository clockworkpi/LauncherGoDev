package UI


import (
	"fmt"
	
	"log"
	"strings"
	"strconv"
	
	"github.com/go-ini/ini"
	
	"github.com/cuu/gogame/color"
	
	"github.com/cuu/LauncherGo/sysgo"
)

type SkinManager struct {
	Colors map[string]*color.Color
	
}

func NewSkinManager() *SkinManager {
	s := &SkinManager{}

	return s
}


func (self *SkinManager) ConvertToRGB(hexstr string) *color.Color {
	if len(hexstr) < 7 || string(hexstr[0]) != "#" { // # 00 00 00
		log.Fatalf("ConvertToRGB hex string format error %s", hexstr)
		//fmt.Printf("ConvertToRGB hex string format error %s", hexstr)
		return nil
	}
	
	h := strings.TrimLeft(hexstr,"#")

	r,_ := strconv.ParseInt(h[0:2], 16,0)
	g,_ := strconv.ParseInt(h[2:4], 16,0)
	b,_ := strconv.ParseInt(h[4:6], 16,0)
	
	col := &color.Color{ uint32(r),uint32(g),uint32(b),255 }
	return col
}

func (self *SkinManager) ChangeSkin( skin_name string ) {
	
}

func (self *SkinManager) Init() {
	self.Colors = make(map[string]*color.Color)

	self.Colors["High"] = &color.Color{51,166,255,255}
	self.Colors["Text"] = &color.Color{83,83,83,255}
  self.Colors["ReadOnlyText"] = &color.Color{130,130,130,255}
	self.Colors["Front"] =  &color.Color{131,199,219,255}
	self.Colors["URL"]   = &color.Color{51,166,255,255}
	self.Colors["Line"]  =  &color.Color{169,169,169,255}
	self.Colors["TitleBg"] = &color.Color{228,228,228,255}
	self.Colors["Active"]  =  &color.Color{175,90,0,255}
	self.Colors["White"]  = &color.Color{255,255,255,255}
  self.Colors["Black"]  = &color.Color{0,0,0,255}

	fname := "skin/"+sysgo.SKIN+"/config.cfg"

	load_opts := ini.LoadOptions{
		IgnoreInlineComment:true,
	}
	cfg, err := ini.LoadSources(load_opts, fname )
	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		return
	}
	
	section := cfg.Section("Colors")
	if section != nil {
		colour_opts := section.KeyStrings()
		for _,v := range colour_opts {
			if _, ok := self.Colors[v]; ok { // has this Color key
				parsed_color := self.ConvertToRGB( section.Key(v).String() )
				if parsed_color != nil {
					self.Colors[v] = parsed_color
				}
			}
		}
	}
}


func (self *SkinManager) GiveColor(name string) *color.Color {

	if val,ok := self.Colors[name]; ok {
		return val
	}else {
		return &color.Color{255,0,0,255}
	}
}

