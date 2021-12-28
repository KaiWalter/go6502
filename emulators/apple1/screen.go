package apple1

import (
	"github.com/veandco/go-sdl2/sdl"

	"github.com/KaiWalter/go6502/pkg/mc6821"
)

const (
	nRows       = 24
	nCols       = 40
	nCharHeight = 8
	nCharWidth  = 8
	nPixelSize  = 4
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer

	cScreenBuffer []uint8
	nCursorY      uint8
	nCursorX      uint8

	outputChannel chan uint8
)

func InitScreen() {

	loadRoms()

	var err error

	// init SDL
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err = sdl.CreateWindow("Apple1", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		nCharWidth*nCols*nPixelSize, nCharHeight*nRows*nPixelSize, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// init screen buffer (for scrolling)
	cScreenBuffer = make([]uint8, nCols*nRows)
	for i := 0; i < len(cScreenBuffer); i++ {
		cScreenBuffer[i] = ' '
	}
	nCursorX, nCursorY = 0, 0

	// init PIA
	outputChannel = make(chan uint8, 10)
	mc6821.SetOutputChannelB(outputChannel)
	go ReceiveOutput()

}

func SendOutput(b uint8) {
	outputChannel <- b
}

func ReceiveOutput() {

	for dsp := range outputChannel {

		// make lower case key upper
		if dsp >= 0x61 && dsp <= 0x7A {
			dsp &= 0x5F
		}

		// clear old cursor
		renderCharacter(nCursorX, nCursorY, cScreenBuffer[nCursorY*nCols+nCursorX], false)

		// display new character
		switch dsp {
		case 0x0D:
			nCursorX = 0
			nCursorY++
		default:
			if dsp >= 0x20 && dsp <= 0x5F {
				cScreenBuffer[nCursorY*nCols+nCursorX] = dsp

				renderCharacter(nCursorX, nCursorY, dsp, false)

				nCursorX++
			}
		}

		// check cursor position
		if nCursorX == nCols {
			nCursorX = 0
			nCursorY++
		}

		if nCursorY == nRows {
			// scroll up
			for y := 0; y < nRows-1; y++ {
				for x := 0; x < nCols; x++ {
					cScreenBuffer[y*nCols+x] = cScreenBuffer[(y+1)*nCols+x]
					renderCharacter(uint8(x), uint8(y), cScreenBuffer[y*nCols+x], false)
				}

			}

			y := (nRows - 1)
			for x := 0; x < nCols; x++ {
				cScreenBuffer[y*nCols+x] = ' '
				renderCharacter(uint8(x), uint8(y), cScreenBuffer[y*nCols+x], false)
			}

			nCursorY--
		}

		// draw new cursor
		renderCharacter(nCursorX, nCursorY, cScreenBuffer[nCursorY*nCols+nCursorX], true)
	}

}

func DestroyScreen() {
	window.Destroy()
	sdl.Quit()
}
