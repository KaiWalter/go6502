package apple1

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Run() {
	InitScreen()
	defer DestroyScreen()

	// rect := sdl.Rect{0, 0, 200, 200}
	// surface.FillRect(&rect, 0xffff0000)
	// window.UpdateSurface()

	mainLoop()
}

func mainLoop() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
	}
}
