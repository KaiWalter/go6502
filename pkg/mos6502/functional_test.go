package mos6502

import (
	"testing"

	"github.com/KaiWalter/go6502/pkg/addressbus"
	"github.com/KaiWalter/go6502/pkg/memory"
)

const (
	endOfFunctionalTest = 0x3469
)

func TestFunctional(t *testing.T) {

	// arrange
	ramContent, err := RetrieveROM("6502_functional_test.bin")
	if err != nil {
		t.Errorf("could not retrieve ROM: %v", err)
	}

	ram := memory.Memory{AddressOffset: 0, AddressSpace: ramContent[:]}
	addressbus.InitBus(0x4000)
	addressbus.RegisterComponent(0, len(ramContent)-1, &ram)

	Init()
	WaitForSystemResetCycles()
	PC = 0x400

	prevPC := uint16(0xFFFF)
	newInstruction := true

	// act
	for int(CurrentPC) != endOfFunctionalTest {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}

		if newInstruction {
			if CurrentPC == prevPC {
				t.Errorf("functional test loops on %x", PC)
				break
			}
			// uncomment for debugging:
			// if CurrentPC >= 0x3480 && CurrentPC <= 0x3489 {
			// 	fmt.Printf("%s %04x %04x SP:%02x A:%02x X:%02x Y:%02x abs:%04x fetched:%02x Status:%02x %08b\n",
			// 		opDef.memnonic, CurrentPC, prevPC, SP, A, X, Y,
			// 		absoluteAddress, fetched, Status, Status,
			// 	)
			// }
			prevPC = CurrentPC
			newInstruction = false
		}

		if CyclesCompleted() {
			newInstruction = true
		}
	}

}
