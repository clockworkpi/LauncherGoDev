package main

import (	
  //"os"
	"fmt"
  "log"
  "io/ioutil"
//  "strconv"
  "strings"
//  "runtime"
  "path/filepath"
  //os/exec"
  "encoding/json"

  "github.com/cuu/LauncherGoDev/sysgo/UI"
  "github.com/cuu/LauncherGoDev/Menu/GameShell/10_Settings"
)

var (
  UIPluginList = []*UI.UIPlugin {
    &UI.UIPlugin{1,"", "Menu/GameShell/10_Settings",     "Settings",  &Settings.APIOBJ},
  }
)


func ReadTheDirIntoPages(self *UI.MainScreen, _dir string, pglevel int, cur_page UI.PageInterface) {
	
	if UI.FileExists(_dir) == false && UI.IsDirectory(_dir) == false {
		return
	}

	files,err := ioutil.ReadDir(_dir)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _,f := range files { // already sorted
		if UI.IsDirectory( _dir +"/"+f.Name()) {
			if pglevel == 0 {
				page := UI.NewPage()
				page.Name = self.ExtraName(f.Name())
				self.Pages = append(self.Pages, page)
				ReadTheDirIntoPages(self,_dir+"/"+f.Name(),pglevel+1, self.Pages[ len(self.Pages) - 1] )
			}else{ // on cur_page now
				i2:= self.ExtraName(f.Name())
				iconitem := UI.NewIconItem()
				iconitem.AddLabel(i2,self.IconFont)
				if UI.FileExists( UI.SkinMap(_dir+"/"+i2+".png")) {
					iconitem.ImageName = UI.SkinMap(_dir+"/"+i2+".png")
				}else {
					fmt.Println(  UI.SkinMap(_dir+"/"+i2+".png") )
					untitled := UI.NewUntitledIcon()
					untitled.Init()
					if len(i2) > 1 {
						untitled.SetWords(string(i2[0]),string(i2[1]))
					}else if len(i2) == 1 {
						untitled.SetWords(string(i2[0]),string(i2[0]))
					}else {
						untitled.SetWords("G","s")
					}
					iconitem.ImgSurf = untitled.Surface()
					iconitem.ImageName = ""
				}

				if self.IsPluginPackage(_dir+"/"+f.Name()) {
					p_c := UI.PluginConfig{}
          
					dat, err := ioutil.ReadFile(_dir+"/"+f.Name()+"/" +UI.Plugin_flag)
					UI.ShowErr(err)

					err = json.Unmarshal(dat, &p_c)
					if err == nil {
						if p_c.NAME == "" {
							p_c.NAME = f.Name()
						}
            
            so_file := filepath.Join(_dir,f.Name(),p_c.SO_FILE)
            if UI.FileExists(so_file) && UI.IsAFile(so_file) {
              pi,err := UI.LoadPlugin(_dir+"/"+f.Name()+"/"+p_c.SO_FILE)
              UI.Assert(err)
              iconitem.CmdInvoke = UI.InitPlugin(pi,self)
              if iconitem.CmdInvoke != nil {
							
                iconitem.MyType = UI.ICON_TYPES["FUNC"]
                iconitem.CmdPath = f.Name()
                cur_page.AppendIcon(iconitem)
              }
            } else {
              for _,v := range UIPluginList {
                if v.LabelText == p_c.NAME {
                  v.EmbInterface.Init(self)
                  iconitem.CmdInvoke = v.EmbInterface
                  if iconitem.CmdInvoke != nil {
                    iconitem.MyType = UI.ICON_TYPES["FUNC"]
                    iconitem.CmdPath = f.Name()
                    cur_page.AppendIcon(iconitem)
                  }                 
                }
              }
            }
          }
					//Init it 
				}else {
					iconitem.MyType = UI.ICON_TYPES["DIR"]
					linkpage := UI.NewPage()
					linkpage.Name = i2					
					iconitem.LinkPage = linkpage
					cur_page.AppendIcon(iconitem)
					ReadTheDirIntoPages(self,_dir+"/"+f.Name(),pglevel+1, iconitem.LinkPage)
				}
				
			}
		} else if UI.IsAFile(_dir+"/"+f.Name()) && (pglevel > 0) {
			if strings.HasSuffix(strings.ToLower(f.Name()),UI.IconExt) {
				i2 := self.ExtraName(f.Name())
				iconitem := UI.NewIconItem()
				iconitem.CmdPath = _dir+"/"+f.Name()
				UI.MakeExecutable( iconitem.CmdPath )
				iconitem.MyType = UI.ICON_TYPES["EXE"]
				if UI.FileExists( UI.SkinMap( _dir+"/"+ UI.ReplaceSuffix(i2,"png"))) {
					iconitem.ImageName = UI.SkinMap( _dir+"/"+ UI.ReplaceSuffix(i2,"png"))
				}else {
					
					untitled:= UI.NewUntitledIcon()
					untitled.Init()
					if len(i2) > 1 {
						untitled.SetWords(string(i2[0]),string(i2[1]))
					}else if len(i2) == 1 {
						untitled.SetWords(string(i2[0]),string(i2[0]))
					}else {
						untitled.SetWords("G","s")
					}
					iconitem.ImgSurf = untitled.Surface()
					iconitem.ImageName = ""
				}

				iconitem.AddLabel(strings.Split(i2,".")[0], self.IconFont)
				iconitem.LinkPage = nil
				cur_page.AppendIcon(iconitem)
			}
		}
	}
}
