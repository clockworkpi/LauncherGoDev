package UI

import (
	"os"
	"log"
	"path/filepath"
	"strings"
	
	"github.com/cuu/gogame/display"
	
)

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
