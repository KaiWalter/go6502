package mos6502

// used https://www.masswerk.at/6502/assembler.html to convert from assembler>binary

import "testing"

// wait for system reset cycles
func WaitForSystemResetCycles() {
	for !CyclesCompleted() {
		Cycle()
	}
}

func Test_LDA_IMM_STA_ZP(t *testing.T) {

	// arrange
	testValue := byte(0x42)
	//   lda #$testValue
	//   sta $0
	ram := [...]byte{0x00, 0xA9, testValue, 0x85, 0x00}

	testRead := func(addr uint16) uint8 {
		return ram[addr]
	}

	testWrite := func(addr uint16, data uint8) {
		ram[addr] = data
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 1

	// act
	for int(PC) < len(ram) {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}
	}

	// assert
	checkSetAccumulator(t, testValue)
	if ram[0] != testValue {
		t.Errorf("failed - value %x expected", testValue)
	}
}

func Test_LDA_IMM_STA_ABS(t *testing.T) {

	// arrange
	testValue := byte(0x42)
	//   lda #$testValue
	//   sta $0
	ram := [...]byte{0x00, 0xA9, testValue, 0x8D, 0x00, 0x00}

	testRead := func(addr uint16) uint8 {
		return ram[addr]
	}

	testWrite := func(addr uint16, data uint8) {
		ram[addr] = data
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 1

	// act
	for int(PC) < len(ram) {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}
	}

	// assert
	checkSetAccumulator(t, testValue)
	if ram[0] != testValue {
		t.Errorf("failed - value %x expected", testValue)
	}
}

func checkSetAccumulator(t *testing.T, testValue uint8) {
	if testValue != 0 && GetFlag(Z) {
		t.Errorf("expected Z flag set for set accumulator to 0x%x", testValue)
	} else if testValue == 0 && !GetFlag(Z) {
		t.Errorf("not expected Z flag set for set accumulator to 0x%x", testValue)
	}

	if testValue&0x80 == 0 && GetFlag(N) {
		t.Errorf("expected N flag set for set accumulator to 0x%x", testValue)
	} else if testValue&0x80 != 0 && !GetFlag(N) {
		t.Errorf("not expected N flag set for set accumulator to 0x%x", testValue)
	}
}
