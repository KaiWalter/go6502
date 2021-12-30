package addressbus

import (
	"testing"

	"github.com/KaiWalter/go6502/pkg/memory"
)

func TestOnlyRam(t *testing.T) {

	const memSize = 0x200

	// arrange
	bus := &MultiBus{}
	bus.InitBus(0x100)
	ram := memory.Memory{AddressOffset: 0, AddressSpace: make([]byte, memSize)}
	bus.RegisterComponent(0, len(ram.AddressSpace)-1, &ram)

	// act & assert
	for addr := uint16(0); addr < memSize; addr++ {
		data, err := bus.Read(addr)
		if err != nil {
			t.Errorf("reading memory failed %v", err)
			break
		}
		if data != 0 {
			t.Errorf("failed - value actual %x / 0 expected", data)
		}
	}

	_, err := bus.Read(memSize + 1)
	if err == nil {
		t.Errorf("expected AddressingError")
	}

}

func TestWithRom(t *testing.T) {

	// arrange
	bus := &MultiBus{}
	bus.InitBus(0x200)
	ram := memory.Memory{AddressOffset: 0, AddressSpace: make([]byte, 0x200)}
	bus.RegisterComponent(0, len(ram.AddressSpace)-1, &ram)

	romContent, err := retrieveROM("dummy01.rom")
	if err != nil {
		t.Errorf("could not retrieve ROM: %v", err)
	}
	rom := memory.Memory{AddressOffset: 0x200, AddressSpace: romContent[:]}
	bus.RegisterComponent(0x200, 0x200+len(romContent)-1, &rom)

	// act & assert
	for addr := uint16(0); addr < 0x200+uint16(len(romContent)); addr++ {
		data, err := bus.Read(addr)
		if err != nil {
			t.Errorf("reading memory failed %v", err)
			break
		}
		if addr < 0x200 && data != 0 {
			t.Errorf("failed at address %x - value actual %x / 0 expected", addr, data)
		}
		if addr >= 0x200 && data != 1 {
			t.Errorf("failed at address %x - value actual %x / 1 expected", addr, data)
		}
	}

	_, err = bus.Read(0x200 + uint16(len(romContent)) + 1)
	if err == nil {
		t.Errorf("expected AddressingError")
	}

}
