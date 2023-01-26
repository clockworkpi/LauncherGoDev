package UI

import (
	"path/filepath"
)

type CommercialSoftwarePackage struct {
        BinLocation   string
        MenuLocation string
}

func NewCommercialSoftwarePackage(b,m string) *CommercialSoftwarePackage{
	return &CommercialSoftwarePackage{b,m}
}

func (self *CommercialSoftwarePackage) Init()  {

	script := filepath.Join(self.MenuLocation,"Setup.sh")
	MakeExecutable(script)
	script = filepath.Join(self.MenuLocation,"Run.sh")
	MakeExecutable(script)
}

func (self *CommercialSoftwarePackage) IsInstalled() bool {
	return FileExists(self.BinLocation) 
}

func (self *CommercialSoftwarePackage) IsConfiged() bool {
	return FileExists(filepath.Join(self.MenuLocation,".done"))
}

func (self *CommercialSoftwarePackage) GetRunScript() string {
	return filepath.Join(self.MenuLocation,"Run.sh")
}

func (self *CommercialSoftwarePackage) RunSetup() {
	if self.IsConfiged() == false {
		script := filepath.Join(self.MenuLocation,"Setup.sh")
		MakeExecutable(script)
		System(script) /// Scripts with very short runtime
	}
}


