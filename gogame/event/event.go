package event

import (
//	"fmt"

	"strconv"
	
	"github.com/veandco/go-sdl2/sdl"
	
)

const (
	NOEVENT = iota
	QUIT
	KEYDOWN
	KEYUP
	USEREVENT

)

var sdlKeyDict = map[int]string{
sdl.K_UNKNOWN:"",
sdl.K_RETURN:"Return",
sdl.K_ESCAPE:"Escape",
sdl.K_BACKSPACE:"Backspace",
sdl.K_TAB:"Tab",
sdl.K_SPACE:"Space",
sdl.K_EXCLAIM:"!",
sdl.K_QUOTEDBL:"\"",
sdl.K_HASH:"#",
sdl.K_PERCENT:"%",
sdl.K_DOLLAR:"$",
sdl.K_AMPERSAND:"&",
sdl.K_QUOTE:"'",
sdl.K_LEFTPAREN:"(",
sdl.K_RIGHTPAREN:")",
sdl.K_ASTERISK:"*",
sdl.K_PLUS:"+",
sdl.K_COMMA:",",
sdl.K_MINUS:"-",
sdl.K_PERIOD:".",
sdl.K_SLASH:"/",
sdl.K_0:"0",
sdl.K_1:"1",
sdl.K_2:"2",
sdl.K_3:"3",
sdl.K_4:"4",
sdl.K_5:"5",
sdl.K_6:"6",
sdl.K_7:"7",
sdl.K_8:"8",
sdl.K_9:"9",
sdl.K_COLON:":",
sdl.K_SEMICOLON:";",
sdl.K_LESS:"<",
sdl.K_EQUALS:"=",
sdl.K_GREATER:">",
sdl.K_QUESTION:"?",
sdl.K_AT:"@",
sdl.K_LEFTBRACKET:"[",
sdl.K_BACKSLASH:"\\",
sdl.K_RIGHTBRACKET:"]",
sdl.K_CARET:"^",
sdl.K_UNDERSCORE:"_",
sdl.K_BACKQUOTE:"`",
sdl.K_a:"A",
sdl.K_b:"B",
sdl.K_c:"C",
sdl.K_d:"D",
sdl.K_e:"E",
sdl.K_f:"F",
sdl.K_g:"G",
sdl.K_h:"H",
sdl.K_i:"I",
sdl.K_j:"J",
sdl.K_k:"K",
sdl.K_l:"L",
sdl.K_m:"M",
sdl.K_n:"N",
sdl.K_o:"O",
sdl.K_p:"P",
sdl.K_q:"Q",
sdl.K_r:"R",
sdl.K_s:"S",
sdl.K_t:"T",
sdl.K_u:"U",
sdl.K_v:"V",
sdl.K_w:"W",
sdl.K_x:"X",
sdl.K_y:"Y",
sdl.K_z:"Z",
sdl.K_CAPSLOCK:"CapsLock",
sdl.K_F1:"F1",
sdl.K_F2:"F2",
sdl.K_F3:"F3",
sdl.K_F4:"F4",
sdl.K_F5:"F5",
sdl.K_F6:"F6",
sdl.K_F7:"F7",
sdl.K_F8:"F8",
sdl.K_F9:"F9",
sdl.K_F10:"F10",
sdl.K_F11:"F11",
sdl.K_F12:"F12",
sdl.K_PRINTSCREEN:"PrintScreen",
sdl.K_SCROLLLOCK:"ScrollLock",
sdl.K_PAUSE:"Pause",
sdl.K_INSERT:"Insert",
sdl.K_HOME:"Home",
sdl.K_PAGEUP:"PageUp",
sdl.K_DELETE:"Delete",
sdl.K_END:"End",
sdl.K_PAGEDOWN:"PageDown",
sdl.K_RIGHT:"Right",
sdl.K_LEFT:"Left",
sdl.K_DOWN:"Down",
sdl.K_UP:"Up",
sdl.K_NUMLOCKCLEAR:"Numlock",
sdl.K_KP_DIVIDE:"Keypad /",
sdl.K_KP_MULTIPLY:"Keypad *",
sdl.K_KP_MINUS:"Keypad -",
sdl.K_KP_PLUS:"Keypad +",
sdl.K_KP_ENTER:"Keypad Enter",
sdl.K_KP_1:"Keypad 1",
sdl.K_KP_2:"Keypad 2",
sdl.K_KP_3:"Keypad 3",
sdl.K_KP_4:"Keypad 4",
sdl.K_KP_5:"Keypad 5",
sdl.K_KP_6:"Keypad 6",
sdl.K_KP_7:"Keypad 7",
sdl.K_KP_8:"Keypad 8",
sdl.K_KP_9:"Keypad 9",
sdl.K_KP_0:"Keypad 0",
sdl.K_KP_PERIOD:"Keypad .",
sdl.K_APPLICATION:"Application",
sdl.K_POWER:"Power",
sdl.K_KP_EQUALS:"Keypad =",
sdl.K_F13:"F13",
sdl.K_F14:"F14",
sdl.K_F15:"F15",
sdl.K_F16:"F16",
sdl.K_F17:"F17",
sdl.K_F18:"F18",
sdl.K_F19:"F19",
sdl.K_F20:"F20",
sdl.K_F21:"F21",
sdl.K_F22:"F22",
sdl.K_F23:"F23",
sdl.K_F24:"F24",
sdl.K_EXECUTE:"Execute",
sdl.K_HELP:"Help",
sdl.K_MENU:"Menu",
sdl.K_SELECT:"Select",
sdl.K_STOP:"Stop",
sdl.K_AGAIN:"Again",
sdl.K_UNDO:"Undo",
sdl.K_CUT:"Cut",
sdl.K_COPY:"Copy",
sdl.K_PASTE:"Paste",
sdl.K_FIND:"Find",
sdl.K_MUTE:"Mute",
sdl.K_VOLUMEUP:"VolumeUp",
sdl.K_VOLUMEDOWN:"VolumeDown",
sdl.K_KP_COMMA:"Keypad ,",
sdl.K_KP_EQUALSAS400:"Keypad = (AS400)",
sdl.K_ALTERASE:"AltErase",
sdl.K_SYSREQ:"SysReq",
sdl.K_CANCEL:"Cancel",
sdl.K_CLEAR:"Clear",
sdl.K_PRIOR:"Prior",
sdl.K_RETURN2:"Return",
sdl.K_SEPARATOR:"Separator",
sdl.K_OUT:"Out",
sdl.K_OPER:"Oper",
sdl.K_CLEARAGAIN:"Clear / Again",
sdl.K_CRSEL:"CrSel",
sdl.K_EXSEL:"ExSel",
sdl.K_KP_00:"Keypad 00",
sdl.K_KP_000:"Keypad 000",
sdl.K_THOUSANDSSEPARATOR:"ThousandsSeparator",
sdl.K_DECIMALSEPARATOR:"DecimalSeparator",
sdl.K_CURRENCYUNIT:"CurrencyUnit",
sdl.K_CURRENCYSUBUNIT:"CurrencySubUnit",
sdl.K_KP_LEFTPAREN:"Keypad (",
sdl.K_KP_RIGHTPAREN:"Keypad )",
sdl.K_KP_LEFTBRACE:"Keypad {",
sdl.K_KP_RIGHTBRACE:"Keypad }",
sdl.K_KP_TAB:"Keypad Tab",
sdl.K_KP_BACKSPACE:"Keypad Backspace",
sdl.K_KP_A:"Keypad A",
sdl.K_KP_B:"Keypad B",
sdl.K_KP_C:"Keypad C",
sdl.K_KP_D:"Keypad D",
sdl.K_KP_E:"Keypad E",
sdl.K_KP_F:"Keypad F",
sdl.K_KP_XOR:"Keypad XOR",
sdl.K_KP_POWER:"Keypad ^",
sdl.K_KP_PERCENT:"Keypad %",
sdl.K_KP_LESS:"Keypad <",
sdl.K_KP_GREATER:"Keypad >",
sdl.K_KP_AMPERSAND:"Keypad &",
sdl.K_KP_DBLAMPERSAND:"Keypad &&",
sdl.K_KP_VERTICALBAR:"Keypad |",
sdl.K_KP_DBLVERTICALBAR:"Keypad ||",
sdl.K_KP_COLON:"Keypad :",
sdl.K_KP_HASH:"Keypad #",
sdl.K_KP_SPACE:"Keypad Space",
sdl.K_KP_AT:"Keypad @",
sdl.K_KP_EXCLAM:"Keypad !",
sdl.K_KP_MEMSTORE:"Keypad MemStore",
sdl.K_KP_MEMRECALL:"Keypad MemRecall",
sdl.K_KP_MEMCLEAR:"Keypad MemClear",
sdl.K_KP_MEMADD:"Keypad MemAdd",
sdl.K_KP_MEMSUBTRACT:"Keypad MemSubtract",
sdl.K_KP_MEMMULTIPLY:"Keypad MemMultiply",
sdl.K_KP_MEMDIVIDE:"Keypad MemDivide",
sdl.K_KP_PLUSMINUS:"Keypad +/-",
sdl.K_KP_CLEAR:"Keypad Clear",
sdl.K_KP_CLEARENTRY:"Keypad ClearEntry",
sdl.K_KP_BINARY:"Keypad Binary",
sdl.K_KP_OCTAL:"Keypad Octal",
sdl.K_KP_DECIMAL:"Keypad Decimal",
sdl.K_KP_HEXADECIMAL:"Keypad Hexadecimal",
sdl.K_LCTRL:"Left Ctrl",
sdl.K_LSHIFT:"Left Shift",
sdl.K_LALT:"Left Alt",
sdl.K_LGUI:"Left GUI",
sdl.K_RCTRL:"Right Ctrl",
sdl.K_RSHIFT:"Right Shift",
sdl.K_RALT:"Right Alt",
sdl.K_RGUI:"Right GUI",
sdl.K_MODE:"ModeSwitch",
sdl.K_AUDIONEXT:"AudioNext",
sdl.K_AUDIOPREV:"AudioPrev",
sdl.K_AUDIOSTOP:"AudioStop",
sdl.K_AUDIOPLAY:"AudioPlay",
sdl.K_AUDIOMUTE:"AudioMute",
sdl.K_MEDIASELECT:"MediaSelect",
sdl.K_WWW:"WWW",
sdl.K_MAIL:"Mail",
sdl.K_CALCULATOR:"Calculator",
sdl.K_COMPUTER:"Computer",
sdl.K_AC_SEARCH:"AC Search",
sdl.K_AC_HOME:"AC Home",
sdl.K_AC_BACK:"AC Back",
sdl.K_AC_FORWARD:"AC Forward",
sdl.K_AC_STOP:"AC Stop",
sdl.K_AC_REFRESH:"AC Refresh",
sdl.K_AC_BOOKMARKS:"AC Bookmarks",
sdl.K_BRIGHTNESSDOWN:"BrightnessDown",
sdl.K_BRIGHTNESSUP:"BrightnessUp",
sdl.K_DISPLAYSWITCH:"DisplaySwitch",
sdl.K_KBDILLUMTOGGLE:"KBDIllumToggle",
sdl.K_KBDILLUMDOWN:"KBDIllumDown",
sdl.K_KBDILLUMUP:"KBDIllumUp",
sdl.K_EJECT:"Eject",
sdl.K_SLEEP:"Sleep",
}


type Event struct {
	Type uint32
	Data map[string]string
}

func map_events( event sdl.Event) Event {
	var ret Event
		if event != nil {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				ret.Type  = QUIT
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					ret.Type = KEYDOWN
				}
				
				if t.Type == sdl.KEYUP {
					ret.Type = KEYUP
				}
				
				ret.Data = make(map[string]string)
				ret.Data["Repeat"]= strconv.Itoa( int(t.Repeat) )
				ret.Data["Key"] = sdlKeyDict[ int(t.Keysym.Sym) ]
				ret.Data["Mod"] = strconv.Itoa( int(t.Keysym.Mod) )
				
			default:
//				fmt.Printf("unknow type %T\n", t)
				ret.Type = NOEVENT
			}
		}

	return ret
}

func Poll() Event {
	var ret Event
	
	sdl.Do(func() {
		event := sdl.PollEvent()
		ret = map_events(event)
	})

	return ret
}

func Wait() Event {
	var ret Event
	
	event := sdl.WaitEvent()
	ret = map_events(event)
	
	return ret
}
