package apple1

import (
	"github.com/KaiWalter/go6502/pkg/mc6821"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	keyboardMapping map[sdl.Keycode]byte
)

func initKeyboardMapping() {

	keyboardMapping = map[sdl.Keycode]byte{
		0x08: 0x08,
		0x0D: 0x0D,
		0x20: 0x20,

		'*': 0x2A,
		'/': 0x2B,
		'+': 0x2F,
		'-': 0x2D,

		'.': 0x2E,
		',': 0x2C,
		// TO DO add missing mappings
	}

	// map characters A-Z
	for i := sdl.Keycode(0x61); i <= 0x7a; i++ {
		keyboardMapping[i] = byte(i - 0x20)
	}
	// map digits 0-9
	for i := sdl.Keycode(0x30); i <= 0x39; i++ {
		keyboardMapping[i] = byte(i)
	}

}

func handleKeypressed(keysym sdl.Keysym) {
	keyvalue, exists := keyboardMapping[keysym.Sym]
	if exists && keyvalue < 0x60 {
		mc6821.SetCA1(mc6821.Fall)        // bring keyboard strobe to low to force active transition
		mc6821.SetInputA(keyvalue | 0x80) // bit 7 is constantly set (+5V)
		mc6821.SetCA1(mc6821.Rise)        // send only pulse
		mc6821.SetCA1(mc6821.Fall)        // 20 micro secs are not worth emulating
	}
}
