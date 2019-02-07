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
  "sort"
  
  "github.com/yookoala/realpath"


  "github.com/clockworkpi/LauncherGoDev/sysgo/UI"
  "github.com/clockworkpi/LauncherGoDev/sysgo/UI/Emulator"
  "github.com/clockworkpi/LauncherGoDev/Menu/GameShell/10_Settings"
  "github.com/clockworkpi/LauncherGoDev/Menu/GameShell/98_TinyCloud"
  "github.com/clockworkpi/LauncherGoDev/Menu/GameShell/99_PowerOFF"
)

var (
  UIPluginList = []*UI.UIPlugin {
    &UI.UIPlugin{1,"", "Menu/GameShell/10_Settings",     "Settings",  &Settings.APIOBJ},
    &UI.UIPlugin{1,"", "Menu/GameShell/98_TinyCloud",    "TinyCloud", &TinyCloud.APIOBJ},
    &UI.UIPlugin{1,"", "Menu/GameShell/99_PowerOFF",     "PowerOFF",  &PowerOFF.APIOBJ},
  }
)

func ReunionPagesIcons(self *UI.MainScreen) {
    type Tup struct {
      FileName string
      OrigIdx  int
    }
    
    var tmp []Tup
    
    for i,p := range self.Pages {
      p_icons := p.GetIcons()
      for i,x := range p_icons {
        var t Tup
        if x.GetFileName() != ""{
          if strings.Contains(x.GetFileName(),"_") == false {
            t = Tup{"98_"+x.GetFileName(),i}
          }else {
            t = Tup{x.GetFileName(),i}
          }
        }else{
          t = Tup{"",i}
        }
        
        tmp = append(tmp,t)
      }
      
      sort.Slice(tmp, func(i, j int) bool { return tmp[i].FileName < tmp[j].FileName })
      //fmt.Println(tmp)
      
      var retro_games_idx []int
      retro_games_dir := "20_Retro Games"
      for _,x := range tmp {
        if strings.HasPrefix(x.FileName, retro_games_dir) {
          retro_games_idx = append(retro_games_idx,x.OrigIdx)
        }
      }
      
      if len(retro_games_idx) > 1 {
        p_icons_0_link_page := p_icons[retro_games_idx[0]].GetLinkPage().(*UI.Page)
        for i:=1;i<len(retro_games_idx);i++ {
          icons_other_page := p_icons[retro_games_idx[i]].GetLinkPage().GetIcons()
          p_icons_0_link_page.Icons = append(p_icons_0_link_page.Icons, icons_other_page...)
        }
        
        var tmpswap []Tup
        for i,x := range tmp {
          if strings.HasPrefix(x.FileName,retro_games_dir) == false{
            tmpswap = append(tmpswap,x)
          }
          
          if strings.HasPrefix(x.FileName,retro_games_dir) == true && i==retro_games_idx[0] {
            tmpswap = append(tmpswap,x)
          }
        }
        
        tmp = tmpswap
      }
      
      var new_icons []UI.IconItemInterface
      for _,x := range tmp {
        new_icons = append(new_icons, p_icons[x.OrigIdx])
      }
      self.Pages[i].(*UI.Page).Icons = new_icons
    }

}


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
        iconitem.FileName = f.Name()
				iconitem.AddLabel(i2,self.IconFont)
        if UI.FileExists(filepath.Join(_dir,f.Name(),i2+".png")) { //eg: 20_Prog/Prog.png , cut 20_ 
          iconitem.ImageName = filepath.Join(_dir,f.Name(),i2+".png")
          
        }else if UI.FileExists( UI.SkinMap(_dir+"/"+i2+".png")) {
					iconitem.ImageName = UI.SkinMap(_dir+"/"+i2+".png")
				}else {
					//fmt.Println(  UI.SkinMap(_dir+"/"+i2+".png") )
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
				}else if self.IsEmulatorPackage(_dir+"/"+f.Name()) {
          a_c := Emulator.ActionConfig{}
          a_c.FILETYPE="file"
          a_c.TITLE = "Game"
          dat, err := ioutil.ReadFile(_dir+"/"+f.Name()+"/" +UI.Emulator_flag)
					UI.ShowErr(err)

					err = json.Unmarshal(dat, &a_c)
					if err == nil {
            //fmt.Println(a_c)
            if UI.FileExists(filepath.Join(_dir,f.Name(),"retroarch-local.cfg")) {
              a_c.RETRO_CONFIG = UI.CmdClean( filepath.Join(_dir,f.Name(),"retroarch-local.cfg") )
              fmt.Println("a local retroarch cfg: ",a_c.RETRO_CONFIG)
            }
            
            em := Emulator.NewMyEmulator()
            em.EmulatorConfig = &a_c
            em.Init(self)
            
            iconitem.CmdInvoke = em
            if iconitem.CmdInvoke != nil {
              iconitem.MyType = UI.ICON_TYPES["Emulator"]
              iconitem.CmdPath = f.Name()
              cur_page.AppendIcon(iconitem)
            }
          }else {
            fmt.Println("ReadTheDirIntoPages EmulatorConfig ",err)
          }
        
        }else if self.IsExecPackage(_dir+"/"+f.Name()) {
          iconitem.MyType = UI.ICON_TYPES["EXE"]
          rel_path,err := realpath.Realpath( filepath.Join(_dir,f.Name(),i2+".sh"))
          if err != nil {
            rel_path,_ = filepath.Abs(filepath.Join(_dir,f.Name(),i2+".sh"))
          }
          iconitem.CmdPath = rel_path
          UI.MakeExecutable( iconitem.CmdPath )
          cur_page.AppendIcon(iconitem)
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
        rel_path,err := realpath.Realpath( _dir+"/"+f.Name() )
        if err != nil {
          rel_path,_ = filepath.Abs(_dir+"/"+f.Name())
        }
        
				iconitem.CmdPath = rel_path
        iconitem.FileName = f.Name()
        
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
