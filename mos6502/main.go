package mos6502

// ReadFunc defines a function where the CPU can read from RAM or Bus
type ReadFunc func(uint16) uint8

// Read points to a function where the CPU can read from RAM or Bus
var Read ReadFunc

// WriteFunc defines a function where the CPU can write to RAM or Bus
type WriteFunc func(uint16, uint8)

// Write points to a function where the CPU can write to RAM or Bus
var Write WriteFunc

// PC = ProgramCounter
var PC uint16

// SP = Stack Pointer
var SP uint8

// A Accumulator
var A uint8

// X Index
var X uint8

// Y Index
var Y uint8

// Test just returns a sample
func Test() string {
	return "XXX"
}
