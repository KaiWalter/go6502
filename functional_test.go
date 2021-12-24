package main

import (
	"testing"
)

const (
	endOfFunctionalTest = 0x3469
)

func TestFunctional(t *testing.T) {

	// arrange
	ram, err := RetrieveROM("6502_functional_test.bin")
	if err != nil {
		t.Errorf("could not retrieve ROM: %v", err)
	}

	testRead := func(addr uint16) uint8 {
		return ram[addr]
	}

	testWrite := func(addr uint16, data uint8) {
		ram[addr] = data
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 0x400

	prevPC := uint16(0xFFFF)
	newInstruction := true

	// act
	for int(currentPC) != endOfFunctionalTest {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}

		if newInstruction {
			if currentPC == prevPC {
				t.Errorf("functional test loops on %x", PC)
				break
			}
			// uncomment for debugging:
			// if currentPC >= 0x3480 && currentPC <= 0x3489 {
			// 	fmt.Printf("%s %04x %04x SP:%02x A:%02x X:%02x Y:%02x abs:%04x fetched:%02x Status:%02x %08b\n",
			// 		opDef.memnonic, currentPC, prevPC, SP, A, X, Y,
			// 		absoluteAddress, fetched, Status, Status,
			// 	)
			// }
			prevPC = currentPC
			newInstruction = false
		}

		if CyclesCompleted() {
			newInstruction = true
		}
	}

}
