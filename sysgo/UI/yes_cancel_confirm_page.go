package UI

import (
	//"fmt"
	//"os"
	//"path/filepath"
	//"strings"

	"github.com/cuu/gogame/event"
	//"github.com/cuu/gogame/time"
	
)

type Func func()

type YesCancelConfirmPage struct {
	ConfirmPage
	StartOrAEvent Func
	KeyXEvent Func
	KeyYEvent Func
}

func NewYesCancelConfirmPage() *YesCancelConfirmPage {
	p := &YesCancelConfirmPage{}
	p.FootMsg = [5]string{"Nav","","","Cancel","Yes"}
	p.ConfirmText = MyLangManager.Tr("Awaiting Input")

	p.StartOrAEvent = nil
	p.KeyXEvent = nil
	p.KeyYEvent = nil
	
	return p
}

func (self *YesCancelConfirmPage) KeyDown(ev *event.Event) {

	if IsKeyMenuOrB(ev.Data["Key"]) {
		self.ReturnToUpLevelPage()
		self.Screen.Draw()
		self.Screen.SwapAndShow()
	}

	if IsKeyStartOrA(ev.Data["Key"]) {
		if self.StartOrAEvent != nil {
			self.StartOrAEvent()
			self.ReturnToUpLevelPage()
		}
	}

	if ev.Data["Key"] == CurKeys["X"] {
		if self.KeyXEvent != nil {
			self.KeyXEvent()
			self.ReturnToUpLevelPage()
		}
	}
	
	if ev.Data["Key"] == CurKeys["Y"] {
		if self.KeyYEvent != nil {
			self.KeyYEvent()
			self.ReturnToUpLevelPage()
		}
	}
	
}
