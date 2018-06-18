package UI

import (
	"io/ioutil"
	"log"
	"strings"
	
	"github.com/veandco/go-sdl2/sdl"
	
	"github.com/cuu/gogame/image"
	
)

type IconPool struct {
	GameShellIconPath  string
	Icons map[string]*sdl.Surface	
}

func NewIconPool() *IconPool {
	i := &IconPool{}
	i.GameShellIconPath = SkinMap("gameshell/icons/")
	i.Icons = make( map[string]*sdl.Surface )
	return i
}


func (self *IconPool) Init() {
	files,err := ioutil.ReadDir(self.GameShellIconPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _,f := range files {
		if f.IsDir() {
			//pass
		}else {
			if strings.HasSuffix(f.Name(),".png") == true {
				keyname := strings.Split(f.Name(),".")
				if len(keyname) > 1 {
					self.Icons[ keyname[0] ] = image.Load( self.GameShellIconPath+ "/"+f.Name() )
				}
			}
		}
	}
}

func (self *IconPool) GetImgSurf(keyname string) *sdl.Surface {
	if _,ok := self.Icons[keyname]; ok {
		return self.Icons[keyname]
	} else {
		return nil
	}
}

var MyIconPool = NewIconPool()

