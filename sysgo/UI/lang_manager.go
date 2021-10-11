package UI

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/veandco/go-sdl2/ttf"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func sliceToInt(s []int) int {
	res := 0
	op := 1
	for i := len(s) - 1; i >= 0; i-- {
		res += s[i] * op
		op *= 10
	}
	return res
}

func ParseNum(s string) []int {
	nLen := 0
	for i := 0; i < len(s); i++ {
		if b := s[i]; '0' <= b && b <= '9' {
			nLen++
		}
	}
	var n = make([]int, 0, nLen)
	for i := 0; i < len(s); i++ {
		if b := s[i]; '0' <= b && b <= '9' {
			n = append(n, int(b)-'0')
		}
	}
	return n
}

func GetNumberFromString(s string) int {
	is := ParseNum(s)

	return sliceToInt(is)

}

type LangManager struct {
	Langs          map[string]string
	ConfigFilename string

	CJKMode bool
}

func NewLangManager() *LangManager {
	p := &LangManager{}

	p.ConfigFilename = "00_English.ini"
	p.CJKMode = false

	return p
}

func (self *LangManager) Init() {
	if self.Langs == nil {
		self.SetLangs()
	}
}

func (self *LangManager) UpdateLang() {

	self.Langs = nil
	self.SetLangs()

}

func (self *LangManager) IsCJKMode() bool {
	var latins = [1]string{"English"}

	self.CJKMode = false

	for _, v := range latins {
		if strings.HasPrefix(self.ConfigFilename, v) {
			self.CJKMode = false
			break
		}
	}

	return self.CJKMode
}

func (self *LangManager) SetLangs() {

	self.Langs = make(map[string]string)
	fname := "sysgo/.lang"

	load_opts := ini.LoadOptions{
		IgnoreInlineComment: true,
	}

	if FileExists(fname) {
		config_bytes, err := ioutil.ReadFile(fname)
		if err == nil {
			self.ConfigFilename = strings.Trim(string(config_bytes), "\r\n ")
			if len(self.ConfigFilename) < 3 {
				self.ConfigFilename = "00_English.ini"
			}
		}
	} else {
		System("touch " + fname)
	}

	config_file_relative_path := filepath.Join("sysgo", "langs", self.ConfigFilename)

	if FileExists(config_file_relative_path) == false {
		return
	}

	//no matter what ,we must have 00_English.ini
	cfg, err := ini.LoadSources(load_opts, config_file_relative_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		return
	}

	section := cfg.Section("Langs")
	if section != nil {
		opts := section.KeyStrings()
		for _, v := range opts {
			self.Langs[v] = section.Key(v).String()
		}
	}

}

func (self *LangManager) Tr(english_key_str string) string {

	if self.Langs == nil {
		return english_key_str
	}

	if len(self.Langs) == 0 {
		return english_key_str
	}

	if v, ok := self.Langs[english_key_str]; ok {

		return v
	}
	return english_key_str
}

func (self *LangManager) TrFont(orig_font_str string) *ttf.Font {

	font_size_number := GetNumberFromString(orig_font_str)
	if font_size_number > 120 {
		panic("font string format error")
	}

	if strings.Contains(self.ConfigFilename, "English.ini") {
		return Fonts[orig_font_str]
	} else {
		if font_size_number > 28 {
			panic("cjk font size over 28")
		}
	}

	return Fonts[fmt.Sprintf("notosanscjk%d", font_size_number)]

}

var MyLangManager *LangManager
