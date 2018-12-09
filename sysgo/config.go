package sysgo

type PowerLevel struct {
  Dim int
  Close  int
  PowerOff int
}

var PowerLevels map[string]*PowerLevel

var (
	CurKeySet = "GameShell" // PC or GameShell
	DontLeave = false
	BackLight = "/proc/driver/backlight"
	Battery   = "/sys/class/power_supply/axp20x-battery/uevent"
	MPD_socket = "/tmp/mpd.socket"

  UPDATE_URL="https://raw.githubusercontent.com/clockworkpi/CPI/master/launchergo_ver.json"

  VERSION="0.22"

	SKIN="default"
  
  //load from dot files   
  CurPowerLevel= "performance"
  Lang        = "English"
  
)



func init() {
  if PowerLevels == nil {
    PowerLevels = make(map[string]*PowerLevel)
    PowerLevels["supersaving"] = &PowerLevel{10, 30,  120}
    PowerLevels["powersaving"] = &PowerLevel{40, 120, 300}
    PowerLevels["server"]      = &PowerLevel{40, 120, 0  }
    PowerLevels["performance"] = &PowerLevel{40, 0,   0  }
  }
}
