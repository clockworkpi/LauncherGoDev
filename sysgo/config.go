package sysgo

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
)

type PowerLevel struct {
	Dim      int
	Close    int
	PowerOff int
}

var PowerLevels map[string]*PowerLevel

var (
	//CurKeySet = "PC" // PC or GameShell
	CurKeySet  = "GameShell"
	DontLeave  = false
	BackLight  = "/proc/driver/backlight"
	Battery    = "/sys/class/power_supply/axp20x-battery/uevent"
	MPD_socket = "/tmp/mpd.socket"

	UPDATE_URL = "https://raw.githubusercontent.com/clockworkpi/CPI/master/launchergo_ver.json"

	VERSION = "0.22"

	SKIN = "skin/default/" // !!!need the last slash!!!

	//load from dot files
	CurPowerLevel = "performance"
	Lang          = "English"
	WifiDev       = "wlan0"

	Aria2Url      = "ws://localhost:6800/jsonrpc"
)

func init() {
	if PowerLevels == nil {
		PowerLevels = make(map[string]*PowerLevel)
		PowerLevels["supersaving"] = &PowerLevel{10, 30, 120}
		PowerLevels["powersaving"] = &PowerLevel{40, 120, 300}
		PowerLevels["server"] = &PowerLevel{40, 120, 0}
		PowerLevels["performance"] = &PowerLevel{40, 0, 0}
	}

	//sudo LauncherGoDev=1 ./launchergo # for develop code on PC
	dev_mode := os.Getenv("LauncherGoDev")

	if len(dev_mode) < 1 {
		return
	}

	if _, err := os.Stat("app-local.ini"); err == nil {
		load_opts := ini.LoadOptions{
			IgnoreInlineComment: true,
		}
		cfg, err := ini.LoadSources(load_opts, "app-local.ini")
		if err != nil {
			fmt.Printf("Fail to read file: %v\n", err)
			return
		}
		section := cfg.Section("GameShell")
		if section != nil {
			gs_opts := section.KeyStrings()
			for i, v := range gs_opts {
				fmt.Println(i, v, section.Key(v).String())
				switch v {
				case "WifiDev":
					WifiDev = section.Key(v).String()
				case "CurKeySet":
					CurKeySet = section.Key(v).String()

				}
			}
		}
	}

}
