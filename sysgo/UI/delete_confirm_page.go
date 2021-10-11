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
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if ev.Data["Key"] == CurKeys["B"] {
		err := os.Remove(filepath.Join(self.TrashDir, filepath.Base(self.FileName)))
		if err != nil {
			fmt.Println("DeleteConfirmPage os.Remove errors :", err)
		}

		err = os.Rename(filepath.Base(self.FileName), filepath.Join(self.TrashDir, filepath.Base(self.FileName)))
		if err != nil {
			if strings.Contains(err.Error(), "exists") {
				self.Screen.MsgBox.SetText("Already Existed")
			} else {
				self.Screen.MsgBox.SetText("Error")
			}
			self.Screen.MsgBox.Draw()
			self.Screen.SwapAndShow()
		} else {
			self.SnapMsg("Deleting")
			self.Screen.Draw()
			self.Screen.SwapAndShow()
			self.Reset()
			time.BlockDelay(300)
			self.ReturnToUpLevelPage()
			self.Screen.Draw()
			self.Screen.SwapAndShow()
		}

		fmt.Println(self.FileName)
	}
}
