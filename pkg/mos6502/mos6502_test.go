package mos6502

// used https://www.masswerk.at/6502/assembler.html to convert from assembler>binary

import (
	"testing"
)

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
	ram := [0x10000]byte{0x00, 0xA9, testValue, 0x85, 0x00}

	testRead := func(addr uint16) byte {
		return ram[addr]
	}

	testWrite := func(addr uint16, data byte) {
		ram[addr] = data
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 1

	// act
	for int(PC) < 5 {
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
	ram := [0x10000]byte{0x00, 0xA9, testValue, 0x8D, 0x00, 0x00}

	testRead := func(addr uint16) byte {
		return ram[addr]
	}

	testWrite := func(addr uint16, data byte) {
		ram[addr] = data
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 1

	// act
	for int(PC) < 6 {
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

func Test_Dec_SBC(t *testing.T) {

	// arrange
	// CLC
	// SED
	// LDA $90
	// SBC $00
	// CMP $89
	ram := [0x10000]byte{0x90, 0x00, 0x89, 0x18, 0xF8, 0xA5, 0x00, 0xE5, 0x01, 0xC5, 0x03}

	testRead := func(addr uint16) byte {
		return ram[addr]
	}

	testWrite := func(addr uint16, data byte) {
		ram[addr] = data
	}

	Init(testRead, testWrite)
	WaitForSystemResetCycles()
	PC = 3

	// act
	for int(PC) < 11 {
		err := Cycle()
		if err != nil {
			t.Errorf("CPU processing failed %v", err)
			break
		}
	}

	// assert
	if A != ram[0x2] {
		t.Errorf("failed - value actual %02x / %02x expected", A, ram[0x2])
	}
}
func checkSetAccumulator(t *testing.T, testValue byte) {
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
