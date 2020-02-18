package mos6502

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

const (
	endOfMain = 0x3469
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
	for int(currentPC) != endOfMain {
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
			// 	fmt.Printf("%s %04x %04x SP:%02x A:%02x X:%02x Y:%02x abs:%04x fetched:%02x Status:%02x %08b %02x-%02x=%02x\n",
			// 		opDef.memnonic, currentPC, prevPC, SP, A, X, Y,
			// 		absoluteAddress, fetched, Status, Status,
			// 		ram[0x0d], ram[0x12], ram[0x0f],
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

// RetrieveROM retrieves contents of a file into memory
// https://github.com/Klaus2m5/6502_65C02_functional_tests/blob/master/bin_files/6502_functional_test.lst
func RetrieveROM(filename string) ([]byte, error) {
	romfile, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romfile.Close()

	buffer := make([]byte, 0x10000)

	bufferreader := bufio.NewReader(romfile)

	_, err = bufferreader.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return buffer, err
}
