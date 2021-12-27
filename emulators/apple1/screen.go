package apple1

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	main_window *sdl.Window
)

func InitScreen() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Apple1", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	main_window = window
}

func DestroyScreen() {
	sdl.Quit()
	main_window.Destroy()

}
