package apple1

// Reference material:
// http://www.myapplecomputer.net/apple-1-specs.html
// http://www.applefritter.com/book/export/html/22

// Apple 1 HEXROM DISASSEMBLY:
// https://gist.github.com/robey/1bb6a99cd19e95c81979b1828ad70612

// Test ROMs:
// https://github.com/Klaus2m5/6502_65C02_functional_tests

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
	ram []byte

	pia = mc6821.MC6821{Name: "Apple1_PIA", StartAddress: 0xD010, EndAddress: 0xD01F}

	screenOutputChannel   chan byte
	keyboardInputChannelA chan byte
	piaCA1Channel         chan mc6821.Signal
)

func init() {

	initKeyboardMapping()

	ram = make([]byte, 64*1024)
	for i := 0; i < len(ram); i++ {
		ram[i] = 0x00
	}
}

func Run() {
	initScreen()
	defer destroyScreen()

	// load ROMs
	loadROMToAddress("./roms/Apple1_HexMonitor.rom", 0xFF00)
	loadROMToAddress("./roms/Apple1_basic.rom", 0xE000)

	// wire up PIA with screen output and keyboard input
	screenOutputChannel = make(chan byte, 10)
	pia.SetOutputChannelB(screenOutputChannel)
	go receiveOutput()

	keyboardInputChannelA = make(chan byte, 10)
	pia.SetInputChannelA(keyboardInputChannelA)

	piaCA1Channel = make(chan mc6821.Signal, 10)
	pia.SetCA1Channel(piaCA1Channel)

	// init 6502
	mos6502.Init(readMemory, writeMemory)
	waitForSystemResetCycles()

	mainLoop()
}

// wait for system reset cycles
func waitForSystemResetCycles() {
	for !mos6502.CyclesCompleted() {
		mos6502.Cycle()
	}
}

func retrieveROM(filename string) ([]byte, error) {
	romFile, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romFile.Close()

	stats, statsErr := romFile.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	buffer := make([]byte, stats.Size())

	bufferReader := bufio.NewReader(romFile)

	_, err = bufferReader.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return buffer, err
}

func loadROMToAddress(filename string, addr uint16) {
	fmt.Printf("loading ROM %v to %x...\n", filename, addr)

	rom, err := retrieveROM(filename)
	if err != nil {
		fmt.Printf("could not retrieve ROM: %v\n", err)
		return
	}

	for i := 0; i < len(rom); i++ {
		ram[addr+uint16(i)] = rom[i]
	}
}

func readMemory(addr uint16) byte {
	if addr >= pia.StartAddress && addr <= pia.EndAddress {
		return pia.CpuRead(addr)
	}
	return ram[addr]
}

func writeMemory(addr uint16, data byte) {
	if addr >= pia.StartAddress && addr <= pia.EndAddress {
		pia.CpuWrite(addr, data)
	} else {
		ram[addr] = data
	}
}

func mainLoop() {

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				if t.State == 0 {
					handleKeypressed(t.Keysym)
				}
			}
		}

		// TO DO https://floooh.github.io/2019/12/13/cycle-stepped-6502.html

		// if mos6502.CyclesCompleted() {
		// 	fmt.Printf("Current PC %x PC %x\n", mos6502.CurrentPC, mos6502.PC)
		// 	if mos6502.PC == 0xfff4 {
		// 		fmt.Println("SEND TO DISPLAY!")
		// 	}
		// }

		err := mos6502.Cycle()
		if err != nil {
			fmt.Printf("CPU processing failed %v\n", err)
			break
		}

		time.Sleep(5 * time.Millisecond)

	}
}
