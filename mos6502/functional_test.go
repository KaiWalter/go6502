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
	for int(PC) != endOfMain {
		if PC == 0x0cc1 {
			PC = 0x0cc1
		}
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
			prevPC = currentPC
			fmt.Printf("%x %x %s\n", currentPC, prevPC, opDef.memnonic)
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
