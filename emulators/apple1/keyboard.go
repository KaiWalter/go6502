package apple1

import (
	"github.com/KaiWalter/go6502/pkg/mc6821"
	"github.com/veandco/go-sdl2/sdl"
)

type keyMap struct {
	unmodified byte
	shifted    byte
	ctrl       byte
}

var (
	keyboardMapping map[sdl.Keycode]keyMap
)

func initKeyboardMapping() {

	// DE!
	keyboardMapping = map[sdl.Keycode]keyMap{
		0x08: {unmodified: 0x5F, shifted: 0x5F, ctrl: 0x5F}, // Apple1 keyboard had not backspace - hence _
		0x0D: {unmodified: 0x0D, shifted: 0x0D, ctrl: 0x0D},

		' ': {unmodified: 0x20, shifted: 0x20, ctrl: 0x20},

		'.': {unmodified: 0x2E, shifted: 0x3A, ctrl: 0x00},
		',': {unmodified: 0x2C, shifted: 0x3B, ctrl: 0x00},
		'+': {unmodified: 0x2B, shifted: 0x2A, ctrl: 0x00},
		'-': {unmodified: 0x2D, shifted: 0x5F, ctrl: 0x00},

		'0': {unmodified: 0x30, shifted: 0x3D, ctrl: 0x00},
		'1': {unmodified: 0x31, shifted: 0x21, ctrl: 0x00},
		'2': {unmodified: 0x32, shifted: 0x22, ctrl: 0x00},
		'3': {unmodified: 0x33, shifted: 0xA7, ctrl: 0x00},
		'4': {unmodified: 0x34, shifted: 0x24, ctrl: 0x00},
		'5': {unmodified: 0x35, shifted: 0x25, ctrl: 0x00},
		'6': {unmodified: 0x36, shifted: 0x26, ctrl: 0x00},
		'7': {unmodified: 0x37, shifted: 0x2F, ctrl: 0x00},
		'8': {unmodified: 0x38, shifted: 0x28, ctrl: 0x00},
		'9': {unmodified: 0x39, shifted: 0x29, ctrl: 0x00},
	}

	// map characters @ A-Z
	for i := sdl.Keycode(0x60); i <= 0x7a; i++ {
		keyboardMapping[i] = keyMap{unmodified: byte(i - 0x20), shifted: byte(i - 0x20), ctrl: byte(i - 0x60)}
	}

}

func handleKeypressed(keysym sdl.Keysym) {
	keyvalue, exists := keyboardMapping[keysym.Sym]
	if exists {

		value := keyvalue.unmodified
		switch keysym.Mod {
		case 0x01:
			value = keyvalue.shifted
		case 0x40:
			value = keyvalue.ctrl
		}

		if value > 0x00 && value < 0x60 {
			piaCA1Channel <- mc6821.Fall            // bring keyboard strobe to low to force active transition
			keyboardInputChannelA <- (value | 0x80) // bit 7 is constantly set (+5V)
			piaCA1Channel <- mc6821.Rise            // send only pulse
			piaCA1Channel <- mc6821.Fall            // 20 micro secs are not worth emulating
		}
	}

}
