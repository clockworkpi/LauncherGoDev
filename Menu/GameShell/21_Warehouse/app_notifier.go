package Warehouse

import(
	"log"
	//"os"
	"github.com/zyxar/argo/rpc"
	
)



type AppNotifier struct {
	Parent *WareHouse
}

func (self AppNotifier) OnDownloadStart(events []rpc.Event) {
	log.Printf("warehouse %s started.", events)
}

func (self AppNotifier) OnDownloadPause(events []rpc.Event){
	log.Printf("warehouse %s paused.", events)
}

func (self AppNotifier) OnDownloadStop(events []rpc.Event){
	log.Printf("warehouse %s stopped.", events)
}


func (self AppNotifier) OnDownloadComplete(events []rpc.Event){
	
	log.Printf("warehouse %s complete",events)
	for _,v := range events {
		self.Parent.OnAria2CompleteCb(v.Gid)
	}
}


func (self AppNotifier) OnDownloadError(events []rpc.Event) { 
	log.Printf("warehouse %s error.", events) 
}

func (self AppNotifier) OnBtDownloadComplete(events []rpc.Event) { 
	
	log.Printf("warehouse bt %s completed.", events) 

}

