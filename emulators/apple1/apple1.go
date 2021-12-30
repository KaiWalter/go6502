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
	"log"
	"os"
	"time"

	"github.com/KaiWalter/go6502/pkg/addressbus"
	"github.com/KaiWalter/go6502/pkg/mc6821"
	"github.com/KaiWalter/go6502/pkg/memory"
	"github.com/KaiWalter/go6502/pkg/mos6502"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	addressMapBlockSize = 0x400 // 1kb
)

var (
	ram  = memory.Memory{AddressOffset: 0, AddressSpace: make([]byte, 4*1024)}
	roms = []memory.Memory{}

	bus = addressbus.MultiBus{}

	pia = mc6821.MC6821{Name: "Apple1_PIA", StartAddress: 0xD010, EndAddress: 0xD01F}

	screenOutputChannel   chan byte
	keyboardInputChannelA chan byte
	piaCA1Channel         chan mc6821.Signal
)

func init() {

	initKeyboardMapping()

	bus.InitBus(addressMapBlockSize)
	bus.RegisterComponent(0, ram.Size()/addressMapBlockSize, &ram)

	// load ROMs
	loadROM("./roms/Apple1_HexMonitor.rom", 0xFF00)
	loadROM("./roms/Apple1_basic.rom", 0xE000)

}

func Run() {
	initScreen()
	defer destroyScreen()

	// wire up PIA with screen output and keyboard input
	bus.RegisterComponent(int(pia.StartAddress), int(pia.EndAddress), &pia)
	screenOutputChannel = make(chan byte, 10)
	pia.SetOutputChannelB(screenOutputChannel)
	go receiveOutput()

	keyboardInputChannelA = make(chan byte, 10)
	pia.SetInputChannelA(keyboardInputChannelA)

	piaCA1Channel = make(chan mc6821.Signal, 10)
	pia.SetCA1Channel(piaCA1Channel)

	// init 6502
	mos6502.Init(&bus)
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

func loadROM(filename string, addr uint16) {
	log.Printf("loading ROM %v to %x...", filename, addr)

	romContent, err := retrieveROM(filename)
	if err != nil {
		log.Printf("could not retrieve ROM: %v", err)
		return
	}

	if len(romContent) == 0 {
		log.Printf("not content in ROM file")
		return
	}

	rom := &memory.Memory{AddressOffset: addr, AddressSpace: romContent}
	bus.RegisterComponent(int(addr), int(addr)+len(romContent)-1, rom)

	roms = append(roms, *rom)
}

func mainLoop() {

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.State == 0 {
					handleKeypressed(t.Keysym)
				}
			}
		}

		if running {

			err := mos6502.Cycle()
			if err != nil {
				log.Printf("CPU processing failed %v", err)
				break
			}

			time.Sleep(5 * time.Millisecond)

		}

	}
}
