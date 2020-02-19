package mos6502

import (
	"fmt"
	"testing"
)

const (
	endOfDecimalTest  = 0x024b
	resultDecimalTest = 0x0b
)

func TestDecimal(t *testing.T) {

	// arrange
	ram, err := RetrieveROM("6502_decimal_test.bin")
	if err != nil {
		t.Errorf("could not retrieve ROM: %v", err)
	}

	testRead := func(addr uint16) uint8 {
		return ram[addr]
	}

	testWrite := func(addr uint16, data uint8) {
		ram[addr] = data
	}

	for i := 0; i < 0x1FF; i++ {
		ram[i+0x200] = ram[i]
		ram[i] = 0
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 0x200

	prevPC := uint16(0xFFFF)
	newInstruction := true

	// act
	for int(PC) != endOfDecimalTest {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}

		if newInstruction {
			fmt.Printf("%s %04x %04x SP:%02x A:%02x X:%02x Y:%02x abs:%04x fetched:%02x Status:%02x NVUBDIZC%08b | N1%02x N2%02x DA%02x | N1L%02x N2L%02x N1H%02x N2H%02x HNVZC%02x AR%02x\n",
				opDef.memnonic, currentPC, prevPC, SP, A, X, Y,
				absoluteAddress, fetched, Status, Status,
				ram[0x00], ram[0x01], ram[0x04],
				ram[0x0c], ram[0x0e], ram[0x0d], ram[0x0f], ram[0x03], ram[0x06],
			)
			prevPC = currentPC
			newInstruction = false
		}

		if CyclesCompleted() {
			newInstruction = true
		}
	}

	// assert
	if ram[resultDecimalTest] != 0 {
		t.Errorf("failed - value actual %x / 0 expected", ram[resultDecimalTest])
	}
}
