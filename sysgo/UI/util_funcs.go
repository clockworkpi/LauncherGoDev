package UI

import (
	"os"
	"log"
	"path/filepath"
	"strings"
	
	"github.com/cuu/gogame/display"

	"../../sysgo"
)

func SkinMap(orig_file_or_dir string) string {
	DefaultSkin := "default"
	ret := ""
	if strings.HasPrefix(orig_file_or_dir, "..") {
		ret = strings.Replace(orig_file_or_dir,"..","../skin/"+sysgo.SKIN,-1)
		if FileExists(ret) == false {
			ret = strings.Replace(orig_file_or_dir,"..", "../skin/"+DefaultSkin)
		}
	}else {
		ret = "../skin/"+sysgo.SKIN+"/sysgo/"+orig_file_or_dir
		if FileExists(ret) == false {
			ret = "../skin/"+DefaultSkin+"/sysgo/"+orig_file_or_dir
		}
	}

	if FileExists(ret) {
		return ret
	}else { // if not existed both in default or custom skin ,return where it is
		return orig_file_or_dir
	}
}

func CmdClean(cmdpath string) string {
	spchars := "\\`$();|{}&'\"*?<>[]!^~-#\n\r "
	for _,v:= range spchars {
		cmdpath = strings.Replace(cmdpath,string(v),"\\"+string(v),-1)
	}
	return cmdpath
}

func FileExists(name string) bool {
	if _, err := os.Stat(name ); err == nil {
		return true
	}else {
		return false
	}
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}else {
    return fileInfo.IsDir()
	}
}


func IsAFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}else {
    return fileInfo.Mode().IsRegular()
	}
}


func MakeExecutable(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatalf("os.Stat %s failed", path)
		return
	}
	mode := fileInfo.Mode()
	mode |= (mode & 0444) >> 2 
	os.Chmod(path,mode)
}

func ReplaceSuffix(orig_file_str string, new_ext string) string {
	orig_ext := filepath.Ext(orig_file_str)
	if orig_ext!= "" {
		las_pos := strings.LastIndex(orig_file_str,".")
		return  orig_file_str[0:las_pos]+"."+new_ext
	}
	return orig_file_str // failed just return back where it came 
}

func SwapAndShow() {
	display.Flip()
}
