package UI

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuu/gogame/event"
	"github.com/cuu/gogame/time"
)

type DeleteConfirmPage struct {
	ConfirmPage
}

func NewDeleteConfirmPage() *DeleteConfirmPage {
	p := &DeleteConfirmPage{}
	p.ListFont = Fonts["veramono20"]
	p.FootMsg = [5]string{"Nav", "", "", "Cancel", "Yes"}
	p.ConfirmText = "Confirm Delete ?"

	return p
}

func (self *DeleteConfirmPage) SetTrashDir(d string) {
	self.TrashDir = d
	if IsDirectory(self.TrashDir) == false {
		panic("DeleteConfirmPage SetTrashDir errors")
	}
}

func (self *DeleteConfirmPage) SetFileName(fn string) {
	self.FileName = fn
}

func (self *DeleteConfirmPage) KeyDown(ev *event.Event) {

	if ev.Data["Key"] == CurKeys["A"] || ev.Data["Key"] == CurKeys["Menu"] {
		self.ReturnToUpLevelPage()
		self.Screen.Refresh()
	}

	if ev.Data["Key"] == CurKeys["B"] {
		err := os.Remove(filepath.Join(self.TrashDir, filepath.Base(self.FileName)))
		if err != nil {
			fmt.Println("DeleteConfirmPage os.Remove errors :", err)
		}

		err = os.Rename(filepath.Base(self.FileName), filepath.Join(self.TrashDir, filepath.Base(self.FileName)))
		if err != nil {
			if strings.Contains(err.Error(), "exists") {
				self.Screen.ShowMsg("Already Existed",0)
			} else {
				self.Screen.ShowMsg("Error",0)
			}
		} else {
			self.SnapMsg("Deleting")
			self.Screen.Refresh()
			self.Reset()
			time.BlockDelay(300)
			self.ReturnToUpLevelPage()
			self.Screen.Refresh()
		}

		fmt.Println(self.FileName)
	}
}
