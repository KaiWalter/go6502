package apple1

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/KaiWalter/go6502/pkg/mc6821"
	"github.com/KaiWalter/go6502/pkg/mos6502"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ram []uint8
)

func init() {

	initKeyboardMapping()

	ram = make([]uint8, 64*1024)
	for i := 0; i < len(ram); i++ {
		ram[i] = 0x00
	}
}

// wait for system reset cycles
func WaitForSystemResetCycles() {
	for !mos6502.CyclesCompleted() {
		mos6502.Cycle()
	}
}

func retrieveROM(filename string) ([]byte, error) {
	romfile, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romfile.Close()

	stats, statsErr := romfile.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	buffer := make([]byte, stats.Size())

	bufferreader := bufio.NewReader(romfile)

	_, err = bufferreader.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return buffer, err
}

func Run() {
	InitScreen()
	defer DestroyScreen()

	// arrange
	rom, err := retrieveROM("./roms/Apple1_HexMonitor.rom")
	if err != nil {
		fmt.Printf("could not retrieve ROM: %v\n", err)
		return
	}

	for i := 0; i < len(rom); i++ {
		ram[0xFF00+i] = rom[i]
	}
	ram[0xFFFC] = 0x00
	ram[0xFFFD] = 0xFF

	testRead := func(addr uint16) uint8 {
		if addr >= 0xD010 && addr <= 0xD01F {
			return mc6821.CpuRead(addr)
		}
		return ram[addr]
	}

	testWrite := func(addr uint16, data uint8) {
		if addr >= 0xD010 && addr <= 0xD01F {
			mc6821.CpuWrite(addr, data)
		} else {
			ram[addr] = data
		}
	}

	mos6502.Init(testRead, testWrite)
	mos6502.PC = 0xFF00
	WaitForSystemResetCycles()

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

		// TO DO https://floooh.github.io/2019/12/13/cycle-stepped-6502.html

		if mos6502.CyclesCompleted() {
			fmt.Printf("Current PC %x PC %x\n", mos6502.CurrentPC, mos6502.PC)
			if mos6502.PC == 0xfff4 {
				fmt.Println("SEND TO DISPLAY!")
			}
		}

		err := mos6502.Cycle()
		if err != nil {
			fmt.Printf("CPU processing failed %v\n", err)
			break
		}

		time.Sleep(10 * time.Millisecond)

	}
}
