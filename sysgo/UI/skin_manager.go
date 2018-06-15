package UI


import (
	"strings"
	
	"github.com/go-ini/ini"
	
	"github.com/cuu/gogame/color"
	
	"../../sysgo"
)

type SkinManager struct {
	Colors map[string]*color.Color
	
}

func NewSkinManager() *SkinManager {
	s := &SkinManager{}

	return s
}


func (self *SkinManager) ConvertToRGB(hexstr string) *color.Color {
	if len(hexstr) < 7 || hexstr[0] != '#' { // # 00 00 00 
		log.Fatalf("ConvertToRGB hex string format error %s", hexstr)
		return nil
	}
	
	h := strings.TrimLeft(hexstr,"#")

	r := strconv.ParseInt(hexstr[0:2], 16,0)
	g := strconv.ParseInt(hexstr[2:4], 16,0)
	b := strconv.ParseInt(hexstr[4:6], 16,0)
	
	col := &color.Color{ r,g,b,255 }
	return col
}

func (self *SkinManager) Init() {
	self.Colors = make(map[string]*color.Color)

	self.Colors["High"] = &color.Color{51,166,255,255}
	self.Colors["Text"] = &color.Color{83,83,83,255}
	self.Colors["Front"] =  &color.Color{131,199,219,255}
	self.Colors["URL"]   = &color.Color{51,166,255,255}
	self.Colors["Line"]  =  &color.Color{169,169,169,255}
	self.Colors["TitleBg"] = &color.Color{228,228,228,255}
	self.Colors["Active"]  =  &color.Color{175,90,0,255}
	self.Colors["White"]  = &color.Color{255,255,255,255}


	fname := "../skin/"+sysgo.SKIN+"/config.cfg"
	
	cfg, err := ini.Load( fname )
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
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

