package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
	"context"
	"time"
	
	"github.com/zyxar/argo/rpc"

    //"database/sql"
    //_"github.com/mattn/go-sqlite3"
)

// The RPC server might send notifications to the client.
// Notifications is unidirectional, therefore the client which receives the notification must not respond to it.
// The method signature of a notification is much like a normal method request but lacks the id key

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func System(cmd string) string {
	ret := ""
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			//exit code !=0 ,but it can be ignored
		} else {
			fmt.Println(err)
		}
	} else {
		ret = string(out)
	}

	return ret
}

type AppNotifier struct{}

func (AppNotifier) OnDownloadStart(events []rpc.Event)      { log.Printf("%s started.", events) }
func (AppNotifier) OnDownloadPause(events []rpc.Event)      { log.Printf("%s paused.", events) }
func (AppNotifier) OnDownloadStop(events []rpc.Event)       { log.Printf("%s stopped.", events) }

func (AppNotifier) OnDownloadComplete(events []rpc.Event)   { 
	log.Printf("AppNotifier %s completed.", events) 
	
	rpcc_, err := rpc.New(context.Background(), rpcURI, rpcSecret, time.Second,nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer rpcc_.Close()
	
	for _,v := range events {
		gid := v.Gid
		msg, err := rpcc_.TellStatus(gid)
		if err == nil {
			log.Printf("%+v\n",msg)
			if msg.Status == "complete" {
				go InstallGame(msg)
			}else if msg.Status == "error" {
				log.Println(msg.ErrorMessage)
				for _,v := range msg.Files {
					if fileExists(v.Path) {
						e := os.Remove(v.Path)
						if e != nil {
							log.Fatal(e)
						}
					}
					if fileExists(  v.Path + ".aria2" ) {
						e := os.Remove(v.Path + ".aria2")
						if e != nil {
							log.Fatal(e)
						}
					}
				}
			}
		}else {
			log.Println("TellStatus err: ",err)
		}
	}
}

func (AppNotifier) OnDownloadError(events []rpc.Event) { 
	log.Printf("%s error.", events) 
}

func (AppNotifier) OnBtDownloadComplete(events []rpc.Event) { 
	
	log.Printf("bt %s completed.", events) 

}



func InstallGame(msg rpc.StatusInfo) {

	if len(msg.Files) <= 0 {
		return
	}

	ret := msg.Files[0].URIs
	home_path,_ := os.UserHomeDir()
	
	remote_file_url := ret[0].URI

	parts := strings.Split(remote_file_url,"raw.githubusercontent.com")
	if len(parts) < 1 {
		return
	}

	menu_file := parts[1]
	local_menu_file := fmt.Sprintf("%s/aria2downloads%s",home_path,menu_file)
	local_menu_file_path := filepath.Dir(local_menu_file)

	if fileExists(local_menu_file) {
		gametype := "launcher"

		if strings.HasSuffix(local_menu_file,".tar.gz") {
			gametype = "launcher"
		}

		if strings.HasSuffix(local_menu_file,".p8.png") {
			gametype = "pico8"
		}

		if strings.HasSuffix(local_menu_file,".tic") {
			gametype = "tic80"
		}

		if gametype == "launcher" {
			_cmd := fmt.Sprintf( "tar zxvf '%s' -C %s",local_menu_file, local_menu_file_path)
			fmt.Println(_cmd)
			System(_cmd)
		}

		if gametype == "pico8" {
			_cmd := fmt.Sprintf("cp -rf '%s' ~/.lexaloffle/pico-8/carts/", local_menu_file)
			fmt.Println(_cmd)
			System(_cmd)
		}

		if gametype == "tic80" {
			_cmd := fmt.Sprintf("cp -rf '%s' ~/games/TIC-80/",local_menu_file)
			fmt.Println(_cmd)
			System(_cmd)
		}
		
	}else {
		fmt.Println(local_menu_file, " not found")
	}

}

