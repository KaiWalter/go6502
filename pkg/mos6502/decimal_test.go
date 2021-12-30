package mos6502

import (
	"testing"

	"github.com/KaiWalter/go6502/pkg/addressbus"
	"github.com/KaiWalter/go6502/pkg/memory"
)

const (
	endOfDecimalTest  = 0x024b
	resultDecimalTest = 0x0b
)

func TestDecimal(t *testing.T) {

	// arrange
	ramContent, err := RetrieveROM("6502_decimal_test.bin")
	if err != nil {
		t.Errorf("could not retrieve ROM: %v", err)
	}

	for i := 0; i < 0x1FF; i++ {
		ramContent[i+0x200] = ramContent[i]
		ramContent[i] = 0
	}

	ram := memory.Memory{AddressOffset: 0, AddressSpace: ramContent[:]}
	bus := addressbus.SimpleBus{}
	bus.InitBus(&ram)

	Init(&bus)
	WaitForSystemResetCycles()
	PC = 0x200

	// prevPC := uint16(0xFFFF)
	newInstruction := true

	// act
	for int(PC) != endOfDecimalTest {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}

		if newInstruction {
			// fmt.Printf("%s %04x %04x SP:%02x A:%02x X:%02x Y:%02x abs:%04x fetched:%02x Status:%02x %08b\n",
			// 	opDef.memnonic, CurrentPC, prevPC, SP, A, X, Y,
			// 	absoluteAddress, fetched, Status, Status,
			// )
			// prevPC = CurrentPC
			newInstruction = false
		}

		if CyclesCompleted() {
			newInstruction = true
		}
	}

	// assert
	if ramContent[resultDecimalTest] != 0 {
		t.Errorf("failed - value actual %x / 0 expected", ramContent[resultDecimalTest])
	}
}
