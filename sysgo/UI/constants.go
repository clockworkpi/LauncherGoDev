package UI

var (
	Width      = 320
	Height     = 240
	IconWidth  = 80
	IconHeight = 80
	IconExt    = ".sh"

	ICON_TYPES = map[string]int{"Emulator": 7, "FILE": 6, "STAT": 5, "NAV": 4, "LETTER": 3, "FUNC": 2, "DIR": 1, "EXE": 0, "None": -1}
	ALIGN      = map[string]int{"HLeft": 0, "HCenter": 1, "HRight": 2, "VMiddle": 3, "SLeft": 4, "VCenter": 5, "SCenter": 6}

	DT = 50
)

var (
	Emulator_flag = "action.config"
	Plugin_flag   = "plugin.json"
)
