package apple1

import (
	"fmt"

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

	renderCharacter(0, 0, 0x41, false)
	renderCharacter(1, 0, 0x42, true)

	// init PIA
	outputChannel = make(chan uint8, 10)
	mc6821.SetOutputChannelA(outputChannel)
	go ReceiveOutput()

}

func SendOutput(b uint8) {
	outputChannel <- b
}

func ReceiveOutput() {

	for b := range outputChannel {
		fmt.Println("received", b)
	}

}

func DestroyScreen() {
	window.Destroy()
	sdl.Quit()
}
