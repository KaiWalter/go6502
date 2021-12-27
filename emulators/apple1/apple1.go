package apple1

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func init() {

	initKeyboardMapping()

}

func Run() {
	InitScreen()
	defer DestroyScreen()

	mainLoop()
}

func mainLoop() {

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				if t.State == 0 {
					handleKeypressed(t.Keysym)
				}
			}
		}
	}
}
