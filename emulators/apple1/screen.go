package apple1

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenRows = 24
	screenCols = 40
	charHeight = 8
	charWidth  = 8
	pixelSize  = 4
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer

	screenBuffer []byte
	cursorY      byte
	cursorX      byte
)

func initScreen() {

	var err error

	// load character ROM
	rom, romInverted, err = loadCharacterRom("./roms/Apple1_charmap.rom", false)
	if err != nil {
		panic(err)
	}

	// init SDL
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("Apple1", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		charWidth*screenCols*pixelSize, charHeight*screenRows*pixelSize, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// init screen buffer (required for scrolling)
	screenBuffer = make([]byte, screenCols*screenRows)
	for i := 0; i < len(screenBuffer); i++ {
		screenBuffer[i] = ' '
	}
	cursorX, cursorY = 0, 0

}

func destroyScreen() {
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()
}

func receiveOutput() {

	for dsp := range screenOutputChannel {

		// make lower case key upper
		if dsp >= 0x61 && dsp <= 0x7A {
			dsp &= 0x5F
		}

		// clear old cursor
		renderCharacter(cursorX, cursorY, screenBuffer[cursorY*screenCols+cursorX], false)

		// display new character
		switch dsp {
		case 0x0D:
			cursorX = 0
			cursorY++
		default:
			if dsp >= 0x20 && dsp <= 0x5F {
				screenBuffer[cursorY*screenCols+cursorX] = dsp

				renderCharacter(cursorX, cursorY, dsp, false)

				cursorX++
			}
		}

		// check cursor position
		if cursorX == screenCols {
			cursorX = 0
			cursorY++
		}

		if cursorY == screenRows {
			// scroll up
			for y := 0; y < screenRows-1; y++ {
				for x := 0; x < screenCols; x++ {
					screenBuffer[y*screenCols+x] = screenBuffer[(y+1)*screenCols+x]
					renderCharacter(byte(x), byte(y), screenBuffer[y*screenCols+x], false)
				}

			}

			y := (screenRows - 1)
			for x := 0; x < screenCols; x++ {
				screenBuffer[y*screenCols+x] = ' '
				renderCharacter(byte(x), byte(y), screenBuffer[y*screenCols+x], false)
			}

			cursorY--
		}

		// draw new cursor
		renderCharacter(cursorX, cursorY, screenBuffer[cursorY*screenCols+cursorX], true)
	}

}
