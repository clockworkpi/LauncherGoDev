package UI

import (
	"os"
	"log"
	"path/filepath"
	"strings"
	"fmt"
	"bufio"
  "bytes"
  
	"github.com/cuu/gogame/display"

	"github.com/cuu/LauncherGo/sysgo"
)

func ShowErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func Assert(e error) {
	if e != nil {
		log.Fatal("Assert: " , e)
	}
}

func CheckAndPanic(e error) {
    if e != nil {
        panic(e)
    }
}

func Abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}


func SkinMap(orig_file_or_dir string) string {
	DefaultSkin := "default"
	ret := ""
	skin_dir_prefix:= "skin/"
	if strings.HasPrefix(orig_file_or_dir, "..") {
		ret = strings.Replace(orig_file_or_dir,"..",skin_dir_prefix + sysgo.SKIN,-1)
		if FileExists(ret) == false {
			ret = strings.Replace(orig_file_or_dir,"..", skin_dir_prefix + DefaultSkin,-1)
		}
	}else {
		ret = skin_dir_prefix+sysgo.SKIN+"/"+orig_file_or_dir
		if FileExists(ret) == false {
			ret = skin_dir_prefix+DefaultSkin+"/"+orig_file_or_dir
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

func ReadLines(path string)(lines [] string,err error){
   var (
       file *os.File
       part [] byte
       prefix bool
   )
   
   if file, err = os.Open(path); err != nil {
       return
   }
   
   reader := bufio.NewReader(file)
   buffer := bytes.NewBuffer(make([]byte,1024))
   
   for {
      if part, prefix, err = reader.ReadLine();err != nil {
          break
      }
      buffer.Write(part)
      if !prefix {
         lines = append(lines,buffer.String())
         buffer.Reset()
      }
   }
   if err == io.EOF {
      err = nil
   }
   return
}

func WriteLines(lines [] string,path string)(err error){
    var file *os.File
    
    if file,err = os.Create(path); err != nil{
         return
    }
    
    defer file.Close()
    
    for _,elem := range lines {
       _,err := file.WriteString(strings.TrimSpace(elem)+"\n")
       if err != nil {
           fmt.Println(err)
           break
       }
    }
    return
}



